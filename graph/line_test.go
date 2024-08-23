package graph

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLine(t *testing.T) {
	Convey("test line", t, func() {
		s := P[int8](1, 2)
		e := P[int8](3, 5)
		l0, l1 := L[int8](s, e), L[int8](s, e)
		So(l0.Equals(l1), ShouldBeTrue)

		So(l0.Start().Equals(s), ShouldBeTrue)
		So(l0.End().Equals(e), ShouldBeTrue)

		t.Log(l0)
	})
}
