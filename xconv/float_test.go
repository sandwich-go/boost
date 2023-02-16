package xconv

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type testFloat struct {
	f1 float32
	f2 float64
}

func (b testFloat) Float32() float32 { return b.f1 }
func (b testFloat) Float64() float64 { return b.f2 }
func TestFloat32(t *testing.T) {
	Convey(`test float32`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is float32
		}{
			{nil, 0}, {float32(3.14), 3.14}, {3.14, 3.14}, {[]byte{0x00, 0x00, 0x48, 0x42}, 50},
			{testFloat{f1: 3.14, f2: 3.15}, 3.14}, {"3.14", 3.14},
		} {
			So(Float32(v.i), ShouldEqual, v.is)
		}
	})
}

func TestFloat64(t *testing.T) {
	Convey(`test float64`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is float64
		}{
			{nil, 0}, {3.14, 3.14}, {[]byte{24, 45, 68, 84, 251, 33, 9, 64}, 3.141592653589793},
			{testFloat{f1: 3.14, f2: 3.15}, 3.15}, {"3.14", 3.14},
		} {
			So(Float64(v.i), ShouldEqual, v.is)
		}
	})
}
