// Code generated by gotemplate. DO NOT EDIT.

package sarray

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestInt32(t *testing.T) {
	Convey("test sync set", t, func() {
		for _, tr := range []*Int32{NewInt32(), NewSyncInt32()} {
			So(tr.Size(), ShouldBeZeroValue)
			var e0 = __formatToInt32(3)
			tr.Add(e0)
			So(tr.Size(), ShouldEqual, 1)
			tr.Add(e0)
			So(tr.Size(), ShouldEqual, 1)
		}
	})
}
