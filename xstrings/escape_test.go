package xstrings

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEscapeStringBackslash(t *testing.T) {
	Convey("escape string backslash", t, func() {
		So(EscapeStringBackslash("foo\x00bar"), ShouldEqual, "foo\\0bar")
		So(EscapeStringBackslash("foo\nbar"), ShouldEqual, "foo\\nbar")
		So(EscapeStringBackslash("foo\rbar"), ShouldEqual, "foo\\rbar")
		So(EscapeStringBackslash("foo\x1abar"), ShouldEqual, "foo\\Zbar")
		So(EscapeStringBackslash("foo\"bar"), ShouldEqual, "foo\\\"bar")
		So(EscapeStringBackslash("foo\\bar"), ShouldEqual, "foo\\\\bar")
		So(EscapeStringBackslash("foo'bar"), ShouldEqual, "foo\\'bar")
	})
}
