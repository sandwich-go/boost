package xpool

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPool(t *testing.T) {
	Convey("pool should work ok", t, func() {
		var nexStr = "new string"
		var p = NewPool[string](func() string {
			return nexStr
		})
		var v = p.Get()
		So(v, ShouldEqual, nexStr)
		v = "string"
		p.Put(v)

		v = p.Get()
		So(v, ShouldNotEqual, nexStr)
	})
}
