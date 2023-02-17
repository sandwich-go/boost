package xcompress

import (
	"github.com/sandwich-go/boost/xrand"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func getTestFrames() [][]byte {
	return [][]byte{
		nil,
		[]byte(""),
		[]byte("time.Duration,[]time.Duration,map[string]*Redis此类的非基础类型的slice或者map都需要辅助指明类型"),
		[]byte(xrand.String(100)),
	}
}

func testFlatAndInflate(c Compressor, frame []byte) {
	out, err1 := c.Flat(frame)
	So(err1, ShouldBeNil)
	So(len(out), ShouldNotBeEmpty)

	in, err2 := c.Inflate(out)
	So(err2, ShouldBeNil)
	So(in, ShouldResemble, frame)
}

func testNFlatAndInflate(c Compressor, frame []byte, maxTimes uint32) {
	var err error
	var out = frame
	var n = maxTimes
	var i uint32
	for i = 0; i < n; i++ {
		out, err = c.Flat(out)
		So(err, ShouldBeNil)
		So(len(out), ShouldNotBeEmpty)
	}
	for i = 0; i < n; i++ {
		out, err = c.Inflate(out)
		So(err, ShouldBeNil)
	}
	So(out, ShouldResemble, frame)
}

func TestGZIP(t *testing.T) {
	Convey("init gzip compressor", t, func() {
		c, err := newGzipCompressor(BestCompression + 1)
		So(err, ShouldNotBeNil)
		So(c, ShouldBeNil)

		c, err = newGzipCompressor(HuffmanOnly - 1)
		So(err, ShouldNotBeNil)
		So(c, ShouldBeNil)

		c, err = newGzipCompressor(HuffmanOnly)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
	})

	Convey("gzip flat/inflate", t, func() {
		for _, frame := range getTestFrames() {
			for lvl := HuffmanOnly; lvl <= BestCompression; lvl++ {
				c, err0 := newGzipCompressor(lvl)
				So(err0, ShouldBeNil)
				So(c, ShouldNotBeNil)

				testFlatAndInflate(c, frame)
			}
		}
	})

	Convey("gzip flat/inflate n times", t, func() {
		for _, frame := range getTestFrames() {
			for lvl := HuffmanOnly; lvl <= BestCompression; lvl++ {
				c, err := newGzipCompressor(lvl)
				So(err, ShouldBeNil)
				So(c, ShouldNotBeNil)

				testNFlatAndInflate(c, frame, 10)
			}
		}
	})
}
