// Code generated by gotemplate. DO NOT EDIT.

package sarray

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestInt64(t *testing.T) {
	Convey("test sync list", t, func() {
		for _, tr := range []*Int64{NewInt64(), NewSyncInt64()} {
			So(tr.Len(), ShouldBeZeroValue)
			var e0 = __formatToInt64(3)
			tr.PushBack(e0)
			So(tr.Len(), ShouldEqual, 1)
			e := tr.PopBack()
			So(e, ShouldEqual, e0)
		}
	})
}
