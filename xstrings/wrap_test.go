package xstrings

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWrap(t *testing.T) {
	Convey("wrap", t, func() {
		So(Wrap("a b\nc", 12), ShouldEqual, `a b
c`)
	})
}
