package xconv

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUint(t *testing.T) {
	Convey(`test uint`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is uint
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 4},
		} {
			So(Uint(v.i), ShouldEqual, v.is)
		}
	})

	Convey(`test uint8`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is uint8
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 4},
		} {
			So(Uint8(v.i), ShouldEqual, v.is)
		}
	})

	Convey(`test uint16`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is uint16
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 4},
		} {
			So(Uint16(v.i), ShouldEqual, v.is)
		}
	})

	Convey(`test uint32`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is int32
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 4},
		} {
			So(Uint32(v.i), ShouldEqual, v.is)
		}
	})

	Convey(`test uint64`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is uint64
		}{
			{nil, 0}, {3, 3}, {uint(3), 3}, {"3", 3}, {float32(3), 3}, {float64(3), 3},
			{int8(3), 3}, {int16(3), 3}, {int32(3), 3}, {int64(3), 3},
			{uint8(3), 3}, {uint16(3), 3}, {uint32(3), 3}, {uint64(3), 3},
			{testInt{f1: 3, f2: 4}, 4},
		} {
			So(Uint64(v.i), ShouldEqual, v.is)
		}
	})
}
