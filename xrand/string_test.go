package xrand

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func TestRandString(t *testing.T) {
	Convey("rand string", t, func() {
		ss := String(10)
		t.Log(ss)
		So(len(ss), ShouldEqual, 10)

		t0 := time.Unix(1000000, 0)
		nowFunc = func() time.Time {
			return t0
		}
		ss = StringWithTimestamp(10)
		So(strings.HasSuffix(ss, "1000000"), ShouldBeTrue)
	})
}
