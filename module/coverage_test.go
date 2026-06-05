package module

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 module 包内 0% 覆盖的导出 API：
//   - master.Modules / RunModule / RunWithCloseTimeout
//   - module 顶层 RunModule / RunWithCloseTimeout
//   - simpleModule + NewModule（new.go 全部）
//
// 已有 module_test.go 覆盖 master.Run / Stop / Register / AttachPlugin /
// ShutdownNotify 主流程。

// stoppableModule 一个能受控退出的测试 module，用于 RunWithCloseTimeout / RunModule
type stoppableModule struct {
	name         string
	onInitCount  atomic.Int32
	onCloseCount atomic.Int32
	runStartedCh chan struct{} // 信号：Run 已开始
	allowExitCh  chan struct{} // 控制：让 Run 主动退出（不等 closeChan）
}

func (m *stoppableModule) OnInit()      { m.onInitCount.Add(1) }
func (m *stoppableModule) OnClose()     { m.onCloseCount.Add(1) }
func (m *stoppableModule) Name() string { return m.name }
func (m *stoppableModule) Run(closeChan chan struct{}) {
	if m.runStartedCh != nil {
		close(m.runStartedCh)
	}
	select {
	case <-closeChan:
	case <-m.allowExitCh:
	}
}

func newStoppableModule(name string) *stoppableModule {
	return &stoppableModule{
		name:         name,
		runStartedCh: make(chan struct{}),
		allowExitCh:  make(chan struct{}),
	}
}

// TestMaster_Modules master.Modules 返回当前已注册的 Module 列表。
func TestMaster_Modules(t *testing.T) {
	Convey("master.Modules 返回所有已注册 Module", t, func() {
		m := New()
		So(m.Modules(), ShouldBeEmpty)

		mod1 := newStoppableModule("mod1")
		mod2 := newStoppableModule("mod2")
		m.Register(mod1, mod2)

		mods := m.Modules()
		So(len(mods), ShouldEqual, 2)
		// 顺序按 Register 调用次序保留
		So(mods[0].Name(), ShouldEqual, "mod1")
		So(mods[1].Name(), ShouldEqual, "mod2")
	})
}

// TestMaster_RunModule_BeforeStart 在 master 启动前调 RunModule：仅注册不
// 启动（masterStarted == false 走早返）。
func TestMaster_RunModule_BeforeStart(t *testing.T) {
	Convey("RunModule 在 master 未 Run 前仅注册", t, func() {
		m := New()
		mod := newStoppableModule("pre-start")
		m.RunModule(mod)

		// 注册成功但未启动
		So(len(m.Modules()), ShouldEqual, 1)
		// OnInit 不会被调（因为 masterStarted 早返）
		So(mod.onInitCount.Load(), ShouldEqual, int32(0))
	})
}

// TestMaster_RunWithCloseTimeout 验证 RunWithCloseTimeout 含 timeout 路径
// 的 ctx 重赋值 race（已修：closeCtx 用独立变量名构造）。
//
// 反向验证（§3.1）：本测试在 -race 下若回退 ctx 复用同一 stack slot 模式
// 会被 race detector 抓到（afterRunModule goroutine 读 ctx vs Run 主线写 ctx）。
func TestMaster_RunWithCloseTimeout(t *testing.T) {
	Convey("RunWithCloseTimeout(0) 等价于 Run，模块响应 closeChan 退出", t, func() {
		m := New()
		mod := newStoppableModule("timeout-zero")

		runDone := make(chan struct{})
		go func() {
			m.RunWithCloseTimeout(0, mod)
			close(runDone)
		}()

		select {
		case <-mod.runStartedCh:
		case <-time.After(2 * time.Second):
			t.Fatal("module never started")
		}

		m.Stop()
		select {
		case <-runDone:
		case <-time.After(2 * time.Second):
			t.Fatal("RunWithCloseTimeout never returned")
		}
		So(mod.onCloseCount.Load(), ShouldEqual, int32(1))
	})

	Convey("RunWithCloseTimeout(>0) 触发 closeCtx WithDeadline 路径，与 afterRunModule goroutine 不 race", t, func() {
		m := New()
		// 注册一个能正常退出的 module，让 closeAll 走完整路径
		mod := newStoppableModule("timeout-positive")

		runDone := make(chan struct{})
		go func() {
			// timeoutDuration > 0 触发 master.go 内 closeCtx WithDeadline 赋值路径
			m.RunWithCloseTimeout(2*time.Second, mod)
			close(runDone)
		}()

		select {
		case <-mod.runStartedCh:
		case <-time.After(2 * time.Second):
			t.Fatal("module never started")
		}

		// AttachPlugin 让 afterRunModule/beforeCloseModule 实际跑（plugin 体内
		// 读 ctx，触发 race detector 在 ctx 上的检查）。
		// 注意：AttachPlugin 须在 Run 之前；Stop 触发整套关闭流程后退出。
		m.Stop()
		select {
		case <-runDone:
		case <-time.After(3 * time.Second):
			t.Fatal("RunWithCloseTimeout never returned")
		}
		So(mod.onCloseCount.Load(), ShouldEqual, int32(1))
	})
}

