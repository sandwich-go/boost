package graph

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRectangle(t *testing.T) {
	HelixRectRangeFromCenterAndMargin[int](P(2, 2), 2, func(p Point[int]) bool {
		t.Log(p.X(), p.Y())
		return true
	})

	Convey("test intersection with line", t, func() {
		x := P(1, 0)
		y := P(1, 4)
		z := x.Add(y)
		So(z.X(), ShouldEqual, 2)
		So(z.Y(), ShouldEqual, 4)
		So(z.String(), ShouldEqual, "(2,4)")
		So(z.Sub(x).Equals(y), ShouldBeTrue)
		q := y.Mul(3)
		So(q.X(), ShouldEqual, y.X()*3)
		So(q.Y(), ShouldEqual, y.Y()*3)
		So(q.Div(3).Equals(y), ShouldBeTrue)
		var r1, r2, r3, r4 int = 1, 0, 5, 5
		r := Rect[int](r1, r2, r3, r4)
		So(x.In(r), ShouldBeTrue)
		So(r.String(), ShouldEqual, fmt.Sprintf("R{(%d,%d)-(%d,%d)}", r1, r2, r3, r4))
		var minX, maxX, minY, maxY int
		var reset = func() { minX, maxX, minY, maxY = 5, 0, 5, 0 }
		var f = func(p Point[int]) bool {
			if p.X() < minX {
				minX = p.X()
			}
			if p.Y() < minY {
				minY = p.Y()
			}
			if p.X() > maxX {
				maxX = p.X()
			}
			if p.Y() > maxY {
				maxY = p.Y()
			}
			return true
		}
		reset()
		r.RangePointsMinClosedMaxOpen(f)
		So(minX, ShouldEqual, r1)
		So(minY, ShouldEqual, r2)
		So(maxX, ShouldEqual, r3-1)
		So(maxY, ShouldEqual, r4-1)
		reset()
		r.RangePointsMinOpenMaxClosed(f)
		So(minX, ShouldEqual, r1+1)
		So(minY, ShouldEqual, r2+1)
		So(maxX, ShouldEqual, r3)
		So(maxY, ShouldEqual, r4)
		reset()
		r.RangePointsMinMaxOpen(f)
		So(minX, ShouldEqual, r1+1)
		So(minY, ShouldEqual, r2+1)
		So(maxX, ShouldEqual, r3-1)
		So(maxY, ShouldEqual, r4-1)
		reset()
		r.RangePointsMinClosedMaxOpen(f)
		So(minX, ShouldEqual, r1)
		So(minY, ShouldEqual, r2)
		So(maxX, ShouldEqual, r3-1)
		So(maxY, ShouldEqual, r4-1)
		So(r.Dx(), ShouldEqual, r3-r1)
		So(r.Dy(), ShouldEqual, r4-r2)
		So(r.Size().Equals(P(r3-r1, r4-r2)), ShouldBeTrue)
		rc := P(1, 1)
		rr := r.Add(rc)
		So(rr.max.X(), ShouldEqual, r.max.X()+rc.X())
		So(rr.max.Y(), ShouldEqual, r.max.Y()+rc.Y())
		So(rr.min.X(), ShouldEqual, r.min.X()+rc.X())
		So(rr.min.Y(), ShouldEqual, r.min.Y()+rc.Y())
		So(rr.Sub(rc).Equals(r), ShouldBeTrue)
		ir := r.Inset(1)
		So(ir.String(), ShouldEqual, fmt.Sprintf("R{(%d,%d)-(%d,%d)}", r1+1, r2+1, r3-1, r4-1))
		So(r.Intersect(ir).Equals(ir), ShouldBeTrue)
		So(r.Union(r.Inset(1)).Equals(r), ShouldBeTrue)
		So(r.Empty(), ShouldBeFalse)
		So(r.Overlaps(ir), ShouldBeTrue)
		So(r.Bounds().Equals(r), ShouldBeTrue)
		So(ir.In(r), ShouldBeTrue)
		So(ir.Expanded(rc).Equals(r), ShouldBeTrue)
		So(ir.ExpandedByMargin(1).Equals(r), ShouldBeTrue)

		rect := Rect[int](0, 0, 3, 3)
		So(rect.IntersectionWithLine(P(1, 0), P(1, 4)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(P(1, 0), P(1, 3)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(P(1, 0), P(1, 2)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(P(1, 0), P(1, 1)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(P(2, 0), P(2, 2)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(P(0, 2), P(1, 2)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(P(0, 1), P(3, 1)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(P(0, 0), P(0, 1)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(P(0, 0), P(0, 2)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(P(0, 0), P(0, 3)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(P(0, 0), P(0, 4)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(P(0, 0), P(1, 0)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(P(0, 0), P(2, 0)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(P(1, 4), P(2, 4)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(P(0, 3), P(3, 3)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(P(1, 3), P(2, 4)), ShouldBeFalse) // should be false
	})
}

// TestRectangle_Accessors 覆盖 Min / Max / Center 三个未覆盖的访问器。
func TestRectangle_Accessors(t *testing.T) {
	Convey("Rectangle Min / Max / Center", t, func() {
		r := Rect[int](0, 0, 10, 6)
		So(r.Min().Equals(P(0, 0)), ShouldBeTrue)
		So(r.Max().Equals(P(10, 6)), ShouldBeTrue)
		// Center: min + (max-min)/2 = (0 + 10/2, 0 + 6/2) = (5, 3)
		So(r.Center().Equals(P(5, 3)), ShouldBeTrue)
	})
}

// TestRectangle_RangePointsMinMaxClosed 覆盖 RangePointsMinMaxClosed
// 三个分支：with==nil 早返 / r==ZR 早返 / 正常遍历 + with 返 false 早退。
func TestRectangle_RangePointsMinMaxClosed(t *testing.T) {
	Convey("RangePointsMinMaxClosed 正常遍历闭区间 [min, max]", t, func() {
		r := Rect[int](0, 0, 2, 1)
		var pts []Point[int]
		r.RangePointsMinMaxClosed(func(p Point[int]) bool {
			pts = append(pts, p)
			return true
		})
		// 闭闭区间：x in [0,2], y in [0,1] → 3x2 = 6 个点
		So(len(pts), ShouldEqual, 6)
	})

	Convey("RangePointsMinMaxClosed with==nil 早返不 panic", t, func() {
		r := Rect[int](0, 0, 5, 5)
		So(func() { r.RangePointsMinMaxClosed(nil) }, ShouldNotPanic)
	})

	Convey("RangePointsMinMaxClosed r==ZR 早返", t, func() {
		var zr Rectangle[int]
		count := 0
		zr.RangePointsMinMaxClosed(func(p Point[int]) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 0)
	})

	Convey("RangePointsMinMaxClosed with 返 false 早退", t, func() {
		r := Rect[int](0, 0, 10, 10)
		count := 0
		r.RangePointsMinMaxClosed(func(p Point[int]) bool {
			count++
			return count < 3 // 第 3 次返 false 早退
		})
		So(count, ShouldEqual, 3)
	})
}

// TestRectangle_RangePoints_EarlyReturns 覆盖 3 个 RangePoints* 子变体的
// with==nil / r==ZR 早返分支（71.4% → 100%）。
func TestRectangle_RangePoints_EarlyReturns(t *testing.T) {
	Convey("RangePointsMinOpenMaxClosed early returns", t, func() {
		r := Rect[int](0, 0, 5, 5)
		So(func() { r.RangePointsMinOpenMaxClosed(nil) }, ShouldNotPanic)
		var zr Rectangle[int]
		count := 0
		zr.RangePointsMinOpenMaxClosed(func(p Point[int]) bool { count++; return true })
		So(count, ShouldEqual, 0)
	})
	Convey("RangePointsMinMaxOpen early returns", t, func() {
		r := Rect[int](0, 0, 5, 5)
		So(func() { r.RangePointsMinMaxOpen(nil) }, ShouldNotPanic)
		var zr Rectangle[int]
		count := 0
		zr.RangePointsMinMaxOpen(func(p Point[int]) bool { count++; return true })
		So(count, ShouldEqual, 0)
	})
	Convey("RangePointsMinClosedMaxOpen early returns", t, func() {
		r := Rect[int](0, 0, 5, 5)
		So(func() { r.RangePointsMinClosedMaxOpen(nil) }, ShouldNotPanic)
		var zr Rectangle[int]
		count := 0
		zr.RangePointsMinClosedMaxOpen(func(p Point[int]) bool { count++; return true })
		So(count, ShouldEqual, 0)
	})
}

// TestRectangle_In_Empty 覆盖 In 的 r.Empty 早返分支（66.7% → 100%）。
func TestRectangle_In_Empty(t *testing.T) {
	Convey("In: 空矩形 r 在任何 s 中都返 true", t, func() {
		var zr Rectangle[int]
		So(zr.In(Rect[int](0, 0, 10, 10)), ShouldBeTrue)
		So(zr.In(Rect[int](100, 100, 200, 200)), ShouldBeTrue)
	})
}

// TestRectangle_HasMethods 覆盖 4 个 Has* 包含判断方法。
func TestRectangle_HasMethods(t *testing.T) {
	Convey("Has* 4 个变体的 边界 / 内部 / 外部 行为", t, func() {
		r := Rect[int](0, 0, 10, 10)

		// 内部点：所有 Has* 都返 true
		inner := P(5, 5)
		So(r.HasMinMaxClosed(inner), ShouldBeTrue)
		So(r.HasMinMaxOpen(inner), ShouldBeTrue)
		So(r.HasMinOpenMaxClosed(inner), ShouldBeTrue)
		So(r.HasMinClosedMaxOpen(inner), ShouldBeTrue)

		// min 角点：MinMaxClosed / MinClosedMaxOpen 含；MinMaxOpen /
		// MinOpenMaxClosed 不含
		minPt := P(0, 0)
		So(r.HasMinMaxClosed(minPt), ShouldBeTrue)
		So(r.HasMinClosedMaxOpen(minPt), ShouldBeTrue)
		So(r.HasMinMaxOpen(minPt), ShouldBeFalse)
		So(r.HasMinOpenMaxClosed(minPt), ShouldBeFalse)

		// max 角点：MinMaxClosed / MinOpenMaxClosed 含；MinMaxOpen /
		// MinClosedMaxOpen 不含
		maxPt := P(10, 10)
		So(r.HasMinMaxClosed(maxPt), ShouldBeTrue)
		So(r.HasMinOpenMaxClosed(maxPt), ShouldBeTrue)
		So(r.HasMinMaxOpen(maxPt), ShouldBeFalse)
		So(r.HasMinClosedMaxOpen(maxPt), ShouldBeFalse)

		// 外部点：所有 Has* 都返 false
		outside := P(15, 15)
		So(r.HasMinMaxClosed(outside), ShouldBeFalse)
		So(r.HasMinMaxOpen(outside), ShouldBeFalse)
		So(r.HasMinOpenMaxClosed(outside), ShouldBeFalse)
		So(r.HasMinClosedMaxOpen(outside), ShouldBeFalse)
	})
}

// TestRectangle_IsIntersect 覆盖 IsIntersect = !Intersect.Empty 的两种结果。
func TestRectangle_IsIntersect(t *testing.T) {
	Convey("IsIntersect 重叠返 true / 不重叠返 false", t, func() {
		r := Rect[int](0, 0, 10, 10)
		// 重叠
		So(r.IsIntersect(Rect[int](5, 5, 15, 15)), ShouldBeTrue)
		// 包含
		So(r.IsIntersect(Rect[int](2, 2, 5, 5)), ShouldBeTrue)
		// 不重叠
		So(r.IsIntersect(Rect[int](20, 20, 30, 30)), ShouldBeFalse)
		// 仅边界相切（max.x = other.min.x）—— Intersect 出 empty 矩形
		So(r.IsIntersect(Rect[int](10, 0, 20, 10)), ShouldBeFalse)
	})
}

// TestRectFromCenterSize 覆盖 RectFromCenterSize（center, size → 矩形）。
func TestRectFromCenterSize(t *testing.T) {
	Convey("RectFromCenterSize", t, func() {
		r := RectFromCenterSize(P(10, 10), P(3, 5))
		// center=(10,10), size=(3,5) → Rect(7, 5, 13, 15)
		So(r.Min().Equals(P(7, 5)), ShouldBeTrue)
		So(r.Max().Equals(P(13, 15)), ShouldBeTrue)
	})
}

// TestRectangle_Union_Empty 覆盖 Union 的 r.Empty / s.Empty 早返分支。
func TestRectangle_Union_Empty(t *testing.T) {
	Convey("Union 处理 empty 矩形", t, func() {
		r := Rect[int](0, 0, 10, 10)
		var zr Rectangle[int]

		// r.Empty → 返 s
		So(zr.Union(r).Equals(r), ShouldBeTrue)
		// s.Empty → 返 r
		So(r.Union(zr).Equals(r), ShouldBeTrue)
		// 两个都 empty
		var zr2 Rectangle[int]
		So(zr.Union(zr2).Empty(), ShouldBeTrue)
	})

	Convey("Union 真扩展（s 在 r 外）4 个 if 都触发", t, func() {
		// r=[5,5 ~ 10,10]; s=[0,0 ~ 15,15]
		// s.min < r.min 让 line 194-196 / 197-199 触发；
		// s.max > r.max 让 line 200-202 / 203-205 触发。
		r := Rect[int](5, 5, 10, 10)
		s := Rect[int](0, 0, 15, 15)
		got := r.Union(s)
		So(got.Min().Equals(P(0, 0)), ShouldBeTrue)
		So(got.Max().Equals(P(15, 15)), ShouldBeTrue)
	})
}

// TestRect_Swap 覆盖 Rect 的 x0>x1 / y0>y1 自动 swap 路径。
func TestRect_Swap(t *testing.T) {
	Convey("Rect 自动 swap 让 min < max", t, func() {
		// 反向坐标 → 自动 swap
		r := Rect[int](10, 10, 0, 0)
		So(r.Min().Equals(P(0, 0)), ShouldBeTrue)
		So(r.Max().Equals(P(10, 10)), ShouldBeTrue)

		// 仅 x 反向
		r2 := Rect[int](10, 0, 0, 10)
		So(r2.Min().Equals(P(0, 0)), ShouldBeTrue)
		So(r2.Max().Equals(P(10, 10)), ShouldBeTrue)

		// 仅 y 反向
		r3 := Rect[int](0, 10, 10, 0)
		So(r3.Min().Equals(P(0, 0)), ShouldBeTrue)
		So(r3.Max().Equals(P(10, 10)), ShouldBeTrue)
	})
}

// TestRectFromMinSize 覆盖 RectFromMinSize（min + size → 矩形）。
func TestRectFromMinSize(t *testing.T) {
	Convey("RectFromMinSize", t, func() {
		r := RectFromMinSize(P(2, 3), P(5, 4))
		So(r.Min().Equals(P(2, 3)), ShouldBeTrue)
		So(r.Max().Equals(P(7, 7)), ShouldBeTrue)
	})
}

// TestRectangle_Inset_Collapse 覆盖 Inset 当 dimension < 2*n 时塌缩到中心
// 的两个分支（Dx < 2n / Dy < 2n）。
func TestRectangle_Inset_Collapse(t *testing.T) {
	Convey("Inset n > Dx/2 时 x 维度塌缩到中心", t, func() {
		r := Rect[int](0, 0, 4, 100)
		// Dx=4 < 2*5=10 → x 塌缩到 (0+4)/2=2; Dy=100 >= 10 → y 正常 inset 5
		ir := r.Inset(5)
		So(ir.Min().Equals(P(2, 5)), ShouldBeTrue)
		So(ir.Max().Equals(P(2, 95)), ShouldBeTrue)
	})

	Convey("Inset n > Dy/2 时 y 维度塌缩到中心", t, func() {
		r := Rect[int](0, 0, 100, 4)
		// Dx=100 >= 10 → x inset 5；Dy=4 < 10 → y 塌缩到 (0+4)/2=2
		ir := r.Inset(5)
		So(ir.Min().Equals(P(5, 2)), ShouldBeTrue)
		So(ir.Max().Equals(P(95, 2)), ShouldBeTrue)
	})

	Convey("Inset 双维度都塌缩", t, func() {
		r := Rect[int](0, 0, 4, 4)
		ir := r.Inset(10)
		// 两个维度都 < 2*10，min/max 都塌到 (2,2)
		So(ir.Min().Equals(P(2, 2)), ShouldBeTrue)
		So(ir.Max().Equals(P(2, 2)), ShouldBeTrue)
		So(ir.Empty(), ShouldBeTrue)
	})
}
