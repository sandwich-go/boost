package graph

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRectangle_SnakeRange(t *testing.T) {
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
