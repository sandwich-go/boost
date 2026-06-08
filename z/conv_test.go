package z

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConv(t *testing.T) {
	Convey("conv", t, func() {
		s := `a {{ .val1 }} {{ .val2 }}`
		s1 := StringToBytes(s)
		So(BytesToString(s1), ShouldEqual, s)
	})

	Convey("conv 空 string / nil byte slice 早返（commit 69dcc68 加的边界）", t, func() {
		// StringToBytes(""): len(s) == 0 → return nil
		So(StringToBytes(""), ShouldBeNil)

		// BytesToString(nil): len(b) == 0 → return ""
		So(BytesToString(nil), ShouldEqual, "")
		// BytesToString([]byte{}): 同上
		So(BytesToString([]byte{}), ShouldEqual, "")
	})
}
