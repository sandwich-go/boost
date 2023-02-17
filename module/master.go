package module

import (
	"context"
	"fmt"
	"github.com/sandwich-go/boost"
	"github.com/sandwich-go/boost/paniccatcher"
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
	paniccatcher.Do(func() {
		a.OnClose()
		boost.LogInfof("ModuleName %s closed", a.Name())
	}, func(p *paniccatcher.Panic) {
		boost.LogInfof("ModuleName %s closed with reason: %v", a.Name(), p.Reason)
	})
}

// master Module管理器
type master struct {
	masterStarted      xsync.AtomicBool
	timeoutDuration    time.Duration
	allAgents          []*agent
	runningCount       xsync.AtomicInt32
	chanHasShutdown    chan struct{}
	chanStoppedByLogic chan string // 逻辑导致的退出,用户主动停止,逻辑异常停止
	plugins            []Plugin
}

// New 新建一个 Module 管理器,一般情况下使用默认 default 即可
func New() *master {
	return &master{
		chanHasShutdown:    make(chan struct{}),
		chanStoppedByLogic: make(chan string, 1),
	}
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
	m.allAgents = append(m.allAgents, a)
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
	for i := 0; i < len(m.allAgents); i++ {
		m.allAgents[i].OnInit()
	}

	for i := 0; i < len(m.allAgents); i++ {
		boost.LogInfof("ModuleName %s starting ...", m.allAgents[i].Name())
		m.allAgents[i].wg.Add(1)
		go m.allAgents[i].run()
		boost.LogInfof("ModuleName %s started ...", m.allAgents[i].Name())
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
	for i := len(m.allAgents) - 1; i >= 0; i-- {
		a := m.allAgents[i]
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

	ctx := context.Background()
	go func() {
		m.afterRunModule(ctx)
	}()
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

	if m.timeoutDuration != 0 {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithDeadline(ctx, time.Now().Add(m.timeoutDuration))
		defer cancelFunc()
	}
	beforeCloseModuleDone := make(chan struct{})
	go func() {
		m.beforeCloseModule(ctx)
		close(beforeCloseModuleDone)
	}()

	select {
	case <-beforeCloseModuleDone:
	case <-ctx.Done():
	}

	m.closeAll(ctx)
	close(m.chanHasShutdown)
}
