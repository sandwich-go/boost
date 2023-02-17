package xcompress

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCompress(t *testing.T) {
	Convey("init compressor", t, func() {
		c, err := New(WithLevel(BestCompression + 1))
		So(err, ShouldNotBeNil)
		So(c, ShouldBeNil)

		c, err = New(WithLevel(HuffmanOnly - 1))
		So(err, ShouldNotBeNil)
		So(c, ShouldBeNil)

		c, err = New(WithLevel(HuffmanOnly))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		c, err = New(WithType(GZIP))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		c, err = New(WithType(Snappy))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		c, err = New(WithType(Dummy))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
	})

	Convey("flat/inflate", t, func() {
		for _, frame := range getTestFrames() {
			testFlatAndInflate(Default, frame)
		}
		for _, frame := range getTestFrames() {
			for i := Dummy; i <= Snappy; i++ {
				switch i {
				case GZIP:
					for lvl := HuffmanOnly; lvl <= BestCompression; lvl++ {
						c, err := New(WithType(i), WithLevel(lvl))
						So(err, ShouldBeNil)
						So(c, ShouldNotBeNil)
						testFlatAndInflate(c, frame)
					}
				default:
					c, err := New(WithType(i))
					So(err, ShouldBeNil)
					So(c, ShouldNotBeNil)
					testFlatAndInflate(c, frame)
				}
			}
		}
	})
}
