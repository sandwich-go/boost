//go:build timingwheel

package xtime

import (
	"time"

	"github.com/RussellLuo/timingwheel"
)

// var timeAfterFunc = time.AfterFunc
var (
	DefaultTiming = timingwheel.NewTimingWheel(time.Millisecond*100, 128)
	timeAfterFunc = func(d time.Duration, f func()) internalTimer {
		t := DefaultTiming.AfterFunc(d, f)
		return &timer{t: t, fn: f}
	}
)

func init() {
	DefaultTiming.Start()
}

type timer struct {
	t  *timingwheel.Timer
	fn func()
}

func (t *timer) Stop() bool {
	return t.t.Stop()
}

func (t *timer) Reset(duration time.Duration) bool {
	active := true
	if stopped := t.t.Stop(); !stopped {
		active = false
	}
	t2 := DefaultTiming.AfterFunc(duration, t.fn)
	t.t = t2
	return active
}
