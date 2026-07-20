// Package xts 提供进程级全局的秒级时间戳工具：维护一份可初始化的全局时钟
// （时区偏移 + 时间提供者），并允许注入自定义的“当前时间”函数（例如接入带
// Tasklet 感知的调度框架时钟）。所有日期边界计算基于固定时区偏移，不涉及 DST。
//
// 时间提供者（TimeProvider）复用 github.com/sandwich-go/boost/xtime 的无状态实现。
package xts

import (
	"sync"
	"time"

	"github.com/sandwich-go/boost/xtime"
)

// Timestamp 是 Unix 秒级时间戳，日期边界计算基于全局配置的固定时区偏移。
type Timestamp int64

// Second 秒数，作为 Timestamp 各算术方法的单位类型。
type Second = int64

// TimeProvider 服务器时间提供者，别名自 xtime。
type TimeProvider = xtime.TimeProvider

// DefaultTimeProvider 与系统时钟等价的 TimeProvider，别名自 xtime。
type DefaultTimeProvider = xtime.DefaultTimeProvider

const (
	Minute = Second(60)
	Hour   = Second(3600)
	Day    = Second(86400)
	Week   = 7 * Day

	// SundayInMondayAsWeekBeginning 在“周一为一周起始”的语义下，把周日视为一周的第 7 天
	// （time.Saturday+1），使其排在周一~周六之后而非之前。
	SundayInMondayAsWeekBeginning = time.Saturday + 1
)

var (
	mu         sync.RWMutex
	provider   TimeProvider = DefaultTimeProvider{}
	zoneOffset int64        = 0
	loc                     = time.UTC
	// nowFunc 为“当前时间”的注入点，nil 表示回退到基于 provider + loc 的默认实现。
	nowFunc func() time.Time
)

// Option 配置全局时钟。
type Option func()

// WithTimeProvider 设置服务器时间提供者。
func WithTimeProvider(p TimeProvider) Option {
	return func() { provider = p }
}

// WithTimeZoneOffset 设置时区偏移，单位小时。
func WithTimeZoneOffset(offsetHours int64) Option {
	return func() { zoneOffset = offsetHours }
}

// WithNowFunc 注入“当前时间”函数。传 nil 恢复为基于 TimeProvider 的默认实现。
// 典型用途：接入带 Tasklet 感知的调度框架，使 Now 返回调度执行时间而非系统时间。
func WithNowFunc(f func() time.Time) Option {
	return func() { nowFunc = f }
}

// Initialize 应用配置并派生全局时区。未提供某项时保留其当前值（首次为默认值）。
func Initialize(opts ...Option) {
	mu.Lock()
	defer mu.Unlock()
	for _, o := range opts {
		o()
	}
	loc = time.FixedZone("UTC", int(zoneOffset)*3600)
	if provider == nil {
		panic("xts: TimeProvider is nil")
	}
}

// now 返回当前全局时钟时间：优先使用注入的 nowFunc，否则基于 TimeProvider + 时区计算。
func now() time.Time {
	mu.RLock()
	f, p, l := nowFunc, provider, loc
	mu.RUnlock()
	if f != nil {
		return f()
	}
	return xtime.ProviderNow(p, l)
}

// location 返回当前全局时区。
func location() *time.Location {
	mu.RLock()
	defer mu.RUnlock()
	return loc
}

// zoneOffsetHours 返回当前全局时区偏移，单位小时。
func zoneOffsetHours() int64 {
	mu.RLock()
	defer mu.RUnlock()
	return zoneOffset
}

// Location 返回当前配置的时区。
func Location() *time.Location { return location() }
