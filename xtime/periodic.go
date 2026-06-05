package xtime

import "time"

// PeriodicHandle 周期任务句柄，用于取消尚未触发的下一轮回调。
//
// Stop 是幂等的；并发调用安全。
// Stop 不会等待已经开始执行的 f() 结束（与 time.AfterFunc / *time.Timer.Stop 一致）。
type PeriodicHandle interface {
	Stop()
}

// Periodic 每隔 d ± d*jitterPercent/100 调用一次 f()。
//
// 与 AfterFunc 自手动 reschedule 的区别：单个 Periodic 整个生命周期内只创建
// 一个底层 timer 对象，每轮通过 Reset / 内部 reschedule 复用，避免稳态产生
// per-tick 的 timer / closure 分配——大量长期运行的周期任务下能显著减轻 GC。
//
// 行为：
//   - 第一轮在 d ± jitter 之后触发，与 time.AfterFunc 行为一致；
//   - f() 在底层 timer 派发的 goroutine 上执行，同一 Periodic 的 f 不会重叠
//     （reschedule 在 f 返回后进行）；
//   - jitterPercent <= 0 时不抖动；jitterPercent > 0 时按 xmath.Disturb 计算
//     每轮真实间隔；
//   - Stop 之后下一轮 reschedule 被跳过，已经在执行的 f 不会被取消。
func Periodic(d time.Duration, jitterPercent int, f func()) PeriodicHandle {
	return PeriodicWithShutdown(d, jitterPercent, nil, f)
}

// PeriodicWithShutdown 与 Periodic 相同，但额外接受 shutdown channel；当
// channel 被关闭后等价于自动 Stop——下一轮 reschedule 不再发生。
//
// 与 Stop 的差别：shutdown 是只读监听，关闭后所有共享同一 channel 的
// Periodic 一起停止；不需要外层持有每个 handle。Stop 仍然可用作单独取消。
//
// shutdown 为 nil 时与 Periodic 等价。
func PeriodicWithShutdown(d time.Duration, jitterPercent int, shutdown <-chan struct{}, f func()) PeriodicHandle {
	if d <= 0 {
		panic("xtime.Periodic: duration must be positive")
	}
	if f == nil {
		panic("xtime.Periodic: f must not be nil")
	}
	return newPeriodic(d, jitterPercent, shutdown, f)
}

// newPeriodic 由 periodic_std.go 提供，使用 time.AfterFunc + *Timer.Reset。
