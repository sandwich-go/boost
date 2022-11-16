package xz

import (
	"time"
	_ "unsafe"
)

// runtimeNanotime 返回从某一个节点开始(如系统启动)的单调时间
//go:noescape
func runtimeNanotime() int64

// MonoTimeDuration 防止MonoSince时参数传递错误
type MonoTimeDuration time.Duration

//go:linkname runtimeNanotime runtime.nanotime
func MonoOffset() MonoTimeDuration {
	return MonoTimeDuration(runtimeNanotime())
}

// MonoSince returns the time elapsed since t, obtained previously using Now.
func MonoSince(t MonoTimeDuration) time.Duration {
	return time.Duration(MonoOffset() - t)
}

// BusyDelay waits for given duration using busy waiting
func BusyDelay(duration time.Duration) {
	end := runtimeNanotime() + int64(duration)
	for runtimeNanotime() < end {
	}
}
