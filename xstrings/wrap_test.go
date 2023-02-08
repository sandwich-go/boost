package xstrings

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWrap(t *testing.T) {
	Convey("wrap", t, func() {
		So(Wrap("a b\nc", 12), ShouldEqual, `a b
c`)
	})
}
