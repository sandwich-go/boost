// Code generated by tools. DO NOT EDIT.
package geom

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRectangleInt_SnakeRange(t *testing.T) {
	HelixRectRangeFromCenterAndMarginInt(PtInt(2, 2), 2, func(p PointInt) bool {
		t.Log(p.X, p.Y)
		return true
	})

	Convey("test intersection with line", t, func() {
		x := PtInt8(1, 0)
		y := PtInt8(1, 4)
		z := x.Add(y)
		So(z.X, ShouldEqual, 2)
		So(z.Y, ShouldEqual, 4)
		So(z.String(), ShouldEqual, "(2,4)")
		So(z.Sub(x).Eq(y), ShouldBeTrue)
		q := y.Mul(3)
		So(q.X, ShouldEqual, y.X*3)
		So(q.Y, ShouldEqual, y.Y*3)
		So(q.Div(3).Eq(y), ShouldBeTrue)
		var r1, r2, r3, r4 int8 = 1, 0, 5, 5
		r := RectInt8(r1, r2, r3, r4)
		So(x.In(r), ShouldBeTrue)
		So(r.String(), ShouldEqual, fmt.Sprintf("(%d,%d)-(%d,%d)", r1, r2, r3, r4))
		var minX, maxX, minY, maxY int8
		var reset = func() { minX, maxX, minY, maxY = 5, 0, 5, 0 }
		var f = func(p PointInt8) bool {
			if p.X < minX {
				minX = p.X
			}
			if p.Y < minY {
				minY = p.Y
			}
			if p.X > maxX {
				maxX = p.X
			}
			if p.Y > maxY {
				maxY = p.Y
			}
			return true
		}
		reset()
		r.RangePoints(f)
		So(minX, ShouldEqual, r1)
		So(minY, ShouldEqual, r2)
		So(maxX, ShouldEqual, r3)
		So(maxY, ShouldEqual, r4)
		reset()
		r.RangePointsMinClosedMaxOpen(f)
		So(minX, ShouldEqual, r1+1)
		So(minY, ShouldEqual, r2+1)
		So(maxX, ShouldEqual, r3)
		So(maxY, ShouldEqual, r4)
		reset()
		r.RangePointsMinMaxClosed(f)
		So(minX, ShouldEqual, r1+1)
		So(minY, ShouldEqual, r2+1)
		So(maxX, ShouldEqual, r3-1)
		So(maxY, ShouldEqual, r4-1)
		reset()
		r.RangePointsMinOpenMaxClosed(f)
		So(minX, ShouldEqual, r1)
		So(minY, ShouldEqual, r2)
		So(maxX, ShouldEqual, r3-1)
		So(maxY, ShouldEqual, r4-1)
		So(r.Dx(), ShouldEqual, r3-r1)
		So(r.Dy(), ShouldEqual, r4-r2)
		So(r.Size().Eq(PtInt8(r3-r1, r4-r2)), ShouldBeTrue)
		rc := PtInt8(1, 1)
		rr := r.Add(rc)
		So(rr.Max.X, ShouldEqual, r.Max.X+rc.X)
		So(rr.Max.Y, ShouldEqual, r.Max.Y+rc.Y)
		So(rr.Min.X, ShouldEqual, r.Min.X+rc.X)
		So(rr.Min.Y, ShouldEqual, r.Min.Y+rc.Y)
		So(rr.Sub(rc).Eq(r), ShouldBeTrue)
		ir := r.Inset(1)
		So(ir.String(), ShouldEqual, fmt.Sprintf("(%d,%d)-(%d,%d)", r1+1, r2+1, r3-1, r4-1))
		So(r.Intersect(ir).Eq(ir), ShouldBeTrue)
		So(r.Union(r.Inset(1)).Eq(r), ShouldBeTrue)
		So(r.Empty(), ShouldBeFalse)
		So(ZRInt8.Empty(), ShouldBeTrue)
		So(r.Overlaps(ir), ShouldBeTrue)
		So(r.Bounds().Eq(r), ShouldBeTrue)
		So(ir.In(r), ShouldBeTrue)
		So(ir.Expanded(rc).Eq(r), ShouldBeTrue)
		So(ir.ExpandedByMargin(1).Eq(r), ShouldBeTrue)

		rect := RectInt(0, 0, 3, 3)
		So(rect.IntersectionWithLine(PtInt(1, 0), PtInt(1, 4)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(PtInt(1, 0), PtInt(1, 3)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(PtInt(1, 0), PtInt(1, 2)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(PtInt(1, 0), PtInt(1, 1)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(PtInt(2, 0), PtInt(2, 2)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(PtInt(0, 2), PtInt(1, 2)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(PtInt(0, 1), PtInt(3, 1)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(PtInt(0, 0), PtInt(0, 1)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(PtInt(0, 0), PtInt(0, 2)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(PtInt(0, 0), PtInt(0, 3)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(PtInt(0, 0), PtInt(0, 4)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(PtInt(0, 0), PtInt(1, 0)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(PtInt(0, 0), PtInt(2, 0)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(PtInt(1, 4), PtInt(2, 4)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(PtInt(0, 3), PtInt(3, 3)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(PtInt(1, 3), PtInt(2, 4)), ShouldBeFalse) // should be false
	})
}
