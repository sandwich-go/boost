package xsync

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWaitGroupWithoutTimeout(t *testing.T) {
	var wg WaitGroupTimeout
	wg.Add(1)
	go func() {
		time.Sleep(time.Second)
		wg.Done()
	}()
	t0 := time.Now()
	wg.Wait()
	t.Logf("wait for %s", time.Since(t0))
}

func TestWaitGroupWithTimeout(t *testing.T) {
	Convey("wait group with timeout", t, func() {
		var wg WaitGroupTimeout
		wg.Add(1)
		go func() {
			time.Sleep(time.Second * 10)
			wg.Done()
		}()
		t0 := time.Now()
		tm := wg.WaitTimeout(time.Second)
		So(tm, ShouldBeTrue)
		t.Logf("wait for %s", time.Since(t0))
	})
}
