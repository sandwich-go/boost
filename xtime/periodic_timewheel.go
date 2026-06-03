//go:build timewheel

package xtime

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/RussellLuo/timingwheel"
	"github.com/sandwich-go/boost/xmath"
)

// periodicTimewheel 在 timewheel build tag 下实现 PeriodicHandle。
//
// 实现选择说明：
// 上游 timingwheel 的 ScheduleFunc 能复用同一个 *timingwheel.Timer，但它在
// f() 之前就 reschedule 了下一轮，导致同一周期任务的 f 可能并发执行——
// 这与 std time 实现的「f 不重叠」语义冲突。为了让 Periodic 在两种 build
// tag 下行为一致（这是公开 API 的契约），timewheel 实现选择 *不* 用
// ScheduleFunc，而是按 std 实现一样在 f() 完成后再 reschedule。
//
// 代价：每轮 reschedule 仍调用 timingwheel.AfterFunc，会 alloc 新
// *timingwheel.Timer。本实现相比「直接调 xtime.AfterFunc 自手动 reschedule」
// 节省了 channel hop + closure capture 这一份；但拿不到 std 模式那种「整个
// 周期生命周期只 alloc 一个 timer」的全部收益。后续若要消除剩余 alloc，
// 可以在 timingwheel fork 里暴露公开的 reuse API。
type periodicTimewheel struct {
	d        time.Duration
	jitterPc int
	shutdown <-chan struct{}
	f        func()
	tick     func() // 复用的回调闭包，避免每轮新建

	stopped atomic.Bool
	mu      sync.Mutex             // 保护 timer 字段：fire 回调与 Stop 可能并发访问
	timer   *timingwheel.Timer
}

func newPeriodic(d time.Duration, jitterPercent int, shutdown <-chan struct{}, f func()) PeriodicHandle {
	p := &periodicTimewheel{
		d:        d,
		jitterPc: jitterPercent,
		shutdown: shutdown,
		f:        f,
	}
	// tick 闭包只创建一次，作为 timingwheel.AfterFunc 的回调反复使用——
	// 这是相对「每轮 closure 装箱」的关键节省点。
	p.tick = p.onFire
	p.schedule(p.nextDelay())
	return p
}

func (p *periodicTimewheel) onFire() {
	if p.shouldStop() {
		return
	}
	p.f()
	if p.shouldStop() {
		return
	}
	p.schedule(p.nextDelay())
}

func (p *periodicTimewheel) schedule(d time.Duration) {
	t := DefaultTiming.AfterFunc(d, p.tick)
	p.mu.Lock()
	if p.stopped.Load() {
		// 与 Stop() 竞争：Stop 已先执行；立即取消新 schedule 出来的 timer。
		p.mu.Unlock()
		t.Stop()
		return
	}
	p.timer = t
	p.mu.Unlock()
}

func (p *periodicTimewheel) nextDelay() time.Duration {
	if p.jitterPc <= 0 {
		return p.d
	}
	return time.Duration(xmath.Disturb(int64(p.d), int64(p.jitterPc)))
}

func (p *periodicTimewheel) shouldStop() bool {
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

func (p *periodicTimewheel) Stop() {
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