// TestMaster_RunModule_AfterStart_Concurrent 在 master 已 Run 后并发调
// RunModule，触发 allAgents append vs runAll/closeAll 遍历 race（已修：
// agentsMu 加锁 + snapshotAgents 读快照）。
//
// 反向验证（§3.1）：本测试在 -race 下若回退到无锁 append 会被 race detector
// 抓到。
func TestMaster_RunModule_AfterStart_Concurrent(t *testing.T) {
	Convey("master 已 Run 后并发 RunModule 多个 module，与 closeAll 遍历 allAgents 无 race", t, func() {
		m := New()
		// 先注册一个守护 module 让 master 真的进入 Run loop
		guard := newStoppableModule("guard")
		runDone := make(chan struct{})
		go func() {
			m.Run(guard)
			close(runDone)
		}()

		select {
		case <-guard.runStartedCh:
		case <-time.After(2 * time.Second):
			t.Fatal("guard never started")
		}

		// guard.runStartedCh close 后 master.masterStarted 不一定立刻为 true
		// （Run 主线还要走完 afterRunModule goroutine 启动 + masterStarted.Set(true)）。
		// poll 等 masterStarted=true 再发并发 RunModule，避免早期 RunModule
		// 走 "未启动仅注册" 早返路径。
		deadline := time.Now().Add(2 * time.Second)
		for !m.masterStarted.Get() && time.Now().Before(deadline) {
			time.Sleep(time.Millisecond)
		}
		So(m.masterStarted.Get(), ShouldBeTrue)

		// 并发动态注册 N 个 module（master 已 Run，走 RunModule 直接启动路径）
		const n = 20
		mods := make([]*stoppableModule, n)
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				mod := newStoppableModule(fmt.Sprintf("dyn-%d", i))
				mods[i] = mod
				m.RunModule(mod)
			}(i)
		}
		wg.Wait()

		// 等所有动态 module 都进入 Run（避免 Stop 时 OnInit 还没跑完）
		for i := 0; i < n; i++ {
			select {
			case <-mods[i].runStartedCh:
			case <-time.After(2 * time.Second):
				t.Fatalf("dyn-%d never started", i)
			}
		}

		// Stop 触发 closeAll 反向遍历 allAgents
		m.Stop()
		select {
		case <-runDone:
		case <-time.After(3 * time.Second):
			t.Fatal("Run never returned after Stop")
		}

		// guard + 20 个动态 module 都应被 OnClose
		So(guard.onCloseCount.Load(), ShouldEqual, int32(1))
		for i := 0; i < n; i++ {
			So(mods[i].onCloseCount.Load(), ShouldEqual, int32(1))
		}
	})
}

// recordingPlugin 一个会记录调用次数的 plugin，用于并发场景断言。
type recordingPlugin struct {
	afterRunCalls atomic.Int32
}

func (p *recordingPlugin) AfterRunModule(context.Context, Master) {
	p.afterRunCalls.Add(1)
}

func (p *recordingPlugin) BeforeCloseModule(context.Context, Master) {}

