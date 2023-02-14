// Code generated by gotemplate. DO NOT EDIT.

package syncmap

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestIntUint64(t *testing.T) {
	Convey("test sync map", t, func() {
		for _, tr := range []*IntUint64{NewIntUint64()} {
			So(tr.Len(), ShouldBeZeroValue)
			var k, v = __formatKTypeToIntUint64(3), __formatVTypeToIntUint64(4)
			So(tr.Len(), ShouldEqual, 0)
			tr.Store(k, v)
			v1, ok := tr.Load(k)
			So(ok, ShouldBeTrue)
			So(v1, ShouldEqual, v)
		}
	})
}
