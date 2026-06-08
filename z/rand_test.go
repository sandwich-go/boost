package z

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRand(t *testing.T) {
	Convey("rand", t, func() {
		So(FastRand(), ShouldNotEqual, FastRand())
		So(FastRandUint32n(100), ShouldBeLessThan, 100)
	})

	Convey("FastRandUint32n maxN=0 早返 0", t, func() {
		So(FastRandUint32n(0), ShouldEqual, uint32(0))
	})

	Convey("FastRandInt", t, func() {
		// min < max 正常
		v := FastRandInt(10, 20)
		So(v, ShouldBeBetweenOrEqual, 10, 20)
		// min == max 返 min
		So(FastRandInt(5, 5), ShouldEqual, 5)
		// min > max 走 'return min' 早返
		So(FastRandInt(100, 50), ShouldEqual, 100)
	})
}
