package xtime

import (
	"math"
	"time"
)

// TimeProvider 描述一个可缩放的虚拟时钟：以 ActualBase 为现实基准点、VirtualBase
// 为虚拟基准点，按 Scale 倍率映射现实流逝的时间。Scale=1 且两 Base 相等即等价于系统时钟。
// 用于压测加速、回放等需要偏移或加速时间的场景。
type TimeProvider interface {
	// VirtualBase 虚拟世界的基准时刻
	VirtualBase() time.Time
	// ActualBase 现实世界的基准时刻
	ActualBase() time.Time
	// Scale 时间流速倍率，>1 加速、<1 减速
	Scale() float64
}

// timeProviderBase 是 DefaultTimeProvider 的虚拟/现实公共基准点。
// 取一个固定历史时刻只是为了让 VirtualBase 与 ActualBase 一致（Scale=1 时等价系统时钟），
// 具体取值不影响 DefaultTimeProvider 的行为。
var timeProviderBase = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

// DefaultTimeProvider 是与系统时钟等价的 TimeProvider：Scale=1 且虚拟基准与现实基准相同。
type DefaultTimeProvider struct{}

func (d DefaultTimeProvider) VirtualBase() time.Time { return timeProviderBase }

func (d DefaultTimeProvider) ActualBase() time.Time { return timeProviderBase }

func (d DefaultTimeProvider) Scale() float64 { return 1 }

// ProviderNow 按 TimeProvider 的虚拟基准 + 缩放后的现实流逝时间计算当前虚拟时刻，
// 并转换到 loc 时区返回。
func ProviderNow(p TimeProvider, loc *time.Location) time.Time {
	past := time.Duration(float64(time.Since(p.ActualBase())) / p.Scale())
	return p.VirtualBase().Add(past).In(loc)
}

// ScaleDuration 把现实时长 d 按 TimeProvider 的 Scale 换算成虚拟世界时长。
func ScaleDuration(p TimeProvider, d time.Duration) time.Duration {
	return time.Duration(math.Round(float64(d) * p.Scale()))
}
