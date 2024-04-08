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
			t.Log("AfterFunc")
		})
		timer := dispatcher.AfterFuncInDomain(time.Second, func() {
			a++
			t.Log("AfterFuncInDomain")
		}, "goconvey")
		timer.t.Reset(0)
		dispatcher.RemoveAllTimerInDomain(DefaultTimerDomain)
		time.Sleep(500 * time.Millisecond)
		dispatcher.Close()
		dispatcher.Close()
		So(a, ShouldEqual, 1)

	})
}

func TestTimerResetDispatcher(t *testing.T) {
	Convey("normal timer", t, func() {
		a := 0
		var tm *time.Timer
		tm = time.AfterFunc(1*time.Second, func() {
			a++
			tm.Reset(1 * time.Second)
		})
		time.Sleep(2500 * time.Millisecond)
		tm.Stop()
		So(a, ShouldEqual, 2)
	})
	Convey("test reset AfterFunc", t, func() {
		dp := NewDispatcher(10)
		go func() {
			for {
				select {
				case tn := <-dp.TimerNotify():
					if tn != nil {
						tn.Cb()
					}
				}
			}
		}()
		a := 0
		var tt *DanglingTimer
		tt = dp.AfterFuncWithOwnershipTransferInDomain(time.Second, func() {
			a++
			t.Log("reset AfterFunc")
			tt.Reset(time.Second)
		}, "reset")
		time.Sleep(3500 * time.Millisecond)
		So(a, ShouldEqual, 3)
	})
}
