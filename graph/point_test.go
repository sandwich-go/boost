package graph

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPoint(t *testing.T) {
	Convey("test point", t, func() {
		p0, p1 := P[int8](1, 2), P[int8](1, 2)
		So(p0.Equals(p1), ShouldBeTrue)

		So(p0.X(), ShouldEqual, int8(1))
		So(p0.Y(), ShouldEqual, int8(2))

		t.Log(p0)

		p2 := p0.Add(p1)
		So(p2.X(), ShouldEqual, int8(1*2))
		So(p2.Y(), ShouldEqual, int8(2*2))

		p3 := p2.Sub(p1)
		So(p0.Equals(p3), ShouldBeTrue)

		p4 := p0.Mul(2)
		So(p4.Equals(p2), ShouldBeTrue)

		p5 := p4.Div(2)
		So(p5.Equals(p0), ShouldBeTrue)
	})

	Convey("test point in rectangle", t, func() {
		r := R[int16](P[int16](1, 1), P[int16](4, 4))
		So(P[int16](1, 1).In(r), ShouldBeTrue)
		So(P[int16](1, 2).In(r), ShouldBeTrue)
		So(P[int16](1, 3).In(r), ShouldBeTrue)
		So(P[int16](1, 4).In(r), ShouldBeTrue)
		So(P[int16](1, 5).In(r), ShouldBeFalse)

		So(P[int16](2, 1).In(r), ShouldBeTrue)
		So(P[int16](2, 2).In(r), ShouldBeTrue)
		So(P[int16](2, 3).In(r), ShouldBeTrue)
		So(P[int16](2, 4).In(r), ShouldBeTrue)
		So(P[int16](2, 5).In(r), ShouldBeFalse)

		So(P[int16](3, 1).In(r), ShouldBeTrue)
		So(P[int16](3, 2).In(r), ShouldBeTrue)
		So(P[int16](3, 3).In(r), ShouldBeTrue)
		So(P[int16](3, 4).In(r), ShouldBeTrue)
		So(P[int16](3, 5).In(r), ShouldBeFalse)

		So(P[int16](4, 1).In(r), ShouldBeTrue)
		So(P[int16](4, 2).In(r), ShouldBeTrue)
		So(P[int16](4, 3).In(r), ShouldBeTrue)
		So(P[int16](4, 4).In(r), ShouldBeTrue)
		So(P[int16](4, 5).In(r), ShouldBeFalse)

		So(P[int16](5, 1).In(r), ShouldBeFalse)
		So(P[int16](0, 1).In(r), ShouldBeFalse)
		So(P[int16](1, 0).In(r), ShouldBeFalse)
		So(P[int16](0, 0).In(r), ShouldBeFalse)
	})

	Convey("test point distance", t, func() {
		So(P[int16](1, 1).Distance(P[int16](1, 3)), ShouldEqual, 2)
		So(P[int32](1, 3).Distance(P[int32](1, 1)), ShouldEqual, 2)

		So(P[int16](3, 1).Distance(P[int16](1, 1)), ShouldEqual, 2)
		So(P[int32](1, 1).Distance(P[int32](3, 1)), ShouldEqual, 2)

		So(P[int16](1, 1).Distance(P[int16](2, 2)), ShouldEqual, 1.4142135623730951)

		So(P[int16](1, 1).Distance(P[int16](4, 2)), ShouldEqual, 3.1622776601683795)
	})
}
