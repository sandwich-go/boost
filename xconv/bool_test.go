package xconv

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type testBool struct{ b bool }

func (b testBool) Bool() bool { return b.b }

type testStruct struct{}

func TestBool(t *testing.T) {
	Convey(`test bool`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is bool
		}{
			{nil, false}, {false, false}, {true, true},
			{1, true}, {0, false}, {int8(1), true}, {int8(0), false},
			{int16(1), true}, {int16(0), false}, {int32(1), true}, {int32(0), false}, {int64(1), true}, {int64(0), false},
			{uint(1), true}, {uint(0), false}, {uint8(1), true}, {uint8(0), false},
			{uint16(1), true}, {uint16(0), false}, {uint32(1), true}, {uint32(0), false}, {uint64(1), true}, {uint64(0), false},
			{testBool{}, false}, {testBool{true}, true},
			{&testStruct{}, true}, {[]int{}, false}, {map[int]int{}, false}, {[]int{1}, true}, {map[int]int{1: 1}, true},
			{testStruct{}, true}, {1.0, true},
		} {
			if v.is {
				So(Bool(v.i), ShouldBeTrue)
			} else {
				So(Bool(v.i), ShouldBeFalse)
			}
		}
	})
}
