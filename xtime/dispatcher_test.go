package xtime

import (
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTimerDispatcher(t *testing.T) {

	Convey("First, create a TestTimerDispatcher", t, func() {
		dispatcher := NewDispatcher(10)
		// a 用 atomic.Int32：timer callback 在 dispatcher goroutine 跑，与主
		// 测试 goroutine 的 So(a, ShouldEqual, 1) 读形成 race。
		var a atomic.Int32
		dispatcher.AfterFunc(time.Second, func() {
			a.Add(1)
			t.Log("AfterFunc")
		})
		timer := dispatcher.AfterFuncInDomain(time.Second, func() {
			a.Add(1)
			t.Log("AfterFuncInDomain")
		}, "goconvey")
		timer.t.Reset(0)
		dispatcher.RemoveAllTimerInDomain(DefaultTimerDomain)
		time.Sleep(500 * time.Millisecond)
		dispatcher.Close()
		dispatcher.Close()
		So(a.Load(), ShouldEqual, int32(1))

	})
}

func TestTimerResetDispatcher(t *testing.T) {
	Convey("normal timer", t, func() {
		// a / tm 用 atomic 包裹：time.AfterFunc 的 callback 在内部 goroutine
		// 跑，与主测试 goroutine 在 tm.Stop() 后读 a 形成 race（Stop 不保证
		// 已发起的 callback 已完成）。
		var a atomic.Int32
		var tm atomic.Pointer[time.Timer]
		tm.Store(time.AfterFunc(1*time.Second, func() {
			a.Add(1)
			tm.Load().Reset(1 * time.Second)
		}))
		time.Sleep(2500 * time.Millisecond)
		tm.Load().Stop()
		// 给 callback 一个 happens-before 时机让它完成（Stop 后 50ms 等收尾）
		time.Sleep(50 * time.Millisecond)
		So(a.Load(), ShouldEqual, int32(2))
	})
	Convey("test reset AfterFunc", t, func() {
		dp := NewDispatcher(10)
		// notify goroutine 用 stopCh 终止，不用 range —— Close 内部自己也
		// 会 drain ChanTimer，两边 range 同一 closed chan 会触发 race。
		stopNotify := make(chan struct{})
		notifyDone := make(chan struct{})
		go func() {
			defer close(notifyDone)
			for {
				select {
				case <-stopNotify:
					return
				case tn := <-dp.TimerNotify():
					if tn != nil {
						tn.Cb()
					}
				}
			}
		}()
		// a / tt 同上，DanglingTimer callback 跑在 dispatcher goroutine。
		var a atomic.Int32
		var tt atomic.Pointer[DanglingTimer]
		tt.Store(dp.AfterFuncWithOwnershipTransferInDomain(time.Second, func() {
			a.Add(1)
			t.Log("reset AfterFunc")
			tt.Load().Reset(time.Second)
		}, "reset"))
		time.Sleep(3500 * time.Millisecond)
		So(a.Load(), ShouldEqual, int32(3))
		// 收尾：停 timer 阻止后续 callback；停 notify goroutine；close
		// dispatcher。顺序：stopNotify 先 → notify goroutine 不再消费 →
		// tt.Stop 阻止 timer 再触发 cb（避免拿不到 notify 而 hang）→ Close
		// dispatcher。
		tt.Load().Stop()
		close(stopNotify)
		<-notifyDone
		dp.Close()
	})
}
