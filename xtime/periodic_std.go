package xtime

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/sandwich-go/boost/xmath"
)

// periodicStd 使用 std 库 time.AfterFunc 创建一个 *time.Timer，
// 在 timer 的回调里执行 f() 后调用 t.Reset(...) 复用同一个 timer，
// 避免每轮新建 *time.Timer 触发 GC。
//
// 关于 time.Timer.Reset 的并发安全：
//   - time.AfterFunc 风格 timer（不是 NewTimer + chan 风格）的回调内调用
//     Reset 是 Go runtime 显式支持的用法；
//   - 我们只在 f() 完成 *之后* 调用 Reset，且回调本身就跑在 timer 派发的
//     goroutine 上，没有 channel 竞争。
//
// timer 字段并发访问：
//   - 构造时主调 goroutine 写 p.timer = time.AfterFunc(...)；
//   - 回调 goroutine 读 p.timer 调 Reset；
//   - Stop goroutine 读 p.timer 调 Stop。
//   构造与回调之间没有 happens-before 关系（time.AfterFunc 内部调度
//   race detector 看不到），所以必须显式 mu 保护写/读。
type periodicStd struct {
	d        time.Duration
	jitterPc int
	shutdown <-chan struct{}
	f        func()

	stopped atomic.Bool
	mu      sync.Mutex // 保护 timer 字段的写/读
	timer   *time.Timer
}

func newPeriodic(d time.Duration, jitterPercent int, shutdown <-chan struct{}, f func()) PeriodicHandle {
	p := &periodicStd{
		d:        d,
		jitterPc: jitterPercent,
		shutdown: shutdown,
		f:        f,
	}
	// 先在 mu 里 schedule，避免 onFire 在 timer 字段尚未发布前就读到。
	p.mu.Lock()
	p.timer = time.AfterFunc(p.nextDelay(), p.onFire)
	p.mu.Unlock()
	return p
}

func (p *periodicStd) onFire() {
	if p.shouldStop() {
		return
	}
	p.f()
	if p.shouldStop() {
		return
	}
	// Reset 复用同一个 *time.Timer，不产生新分配。
	// 与外部 Stop 的竞争：若 Stop 先执行，stopped 已 true，上面 shouldStop
	// 已 return；若 Stop 在 Reset 之后执行，timer.Stop() 会取消已 reschedule
	// 的下一轮。
	p.mu.Lock()
	t := p.timer
	p.mu.Unlock()
	t.Reset(p.nextDelay())
}

func (p *periodicStd) nextDelay() time.Duration {
	if p.jitterPc <= 0 {
		return p.d
	}
	return time.Duration(xmath.Disturb(int64(p.d), int64(p.jitterPc)))
}

func (p *periodicStd) shouldStop() bool {
	if p.stopped.Load() {
		return true
	}
	if p.shutdown != nil {
		select {
		case <-p.shutdown:
			return true
		default:
		}
	}
	return false
}

func (p *periodicStd) Stop() {
	if !p.stopped.CompareAndSwap(false, true) {
		return
	}
	p.mu.Lock()
	t := p.timer
	p.mu.Unlock()
	if t != nil {
		t.Stop()
	}
}
