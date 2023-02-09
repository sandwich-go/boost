package xstrings

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCamelCase(t *testing.T) {
	Convey("camel case", t, func() {
		So(IsASCIILower('b'), ShouldBeTrue)
		So(IsASCIILower('9'), ShouldBeFalse)
		So(IsASCIIDigit('b'), ShouldBeFalse)
		So(IsASCIIDigit('9'), ShouldBeTrue)
		So(CamelCase("__camel__case_"), ShouldEqual, "XCamel_Case_")
		So(CamelCase("camel_case"), ShouldEqual, "CamelCase")
	})
}
