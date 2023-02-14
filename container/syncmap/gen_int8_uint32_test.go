// Code generated by gotemplate. DO NOT EDIT.

package syncmap

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestInt8Uint32(t *testing.T) {
	Convey("test sync map", t, func() {
		for _, tr := range []*Int8Uint32{NewInt8Uint32()} {
			So(tr.Len(), ShouldBeZeroValue)
			var k, v = __formatKTypeToInt8Uint32(3), __formatVTypeToInt8Uint32(4)
			So(tr.Len(), ShouldEqual, 0)
			tr.Store(k, v)
			v1, ok := tr.Load(k)
			So(ok, ShouldBeTrue)
			So(v1, ShouldEqual, v)
		}
	})
}
