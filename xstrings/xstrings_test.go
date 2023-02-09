package xstrings

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString(t *testing.T) {
	Convey("string", t, func() {
		So(FirstUpper("abc"), ShouldEqual, `Abc`)
		So(FirstUpper("_abc"), ShouldEqual, `_abc`)
		So(FirstLower("Abc"), ShouldEqual, `abc`)
		So(FirstLower("_abc"), ShouldEqual, `_abc`)
		So(HasPrefixIgnoreCase("abc", "AB"), ShouldBeTrue)
		So(HasPrefixIgnoreCase("abc", "ac"), ShouldBeFalse)
		So(TrimPrefixIgnoreCase("abc", "AB"), ShouldEqual, "c")
	})
}
