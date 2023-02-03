package xcmd

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCmdBoolean(t *testing.T) {
	Convey("cmd boolean", t, func() {
		for _, cmd := range []struct {
			s       string
			isFalse bool
		}{
			{s: "", isFalse: true},
			{s: "0", isFalse: true},
			{s: "n", isFalse: true}, {s: "N", isFalse: true},
			{s: "no", isFalse: true}, {s: "No", isFalse: true}, {s: "nO", isFalse: true}, {s: "NO", isFalse: true},
			{s: "off", isFalse: true}, {s: "Off", isFalse: true}, {s: "OFf", isFalse: true}, {s: "OFF", isFalse: true},
			{s: "false", isFalse: true}, {s: "False", isFalse: true}, {s: "faLSE", isFalse: true}, {s: "FALSE", isFalse: true},
			{s: "1", isFalse: false}, {s: "Y", isFalse: false}, {s: "Yes", isFalse: false}, {s: "YES", isFalse: false}, {s: "ON", isFalse: false}, {s: "TRUE", isFalse: false},
		} {
			So(IsFalse(cmd.s), ShouldEqual, cmd.isFalse)
			So(IsTrue(cmd.s), ShouldEqual, !cmd.isFalse)
		}
	})
}
