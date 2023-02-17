package ringbuf

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var buf = New(10)

func TestRingbuf(t *testing.T) {
	Convey("ring buf test", t, func() {
		So(buf.Capacity(), ShouldEqual, 10)
		err := buf.Write([]byte("helloworld"))
		//helloworld
		So(err, ShouldBeNil)
		So(buf.Size(), ShouldEqual, 10)
		So(buf.start, ShouldEqual, 0)
		err = buf.Write([]byte("a"))
		So(err, ShouldEqual, NotEnoughMem)

		tmp := make([]byte, 5)
		n := buf.Read(tmp, 5)
		//world
		So(n, ShouldEqual, 5)
		So(len(tmp), ShouldEqual, 5)
		So(buf.Size(), ShouldEqual, 5)
		So(buf.start, ShouldEqual, 5)
		So(string(tmp), ShouldEqual, "hello")

		err = buf.Write([]byte("1234"))
		// world1234
		So(err, ShouldBeNil)
		So(buf.Size(), ShouldEqual, 9)
		So(buf.start, ShouldEqual, 5)

		n = buf.Read(tmp, 4)
		// d1234
		So(n, ShouldEqual, 4)
		So(string(tmp), ShouldEqual, "worlo")
		So(buf.Size(), ShouldEqual, 5)
		So(buf.start, ShouldEqual, 9)

		So(len(buf.PreUse(5)), ShouldEqual, 5)
		buf.Read(tmp, 3)
		So(string(tmp), ShouldEqual, "d12lo")

		buf.RealUse(1)
		buf.Read(tmp, 3)
		So(string(tmp), ShouldEqual, "34olo")

	})
}
