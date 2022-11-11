package xtime

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTimerDispatcher(t *testing.T) {

	Convey("First, create a TestTimerDispatcher", t, func() {
		dispatcher := NewDispatcher(10)
		a := 0
		dispatcher.AfterFunc(time.Second, func() {
			a++
		})
		timer := dispatcher.AfterFuncInDomain(time.Second, func() {
			a++
		}, "goconvey")
		timer.Reset(0)
		dispatcher.RemoveAllTimerInDomain(DefaultTimerDomain)
		time.Sleep(500 * time.Millisecond)
		dispatcher.Close()
		dispatcher.Close()
		So(a, ShouldEqual, 1)

	})

}
