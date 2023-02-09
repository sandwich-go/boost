package z

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
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
}
