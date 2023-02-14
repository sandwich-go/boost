// Code generated by gotemplate. DO NOT EDIT.

package sarray

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestInt32(t *testing.T) {
	Convey("test sync array", t, func() {
		for _, tr := range []*Int32{NewInt32(), NewSyncInt32()} {
			So(tr.Len(), ShouldBeZeroValue)
			_, exists := tr.Get(0)
			So(exists, ShouldBeFalse)
			So(tr.Empty(), ShouldBeTrue)
			var e0 = __formatToInt32(3)
			tr.PushLeft(e0)
			So(tr.Len(), ShouldEqual, 1)
			e := tr.At(0)
			So(e, ShouldEqual, e0)
		}
	})
}
