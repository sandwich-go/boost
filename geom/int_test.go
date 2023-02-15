// Code generated by tools. DO NOT EDIT.
package geom

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRectangleInt_SnakeRange(t *testing.T) {
	HelixRectRangeFromCenterAndMarginInt(PtInt(2, 2), 2, func(p PointInt) bool {
		t.Log(p.X, p.Y)
		return true
	})

	Convey("test intersection with line", t, func() {
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
