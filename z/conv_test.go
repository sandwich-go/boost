package z

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConv(t *testing.T) {
	Convey("conv", t, func() {
		s := `a {{ .val1 }} {{ .val2 }}`
		s1 := StringToBytes(s)
		So(BytesToString(s1), ShouldEqual, s)
	})
}
