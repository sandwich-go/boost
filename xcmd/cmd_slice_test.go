package xcmd

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCmdSlice(t *testing.T) {
	Convey("cmd slice", t, func() {
		for _, cmd := range []struct {
			s        string
			expected string
		}{
			{s: "", expected: ""},
			{s: "hi, hello world", expected: "hi hello world"},
		} {
			So(strings.Join(Slice(cmd.s), ""), ShouldEqual, cmd.expected)
		}
	})
}
