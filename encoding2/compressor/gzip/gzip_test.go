package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const testStr = "wo ai wo jia"

func TestCompress(t *testing.T) {
	Convey("Gzip Compress and Decompress normal data", t, func() {
		SetLevel(gzip.DefaultCompression)
		c := createCompressor()
		buffer := &bytes.Buffer{}
		en, err := c.Compress(buffer)
		So(err, ShouldBeNil)
		_, err = en.Write([]byte(testStr))
		So(err, ShouldBeNil)
		err = en.Flush()
		So(err, ShouldBeNil)
		err = en.Close()
		So(err, ShouldBeNil)
		So(buffer.Bytes(), ShouldNotBeNil)

		de, err := c.Decompress(bytes.NewReader(buffer.Bytes()))
		So(err, ShouldBeNil)
		r, err := ioutil.ReadAll(de)
		So(err, ShouldBeNil)
		So(string(r), ShouldEqual, testStr)
		err = de.Close()
		So(err, ShouldBeNil)
	})
}

func TestCompressNil(t *testing.T) {
	Convey("Gzip Compress and Decompress nil", t, func() {
		c := createCompressor()
		buffer := &bytes.Buffer{}
		en, err := c.Compress(buffer)
		So(err, ShouldBeNil)
		_, err = en.Write(nil)
		So(err, ShouldBeNil)
		err = en.Flush()
		So(err, ShouldBeNil)
		err = en.Close()
		So(err, ShouldBeNil)
		So(buffer.Bytes(), ShouldNotBeNil)

		de, err := c.Decompress(bytes.NewReader(buffer.Bytes()))
		So(err, ShouldBeNil)
		r, err := ioutil.ReadAll(de)
		So(err, ShouldBeNil)
		So(string(r), ShouldEqual, "")
		err = de.Close()
		So(err, ShouldBeNil)
	})
}

func TestCompressMultiEmptyPool(t *testing.T) {
	Convey("Multi gzip Compress and Decompress normal data for empty pool", t, func() {
		c := createCompressor()
		b1 := &bytes.Buffer{}
		en1, err := c.Compress(b1)
		So(err, ShouldBeNil)
		_, err = en1.Write([]byte(testStr))
		So(err, ShouldBeNil)

		b2 := &bytes.Buffer{}
		en2, err := c.Compress(b2)
		So(en2, ShouldNotEqual, en1)
		So(err, ShouldBeNil)
		_, err = en2.Write([]byte(testStr))
		So(err, ShouldBeNil)

		So(en1, ShouldNotEqual, en2)

		err = en1.Flush()
		So(err, ShouldBeNil)
		err = en1.Close()
		So(err, ShouldBeNil)
		So(b1.Bytes(), ShouldNotBeNil)

		err = en2.Flush()
		So(err, ShouldBeNil)
		err = en2.Close()
		So(err, ShouldBeNil)
		So(b2.Bytes(), ShouldNotBeNil)

		de1, err := c.Decompress(bytes.NewReader(b1.Bytes()))
		So(err, ShouldBeNil)
		r1, err := ioutil.ReadAll(de1)
		So(err, ShouldBeNil)
		So(string(r1), ShouldEqual, testStr)

		de2, err := c.Decompress(bytes.NewReader(b2.Bytes()))
		So(err, ShouldBeNil)
		So(de1, ShouldNotEqual, de2)
		r2, err := ioutil.ReadAll(de2)
		So(err, ShouldBeNil)
		So(string(r2), ShouldEqual, testStr)

		err = de1.Close()
		So(err, ShouldBeNil)
		err = de2.Close()
		So(err, ShouldBeNil)
	})
}

func TestCompressMultiNotEmptyPool(t *testing.T) {
	Convey("Multi gzip Compress and Decompress normal data for not empty pool", t, func() {
		c := createCompressor()
		b1 := &bytes.Buffer{}
		en1, err := c.Compress(b1)
		So(err, ShouldBeNil)
		_, err = en1.Write([]byte(testStr))
		So(err, ShouldBeNil)
		err = en1.Flush()
		So(err, ShouldBeNil)
		err = en1.Close()
		So(err, ShouldBeNil)
		So(b1.Bytes(), ShouldNotBeNil)

		b2 := &bytes.Buffer{}
		en2, err := c.Compress(b2)
		So(en2, ShouldEqual, en1)
		So(err, ShouldBeNil)
		_, err = en2.Write([]byte(testStr))
		So(err, ShouldBeNil)
		err = en2.Flush()
		So(err, ShouldBeNil)
		err = en2.Close()
		So(err, ShouldBeNil)
		So(b2.Bytes(), ShouldNotBeNil)

		de1, err := c.Decompress(bytes.NewReader(b1.Bytes()))
		So(err, ShouldBeNil)
		r1, err := ioutil.ReadAll(de1)
		So(err, ShouldBeNil)
		So(string(r1), ShouldEqual, testStr)
		err = de1.Close()
		So(err, ShouldBeNil)

		de2, err := c.Decompress(bytes.NewReader(b1.Bytes()))
		So(de1, ShouldEqual, de2)
		So(err, ShouldBeNil)
		r2, err := ioutil.ReadAll(de2)
		So(err, ShouldBeNil)
		So(string(r2), ShouldEqual, testStr)
		err = de2.Close()
		So(err, ShouldBeNil)
	})
}
