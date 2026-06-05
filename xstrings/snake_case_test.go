package xstrings

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSnakeCase(t *testing.T) {
	Convey("snake case", t, func() {
		So(SnakeCase("XCamel_Case_"), ShouldEqual, "x_camel__case_")
		So(SnakeCase("CamelCase"), ShouldEqual, "camel_case")
	})
}
