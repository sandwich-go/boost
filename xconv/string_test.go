package xconv

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type testString struct {
	f1 string
}

func (b testString) String() string { return b.f1 }

func TestString(t *testing.T) {
	Convey(`test string`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is string
		}{
			{i: nil}, {"3", "3"}, {[]byte("3"), "3"}, {3, "3"}, {uint(3), "3"},
			{int8(3), "3"}, {int16(3), "3"}, {int32(3), "3"}, {int64(3), "3"},
			{uint8(3), "3"}, {uint16(3), "3"}, {uint32(3), "3"}, {uint64(3), "3"},
			{float32(3), "3"}, {float64(3), "3"}, {true, "true"}, {false, "false"},
			{time.Time{}, ""}, {(*time.Time)(nil), ""},
			{testString{"3"}, "3"}, {&testString{"3"}, "3"},
			{map[string]string{"3": "3"}, `{"3":"3"}`},
		} {
			So(String(v.i), ShouldEqual, v.is)
		}
		So(String(func() {}), ShouldStartWith, "0x")
	})
}
