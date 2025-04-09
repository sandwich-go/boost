//go:build timingwheel

package xtime

import (
	"time"

	"github.com/RussellLuo/timingwheel"
)

// var timeAfterFunc = time.AfterFunc
var (
	DefaultTiming = timingwheel.NewTimingWheel(time.Millisecond*100, 128)
	timeAfterFunc = DefaultTiming.AfterFunc
)

type timer *timingwheel.Timer

func init() {
	DefaultTiming.Start()
}
