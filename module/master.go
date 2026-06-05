package module

import (
	"context"
	"fmt"
	"github.com/sandwich-go/boost"
	"github.com/sandwich-go/boost/version"
	"github.com/sandwich-go/boost/xdebug"
	"github.com/sandwich-go/boost/xdebug/race"
	"github.com/sandwich-go/boost/xpanic"
	"github.com/sandwich-go/boost/xsync"
	"os"
	"sync"
	"time"
)

type agent struct {
	master *master
	Module
	wg        sync.WaitGroup
	closeChan chan struct{}
}

func (a *agent) run() {
	a.master.runningCount.Add(1)
	a.Run(a.closeChan)
	a.wg.Done()
	if a.master.runningCount.Add(-1) == 0 {
		a.master.Stop(fmt.Sprintf("%s stopped, no module in running", a.Name()))
	}
}

func (a *agent) close() {
	xpanic.Do(func() {
		a.OnClose()
		boost.LogInfof("ModuleName %s closed", a.Name())
	}, func(p *xpanic.Panic) {
		boost.LogInfof("ModuleName %s closed with reason: %v", a.Name(), p.Reason)
	})
}

// master Module管理器
type master struct {
	masterStarted   xsync.AtomicBool
	timeoutDuration time.Duration
	// agentsMu 保护 allAgents 的并发读写：registerOneModule 写 vs
	// Modules / runAll / closeAll 读，以及 RunModule 在 master 已 Run
	// 后注册新 module 时与 closeAll 反向遍历的并发。读侧用
	// snapshotAgents 拷贝快照后无锁遍历。
	agentsMu           sync.Mutex
	allAgents          []*agent
	runningCount       xsync.AtomicInt32
	chanHasShutdown    chan struct{}
	chanStoppedByLogic chan string // 逻辑导致的退出,用户主动停止,逻辑异常停止
	plugins            []Plugin
}

// snapshotAgents 锁内拷贝 allAgents 快照，让读侧遍历期间无需持锁。
func (m *master) snapshotAgents() []*agent {
	m.agentsMu.Lock()
	out := make([]*agent, len(m.allAgents))
	copy(out, m.allAgents)
	m.agentsMu.Unlock()
	return out
}

// New 新建一个 Module 管理器,一般情况下使用默认 default 即可
func New() *master {
	return &master{
		chanHasShutdown:    make(chan struct{}),
		chanStoppedByLogic: make(chan string, 1),
	}
}

func (m *master) Modules() []Module {
	agents := m.snapshotAgents()
	out := make([]Module, 0, len(agents))
	for i := 0; i < len(agents); i++ {
		out = append(out, agents[i].Module)
	}
	return out
}

func (m *master) AttachPlugin(plugins ...Plugin) {
	if len(plugins) > 0 {
		m.plugins = append(m.plugins, plugins...)
	}
}

func (m *master) afterRunModule(ctx context.Context) {
	for _, v := range m.plugins {
		v.AfterRunModule(ctx, m)
	}
}

func (m *master) beforeCloseModule(ctx context.Context) {
	for _, v := range m.plugins {
		v.BeforeCloseModule(ctx, m)
	}
}

func (m *master) registerOneModule(md Module) *agent {
	a := &agent{Module: md, master: m, closeChan: make(chan struct{}, 1)}
	m.agentsMu.Lock()
	m.allAgents = append(m.allAgents, a)
	m.agentsMu.Unlock()
	return a
}

// Register 注册多个Module
func (m *master) Register(ms ...Module) {
	for _, md := range ms {
		_ = m.registerOneModule(md)
	}
}

// Stop 主动停止Master
func (m *master) Stop(reason ...string) {
	var s = "stop_called"
	if len(reason) > 0 {
		s = reason[0]
	}
	select {
	case <-m.chanStoppedByLogic:
	default:
	}
	m.chanStoppedByLogic <- s
}

func (m *master) runAll() {
	agents := m.snapshotAgents()
	for i := 0; i < len(agents); i++ {
		agents[i].OnInit()
	}

	for i := 0; i < len(agents); i++ {
		boost.LogInfof("ModuleName %s starting ...", agents[i].Name())
		agents[i].wg.Add(1)
		go agents[i].run()
		boost.LogInfof("ModuleName %s started ...", agents[i].Name())
	}
}

// RunModule 运行一个单独的 module
func (m *master) RunModule(md Module) {
	s := m.registerOneModule(md)
	if !m.masterStarted.Get() {
		return
	}
	boost.LogInfof("ModuleName %s starting ...", s.Name())
	s.wg.Add(1)
	s.OnInit()
	go s.run()
	boost.LogInfof("ModuleName %s started", s.Name())
}

func (m *master) closeAll(ctx context.Context) {
	agents := m.snapshotAgents()
	for i := len(agents) - 1; i >= 0; i-- {
		a := agents[i]
		boost.LogInfof("ModuleName %s closing ...", a.Name())
		close(a.closeChan)
		if m.timeoutDuration == 0 {
			a.wg.Wait()
		} else {
			if xsync.WaitContext(&a.wg, ctx) {
				boost.LogInfof("ModuleName %s close with timeout %v ...", a.Name(), m.timeoutDuration)
			}
		}
		a.close()
	}
}

// RunWithCloseTimeout 参考Run,扩充了关闭超时支持，防止逻辑层堵塞导致进程关闭失败
func (m *master) RunWithCloseTimeout(duration time.Duration, ms ...Module) {
	m.timeoutDuration = duration
	m.Run(ms...)
}

// ShutdownNotify master停止的通知信号
func (m *master) ShutdownNotify() chan struct{} { return m.chanHasShutdown }

// Run 运行入口，进程会堵塞在这里直到收到停止信号，可以指定要运行的Module列表
func (m *master) Run(ms ...Module) {
	m.Register(ms...)
	m.runAll()

	xdebug.CheckRequireDependencies()
	boost.LogInfof("progress started, pid: %d, version: %s, race: %t, debug_enabled: %t",
		os.Getpid(), version.String(), race.Enabled, xdebug.Enabled())

	// runCtx 给 afterRunModule 用；按值传给 goroutine 闭包，避免主 goroutine
	// 后续 closeCtx 赋值时引发 stack 变量 race。
	runCtx := context.Background()
	go func(ctx context.Context) {
		m.afterRunModule(ctx)
	}(runCtx)
	m.masterStarted.Set(true)

	reason := "unknown"
	// Block until a signal is received
	select {
	case reason = <-m.chanStoppedByLogic:
	case <-ProcessShutdownNotify():
		reason = fmt.Sprintf("sig(%s)", processShutdownSignal.String())
	}

	boost.LogInfof("progress closing down by signal, pid: %d, reason: %s", os.Getpid(), reason)

	m.masterStarted.Set(false)

	// closeCtx 用独立变量名构造，不复用 runCtx 的 stack slot，避免与
	// afterRunModule goroutine 之间 race。
	closeCtx := context.Background()
	if m.timeoutDuration != 0 {
		var cancelFunc context.CancelFunc
		closeCtx, cancelFunc = context.WithDeadline(closeCtx, time.Now().Add(m.timeoutDuration))
		defer cancelFunc()
	}
	beforeCloseModuleDone := make(chan struct{})
	go func(ctx context.Context) {
		m.beforeCloseModule(ctx)
		close(beforeCloseModuleDone)
	}(closeCtx)

	select {
	case <-beforeCloseModuleDone:
	case <-closeCtx.Done():
	}

	m.closeAll(closeCtx)
	close(m.chanHasShutdown)
}