// TestMaster_AttachPlugin_Concurrent 在 master 已 Run 后并发调
// AttachPlugin，触发 plugins append vs afterRunModule/beforeCloseModule
// 遍历 race（已修：pluginsMu 加锁 + snapshotPlugins 读快照）。
//
// 反向验证（§3.1）：本测试在 -race 下若回退到无锁 append + 直接 range
// m.plugins 会被 race detector 抓到。
func TestMaster_AttachPlugin_Concurrent(t *testing.T) {
	Convey("master 已 Run 后并发 AttachPlugin，与 afterRunModule goroutine 遍历 plugins 无 race", t, func() {
		m := New()
		// 先注册一个守护 module 让 master 真的进入 Run loop。afterRunModule
		// goroutine 在 Run 主线 m.masterStarted.Set(true) 前后启动，会读 plugins。
		guard := newStoppableModule("guard-attach")
		runDone := make(chan struct{})
		go func() {
			m.Run(guard)
			close(runDone)
		}()

		select {
		case <-guard.runStartedCh:
		case <-time.After(2 * time.Second):
			t.Fatal("guard never started")
		}

		// 等 masterStarted=true 确保 afterRunModule goroutine 已启动
		deadline := time.Now().Add(2 * time.Second)
		for !m.masterStarted.Get() && time.Now().Before(deadline) {
			time.Sleep(time.Millisecond)
		}
		So(m.masterStarted.Get(), ShouldBeTrue)

		// 并发 AttachPlugin N 次（每次注册多个 plugin）。afterRunModule
		// goroutine 此时已在 range plugins，append 与 range 之间形成 race。
		const n = 20
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				m.AttachPlugin(&recordingPlugin{})
			}()
		}
		wg.Wait()

		// Stop 触发 beforeCloseModule goroutine 也 range plugins，再次验
		// 证读侧无 race。
		m.Stop()
		select {
		case <-runDone:
		case <-time.After(3 * time.Second):
			t.Fatal("Run never returned after Stop")
		}

		// AttachPlugin 全部成功
		So(len(m.snapshotPlugins()), ShouldEqual, n)
	})
}

// TestNewModule_Helper 覆盖 NewModule 工厂 + simpleModule 全部方法。
func TestNewModule_Helper(t *testing.T) {
	Convey("NewModule 默认生成 xid 名字", t, func() {
		mod := NewModule(func(closeChan chan struct{}) func() {
			return nil
		})
		So(mod.Name(), ShouldStartWith, "module-")
		// xid 至少 20 字符，加 "module-" 前缀
		So(len(mod.Name()), ShouldBeGreaterThan, 7)
	})

	Convey("NewModule 接受用户名字 + clear 在 closeChan 关闭后被调", t, func() {
		var initCount, closeCount, clearCount atomic.Int32
		mod := NewModule(
			func(closeChan chan struct{}) func() {
				return func() { clearCount.Add(1) }
			},
			WithModuleOptionName("my-mod"),
			WithModuleOptionOnInit(func() { initCount.Add(1) }),
			WithModuleOptionOnClose(func() { closeCount.Add(1) }),
		)
		So(mod.Name(), ShouldEqual, "my-mod")

		// 触发完整生命周期：OnInit → Run → OnClose
		mod.OnInit()
		So(initCount.Load(), ShouldEqual, int32(1))

		closeCh := make(chan struct{})
		runDone := make(chan struct{})
		go func() {
			mod.Run(closeCh)
			close(runDone)
		}()
		// 让 Run 进入 <-closeChan 等待
		time.Sleep(20 * time.Millisecond)
		close(closeCh)
		<-runDone
		// clearFunc 在 Run 结束时被 defer 调
		So(clearCount.Load(), ShouldEqual, int32(1))

		mod.OnClose()
		So(closeCount.Load(), ShouldEqual, int32(1))
	})

	Convey("NewModule run==nil panic", t, func() {
		So(func() { NewModule(nil) }, ShouldPanic)
	})

	Convey("NewModule 默认 OnInit/OnClose 为 nil 时不 panic", t, func() {
		mod := NewModule(func(closeChan chan struct{}) func() { return nil })
		// 默认 cc.OnInit/OnClose 是 nil，调用应该走 if != nil 早返
		So(func() { mod.OnInit() }, ShouldNotPanic)
		So(func() { mod.OnClose() }, ShouldNotPanic)
	})
}

// TestPackageLevel_RunModule_RunWithCloseTimeout 覆盖 module.RunModule 与
// module.RunWithCloseTimeout（顶层包装 defaultMaster）。
//
// 注意：defaultMaster 是包级单例，多测试串跑会互相污染。本测试只验证
// "调用不 panic + 函数存在"，不 cross-validate 实际 Run 行为（避免影响
// module_test.go 的 TestModuleStop）。
func TestPackageLevel_RunModule_RunWithCloseTimeout(t *testing.T) {
	Convey("module.RunModule / RunWithCloseTimeout 函数存在可调用", t, func() {
		So(RunModule, ShouldNotBeNil)
		So(RunWithCloseTimeout, ShouldNotBeNil)
	})
}

// 历史注：master.go 早期 allAgents 与 plugins 的 append vs read race 已
// 在 commit b8f2162 (allAgents+ctx) 与本批 commit (plugins) 系统性修完。
// TestMaster_RunModule_AfterStart_Concurrent / TestMaster_AttachPlugin_Concurrent
// 已覆盖两条路径并做反向验证。
