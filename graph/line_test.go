package graph

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLine(t *testing.T) {
	Convey("test line", t, func() {
		s := P[int8](1, 1)
		e := P[int8](4, 4)
		l0, l1 := L[int8](s, e), L[int8](s, e)
		So(l0.Equals(l1), ShouldBeTrue)

		So(l0.Start().Equals(s), ShouldBeTrue)
		So(l0.End().Equals(e), ShouldBeTrue)

		t.Log(l0)

		var l2 Line[int64]
		var l2PointNum int
		l2.RangePoints(func(p Point[int64]) bool {
			l2PointNum++
			return true
		})
		So(l2PointNum, ShouldBeZeroValue)

		var l3 = L[int8](s, s)
		var l3Points []Point[int8]
		l3.RangePoints(func(p Point[int8]) bool {
			l3Points = append(l3Points, p)
			return true
		})
		So(len(l3Points), ShouldEqual, 1)
		for _, p := range l3Points {
			So(p.Equals(s), ShouldBeTrue)
		}

		var l0Points []Point[int8]
		l0.RangePoints(func(p Point[int8]) bool {
			l0Points = append(l0Points, p)
			return true
		})
		So(len(l0Points), ShouldEqual, 4)
	})
}
