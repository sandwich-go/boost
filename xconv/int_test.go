package xconv

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type testInt struct {
	f1 int64
	f2 uint64
}

func (b testInt) Int64() int64   { return b.f1 }
func (b testInt) Uint64() uint64 { return b.f2 }

func TestInt(t *testing.T) {
	Convey(`test int`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is int
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 3},
		} {
			So(Int(v.i), ShouldEqual, v.is)
		}
	})

	Convey(`test int8`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is int8
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 3},
		} {
			So(Int8(v.i), ShouldEqual, v.is)
		}
	})

	Convey(`test int16`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is int16
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 3},
		} {
			So(Int16(v.i), ShouldEqual, v.is)
		}
	})

	Convey(`test int32`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is int32
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 3},
		} {
			So(Int32(v.i), ShouldEqual, v.is)
		}
	})

	Convey(`test int64`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is int64
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 3},
		} {
			So(Int64(v.i), ShouldEqual, v.is)
		}
	})
}
