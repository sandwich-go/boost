package xtime

import (
	"time"
)

var globalCop *Cop = NewCop(func() time.Time { return time.Now() })

// Stop 停止globalCop的时间，使用系统时间
func Stop() { globalCop.Stop() }

// SetNowProvider 设定now提供方法
func SetNowProvider(nowProvider func() time.Time) { globalCop.nowProvider = nowProvider }

// Now 获取当前时间，精度秒
var Now = func() time.Time { return globalCop.Now() }

// Unix 获取当前时间戳，精度秒
// time consume: time.Now().Unix() > atomic.LoadInt64(&_atomic_timestamp) > globalCop.Unix() > globalCop.ts
var Unix = func() int64 { return globalCop.Unix() }

// UnixMilli 当前时间，单位毫秒
var UnixMilli = func() int64 { return globalCop.UnixMilli() }

const mockDeprecated = "Time Mock Deprecated: use https://github.com/sandwich-go/xtime instead"

// mock属于测试类需求，stime对时间mock的支持仅限于Freeze，Scale,Travel，对于mock下的tick，after没有做完备支持，迁移到https://github.com/sandwich-go/xtime独立扩展支持
// Mocked Deprecated: https://github.com/sandwich-go/xtime instead
var Mocked = func() bool { panic(mockDeprecated) }

// Freeze Deprecated: https://github.com/sandwich-go/xtime instead
var Freeze = func(t time.Time) { panic(mockDeprecated) }

// Scale Deprecated: https://github.com/sandwich-go/xtime instead
var Scale = func(scale float64) { panic(mockDeprecated) }

// Travel Deprecated: https://github.com/sandwich-go/xtime instead
var Travel = func(t time.Time) { panic(mockDeprecated) }

// Since Deprecated: https://github.com/sandwich-go/xtime instead
var Since = func(t time.Time) time.Duration { panic(mockDeprecated) }

// Sleep Deprecated: https://github.com/sandwich-go/xtime instead
var Sleep = func(d time.Duration) { panic(mockDeprecated) }

// After Deprecated: https://github.com/sandwich-go/xtime instead
var After = func(d time.Duration) <-chan time.Time { panic(mockDeprecated) }

// Tick Deprecated: https://github.com/sandwich-go/xtime instead
var Tick = func(d time.Duration) <-chan time.Time { panic(mockDeprecated) }

// Return Deprecated: https://github.com/sandwich-go/xtime instead
var Return = func() { panic(mockDeprecated) }
