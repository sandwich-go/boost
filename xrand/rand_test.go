package xrand

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRand(t *testing.T) {
	Convey("rand", t, func() {
		So(FastRand(), ShouldNotEqual, FastRand())
		So(FastRandUint32n(100), ShouldBeLessThan, 100)
	})
}
