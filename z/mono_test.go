package z

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMono(t *testing.T) {
	Convey("mono", t, func() {
		n := MonoOffset()
		So(n, ShouldNotBeZeroValue)
		BusyDelay(1 * time.Second)
		e := MonoSince(n)
		So(e, ShouldBeGreaterThan, 1*time.Second)
		t.Log(Now())
	})

	Convey("NowWithOffset 用给定 monoOffset 计算 wall clock", t, func() {
		// NowWithOffset 把 wallClock.t + (monoOffset - wc.offset) 算出来
		// 调用方传入当前 mono offset 应得到接近 time.Now() 的 wall time
		current := time.Duration(MonoOffset())
		got := NowWithOffset(current)
		// 应在当前 wall time 附近 1 秒内
		So(got.Sub(time.Now()).Abs(), ShouldBeLessThan, 1*time.Second)
	})
}
