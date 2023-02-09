package compressor

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSnappy(t *testing.T) {
	Convey("snappy flat/inflate", t, func() {
		for _, frame := range getTestFrames() {
			c, err0 := newSnappyCompressor()
			So(err0, ShouldBeNil)
			So(c, ShouldNotBeNil)

			testFlatAndInflate(c, frame)
		}
	})

	Convey("snappy flat/inflate n times", t, func() {
		for _, frame := range getTestFrames() {
			c, err := newSnappyCompressor()
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)

			testNFlatAndInflate(c, frame, 10)
		}
	})
}
