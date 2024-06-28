package geom

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRectangle_SnakeRange(t *testing.T) {
	HelixRectRangeFromCenterAndMargin[int](Pt(2, 2), 2, func(p Point[int]) bool {
		t.Log(p.X, p.Y)
		return true
	})

	Convey("test intersection with line", t, func() {
		x := Pt(1, 0)
		y := Pt(1, 4)
		z := x.Add(y)
		So(z.X, ShouldEqual, 2)
		So(z.Y, ShouldEqual, 4)
		So(z.String(), ShouldEqual, "(2,4)")
		So(z.Sub(x).Eq(y), ShouldBeTrue)
		q := y.Mul(3)
		So(q.X, ShouldEqual, y.X*3)
		So(q.Y, ShouldEqual, y.Y*3)
		So(q.Div(3).Eq(y), ShouldBeTrue)
		var r1, r2, r3, r4 int = 1, 0, 5, 5
		r := Rect[int](r1, r2, r3, r4)
		So(x.In(r), ShouldBeTrue)
		So(r.String(), ShouldEqual, fmt.Sprintf("(%d,%d)-(%d,%d)", r1, r2, r3, r4))
		var minX, maxX, minY, maxY int
		var reset = func() { minX, maxX, minY, maxY = 5, 0, 5, 0 }
		var f = func(p Point[int]) bool {
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
		So(r.Size().Eq(Pt(r3-r1, r4-r2)), ShouldBeTrue)
		rc := Pt(1, 1)
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
		So(r.Overlaps(ir), ShouldBeTrue)
		So(r.Bounds().Eq(r), ShouldBeTrue)
		So(ir.In(r), ShouldBeTrue)
		So(ir.Expanded(rc).Eq(r), ShouldBeTrue)
		So(ir.ExpandedByMargin(1).Eq(r), ShouldBeTrue)

		rect := Rect[int](0, 0, 3, 3)
		So(rect.IntersectionWithLine(Pt(1, 0), Pt(1, 4)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(Pt(1, 0), Pt(1, 3)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(Pt(1, 0), Pt(1, 2)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(Pt(1, 0), Pt(1, 1)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(Pt(2, 0), Pt(2, 2)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(Pt(0, 2), Pt(1, 2)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(Pt(0, 1), Pt(3, 1)), ShouldBeTrue)  // should be true
		So(rect.IntersectionWithLine(Pt(0, 0), Pt(0, 1)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(Pt(0, 0), Pt(0, 2)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(Pt(0, 0), Pt(0, 3)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(Pt(0, 0), Pt(0, 4)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(Pt(0, 0), Pt(1, 0)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(Pt(0, 0), Pt(2, 0)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(Pt(1, 4), Pt(2, 4)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(Pt(0, 3), Pt(3, 3)), ShouldBeFalse) // should be false
		So(rect.IntersectionWithLine(Pt(1, 3), Pt(2, 4)), ShouldBeFalse) // should be false
	})
}
