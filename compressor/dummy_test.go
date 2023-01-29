package compressor

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDummy(t *testing.T) {
	Convey("dummy flat/inflate", t, func() {
		for _, frame := range getTestFrames() {
			c, err0 := newDummyCompressor()
			So(err0, ShouldBeNil)
			So(c, ShouldNotBeNil)

			testFlatAndInflate(c, frame)
		}
	})

	Convey("dummy flat/inflate n times", t, func() {
		for _, frame := range getTestFrames() {
			c, err := newDummyCompressor()
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)

			testNFlatAndInflate(c, frame, 10)
		}
	})
}
