// Code generated by tools. DO NOT EDIT.
package xmath

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMath(t *testing.T) {
	Convey("Float64Equals", t, func() {
		So(Float64Equals(0.1, 0.2), ShouldBeFalse)
		So(Float64Equals(0.1, 0.1+EPSILON64), ShouldBeTrue)
	})

	Convey("Float32Equals", t, func() {
		So(Float32Equals(0.1, 0.2), ShouldBeFalse)
		So(Float32Equals(0.1, 0.1+EPSILON32), ShouldBeTrue)
	})

	Convey("IsZeroFloat64", t, func() {
		So(IsZeroFloat64(0.1), ShouldBeFalse)
		So(IsZeroFloat64(0), ShouldBeTrue)
	})

	Convey("IsZeroFloat32", t, func() {
		So(IsZeroFloat32(0.1), ShouldBeFalse)
		So(IsZeroFloat32(0), ShouldBeTrue)
	})

	Convey("IsBelowZeroFloat64", t, func() {
		So(IsBelowZeroFloat64(0.1), ShouldBeFalse)
		So(IsBelowZeroFloat64(0), ShouldBeTrue)
		So(IsBelowZeroFloat64(-0.1), ShouldBeTrue)
	})

	Convey("IsBelowZeroFloat32", t, func() {
		So(IsBelowZeroFloat32(0.1), ShouldBeFalse)
		So(IsBelowZeroFloat32(0), ShouldBeTrue)
		So(IsBelowZeroFloat32(-0.1), ShouldBeTrue)
	})

	Convey("Float32", t, func() {
		So(MaxFloat32(3.000000001, 3.000000002), ShouldEqual, 3.000000002)
		So(MinFloat32(3.000000001, 3.000000002), ShouldEqual, 3.000000001)
		So(AbsFloat32(3.000000001), ShouldEqual, 3.000000001)
		So(AbsFloat32(-3.000000001), ShouldEqual, 3.000000001)
		So(EffectZeroLimitFloat32(3.000000001, 0), ShouldEqual, 3.000000001)
		So(EffectZeroLimitFloat32(3.000000001, -3.000000002), ShouldEqual, 0)
	})

	Convey("Float64", t, func() {
		So(MaxFloat64(3.000000001, 3.000000002), ShouldEqual, 3.000000002)
		So(MinFloat64(3.000000001, 3.000000002), ShouldEqual, 3.000000001)
		So(AbsFloat64(3.000000001), ShouldEqual, 3.000000001)
		So(AbsFloat64(-3.000000001), ShouldEqual, 3.000000001)
		So(EffectZeroLimitFloat64(3.000000001, 0), ShouldEqual, 3.000000001)
		So(EffectZeroLimitFloat64(3.000000001, -3.000000002), ShouldEqual, 0)
	})

	Convey("Int", t, func() {
		So(MaxInt(3, 2), ShouldEqual, 3)
		So(MinInt(3, 2), ShouldEqual, 2)
		So(AbsInt(3), ShouldEqual, 3)
		So(AbsInt(-3), ShouldEqual, 3)
		So(EffectZeroLimitInt(3, 0), ShouldEqual, 3)
		So(EffectZeroLimitInt(3, -4), ShouldEqual, 0)
	})

	Convey("Int16", t, func() {
		So(MaxInt16(3, 2), ShouldEqual, 3)
		So(MinInt16(3, 2), ShouldEqual, 2)
		So(AbsInt16(3), ShouldEqual, 3)
		So(AbsInt16(-3), ShouldEqual, 3)
		So(EffectZeroLimitInt16(3, 0), ShouldEqual, 3)
		So(EffectZeroLimitInt16(3, -4), ShouldEqual, 0)
	})

	Convey("Int32", t, func() {
		So(MaxInt32(3, 2), ShouldEqual, 3)
		So(MinInt32(3, 2), ShouldEqual, 2)
		So(AbsInt32(3), ShouldEqual, 3)
		So(AbsInt32(-3), ShouldEqual, 3)
		So(EffectZeroLimitInt32(3, 0), ShouldEqual, 3)
		So(EffectZeroLimitInt32(3, -4), ShouldEqual, 0)
	})

	Convey("Int64", t, func() {
		So(MaxInt64(3, 2), ShouldEqual, 3)
		So(MinInt64(3, 2), ShouldEqual, 2)
		So(AbsInt64(3), ShouldEqual, 3)
		So(AbsInt64(-3), ShouldEqual, 3)
		So(EffectZeroLimitInt64(3, 0), ShouldEqual, 3)
		So(EffectZeroLimitInt64(3, -4), ShouldEqual, 0)
	})

	Convey("Int8", t, func() {
		So(MaxInt8(3, 2), ShouldEqual, 3)
		So(MinInt8(3, 2), ShouldEqual, 2)
		So(AbsInt8(3), ShouldEqual, 3)
		So(AbsInt8(-3), ShouldEqual, 3)
		So(EffectZeroLimitInt8(3, 0), ShouldEqual, 3)
		So(EffectZeroLimitInt8(3, -4), ShouldEqual, 0)
	})

	Convey("Uint", t, func() {
		So(MaxUint(3, 2), ShouldEqual, 3)
		So(MinUint(3, 2), ShouldEqual, 2)
	})

	Convey("Uint16", t, func() {
		So(MaxUint16(3, 2), ShouldEqual, 3)
		So(MinUint16(3, 2), ShouldEqual, 2)
	})

	Convey("Uint32", t, func() {
		So(MaxUint32(3, 2), ShouldEqual, 3)
		So(MinUint32(3, 2), ShouldEqual, 2)
	})

	Convey("Uint64", t, func() {
		So(MaxUint64(3, 2), ShouldEqual, 3)
		So(MinUint64(3, 2), ShouldEqual, 2)
	})

	Convey("Uint8", t, func() {
		So(MaxUint8(3, 2), ShouldEqual, 3)
		So(MinUint8(3, 2), ShouldEqual, 2)
	})
}