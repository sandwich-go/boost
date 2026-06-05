package xmath

import (
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFloatHelpers(t *testing.T) {
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
}

func TestMaxMin(t *testing.T) {
	Convey("Max/Min on int families", t, func() {
		So(Max(3, 2), ShouldEqual, 3)
		So(Min(3, 2), ShouldEqual, 2)
		So(Max[int8](3, 2), ShouldEqual, int8(3))
		So(Min[int16](3, 2), ShouldEqual, int16(2))
		So(Max[int32](3, 2), ShouldEqual, int32(3))
		So(Max[int64](3, 2), ShouldEqual, int64(3))
		So(Max[uint](3, 2), ShouldEqual, uint(3))
		So(Max[uint8](3, 2), ShouldEqual, uint8(3))
		So(Max[uint16](3, 2), ShouldEqual, uint16(3))
		So(Max[uint32](3, 2), ShouldEqual, uint32(3))
		So(Max[uint64](3, 2), ShouldEqual, uint64(3))
	})

	Convey("Max/Min on string", t, func() {
		So(Max("b", "a"), ShouldEqual, "b")
		So(Min("b", "a"), ShouldEqual, "a")
	})

	Convey("Max/Min on float (Ordered: NaN 不传染)", t, func() {
		So(Max(3.000000001, 3.000000002), ShouldEqual, 3.000000002)
		So(Min(3.000000001, 3.000000002), ShouldEqual, 3.000000001)
		// NaN 通过 Max[Ordered]: NaN 与任何值比较都 false，所以走 else 分支返回 b
		// 这是 Go 内建 max() 行为，与 math.Max 的 NaN 传染不同
		So(math.IsNaN(Max(math.NaN(), 1.0)), ShouldBeFalse)
		So(Max(math.NaN(), 1.0), ShouldEqual, 1.0)
	})

	Convey("MaxFloat/MinFloat (IEEE 754: NaN 传染)", t, func() {
		So(MaxFloat[float32](3.000000001, 3.000000002), ShouldEqual, float32(3.000000002))
		So(MinFloat[float32](3.000000001, 3.000000002), ShouldEqual, float32(3.000000001))
		So(MaxFloat(3.000000001, 3.000000002), ShouldEqual, 3.000000002)
		So(MinFloat(3.000000001, 3.000000002), ShouldEqual, 3.000000001)
		So(math.IsNaN(MaxFloat(math.NaN(), 1.0)), ShouldBeTrue)
		So(math.IsNaN(MinFloat(math.NaN(), 1.0)), ShouldBeTrue)
	})
}

func TestAbs(t *testing.T) {
	Convey("Abs on signed integer (if v < 0)", t, func() {
		So(Abs(3), ShouldEqual, 3)
		So(Abs(-3), ShouldEqual, 3)
		So(Abs[int8](-3), ShouldEqual, int8(3))
		So(Abs[int16](-3), ShouldEqual, int16(3))
		So(Abs[int32](-3), ShouldEqual, int32(3))
		So(Abs[int64](-3), ShouldEqual, int64(3))
	})

	Convey("Abs on unsigned (always v itself)", t, func() {
		So(Abs[uint](3), ShouldEqual, uint(3))
		So(Abs[uint8](3), ShouldEqual, uint8(3))
	})

	Convey("AbsFloat 沿用 EPSILON 近似判负语义", t, func() {
		So(AbsFloat[float32](3.000000001), ShouldEqual, float32(3.000000001))
		So(AbsFloat[float32](-3.000000001), ShouldEqual, float32(3.000000001))
		So(AbsFloat(3.000000001), ShouldEqual, 3.000000001)
		So(AbsFloat(-3.000000001), ShouldEqual, 3.000000001)

		// EPSILON 近似判负的字面行为：v < EPSILON 时取反（含 v=0）
		// AbsFloat(0) -> IsBelowZeroFloat*(0) == true -> 返回 -0
		// math.Signbit 用于辨别 -0 vs +0
		So(math.Signbit(float64(AbsFloat[float64](0))), ShouldBeTrue)
		// v < EPSILON32 / 2 时也会被判作"负"取反
		var tiny float32 = EPSILON32 / 4 // 仍 < EPSILON32
		So(math.Signbit(float64(AbsFloat(tiny))), ShouldBeTrue)
	})
}

func TestEffectZeroLimit(t *testing.T) {
	Convey("EffectZeroLimit on signed integer", t, func() {
		So(EffectZeroLimit(3, 0), ShouldEqual, 3)
		So(EffectZeroLimit(3, -4), ShouldEqual, 0)
		So(EffectZeroLimit[int32](3, -4), ShouldEqual, int32(0))
	})

	Convey("EffectZeroLimitFloat 沿用 EPSILON 近似判负", t, func() {
		So(EffectZeroLimitFloat[float32](3.000000001, 0), ShouldEqual, float32(3.000000001))
		So(EffectZeroLimitFloat[float32](3.000000001, -3.000000002), ShouldEqual, float32(0))
		So(EffectZeroLimitFloat(3.000000001, 0), ShouldEqual, 3.000000001)
		So(EffectZeroLimitFloat(3.000000001, -3.000000002), ShouldEqual, 0.0)
	})
}
