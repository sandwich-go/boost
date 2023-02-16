package xconv

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNum(t *testing.T) {
	Convey(`test parse num`, t, func() {
		for _, v := range []struct {
			i      string
			is     interface{}
			hasErr bool
		}{
			{i: "3", is: 3}, {i: "3x", is: 3, hasErr: true},
			{i: "3", is: uint(3)}, {i: "3x", is: uint(3), hasErr: true},
			{i: "3", is: int8(3)}, {i: "3x", is: int8(3), hasErr: true},
			{i: "3", is: int16(3)}, {i: "3x", is: int16(3), hasErr: true},
			{i: "3", is: int32(3)}, {i: "3x", is: int32(3), hasErr: true},
			{i: "3", is: int64(3)}, {i: "3x", is: int64(3), hasErr: true},
			{i: "3", is: uint8(3)}, {i: "3x", is: uint8(3), hasErr: true},
			{i: "3", is: uint16(3)}, {i: "3x", is: uint16(3), hasErr: true},
			{i: "3", is: uint32(3)}, {i: "3x", is: uint32(3), hasErr: true},
			{i: "3", is: uint64(3)}, {i: "3x", is: uint64(3), hasErr: true},
		} {
			var check = func(val, expected interface{}, err error, hasErr bool) {
				if hasErr {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(val, ShouldEqual, expected)
				}
			}
			switch vv := v.is.(type) {
			case int:
				val, err := ParseInt(v.i)
				check(val, vv, err, v.hasErr)
			case uint:
				val, err := ParseUint(v.i)
				check(val, vv, err, v.hasErr)
			case int8:
				val, err := ParseInt8(v.i)
				check(val, vv, err, v.hasErr)
			case int16:
				val, err := ParseInt16(v.i)
				check(val, vv, err, v.hasErr)
			case int32:
				val, err := ParseInt32(v.i)
				check(val, vv, err, v.hasErr)
			case int64:
				val, err := ParseInt64(v.i)
				check(val, vv, err, v.hasErr)
			case uint8:
				val, err := ParseUint8(v.i)
				check(val, vv, err, v.hasErr)
			case uint16:
				val, err := ParseUint16(v.i)
				check(val, vv, err, v.hasErr)
			case uint32:
				val, err := ParseUint32(v.i)
				check(val, vv, err, v.hasErr)
			case uint64:
				val, err := ParseUint64(v.i)
				check(val, vv, err, v.hasErr)
			}
		}
	})

	Convey(`test format num`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is string
		}{
			{3, "3"}, {uint(3), "3"},
			{int8(3), "3"}, {int16(3), "3"}, {int32(3), "3"}, {int64(3), "3"},
			{uint8(3), "3"}, {uint16(3), "3"}, {uint32(3), "3"}, {uint64(3), "3"},
		} {
			switch vv := v.i.(type) {
			case int:
				So(FormatInt(vv), ShouldEqual, v.is)
			case uint:
				So(FormatUint(vv), ShouldEqual, v.is)
			case int8:
				So(FormatInt8(vv), ShouldEqual, v.is)
			case int16:
				So(FormatInt16(vv), ShouldEqual, v.is)
			case int32:
				So(FormatInt32(vv), ShouldEqual, v.is)
			case int64:
				So(FormatInt64(vv), ShouldEqual, v.is)
			case uint8:
				So(FormatUint8(vv), ShouldEqual, v.is)
			case uint16:
				So(FormatUint16(vv), ShouldEqual, v.is)
			case uint32:
				So(FormatUint32(vv), ShouldEqual, v.is)
			case uint64:
				So(FormatUint64(vv), ShouldEqual, v.is)
			default:
				t.Fatal("unsupported type")
			}
		}
	})

	Convey(`test format num to slice`, t, func() {
		for _, v := range []struct {
			i  interface{}
			is []byte
		}{
			{int(3), []byte{51}}, {int8(3), []byte{51}}, {int16(3), []byte{51}}, {int32(3), []byte{51}}, {int64(3), []byte{51}},
			{uint(3), []byte{51}}, {uint8(3), []byte{51}}, {uint16(3), []byte{51}}, {uint32(3), []byte{51}}, {uint64(3), []byte{51}},
		} {
			switch vv := v.i.(type) {
			case int:
				So(FormatIntToSlice(vv), ShouldResemble, v.is)
			case uint:
				So(FormatUintToSlice(vv), ShouldResemble, v.is)
			case int8:
				So(FormatInt8ToSlice(vv), ShouldResemble, v.is)
			case int16:
				So(FormatInt16ToSlice(vv), ShouldResemble, v.is)
			case int32:
				So(FormatInt32ToSlice(vv), ShouldResemble, v.is)
			case int64:
				So(FormatInt64ToSlice(vv), ShouldResemble, v.is)
			case uint8:
				So(FormatUint8ToSlice(vv), ShouldResemble, v.is)
			case uint16:
				So(FormatUint16ToSlice(vv), ShouldResemble, v.is)
			case uint32:
				So(FormatUint32ToSlice(vv), ShouldResemble, v.is)
			case uint64:
				So(FormatUint64ToSlice(vv), ShouldResemble, v.is)
			default:
				t.Fatal("unsupported type")
			}
		}
	})
}
