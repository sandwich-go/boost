package compressor

import (
	"context"
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

func TestCompressor(t *testing.T) {
	Convey("compress", t, func() {
		c := NewCodec(SnappyType + 1)
		_, err := c.Marshal(context.Background(), "")
		So(err, ShouldEqual, errCodecNoFound)

		c = NewCodec(SnappyType)
		_, err = c.Marshal(context.Background(), "")
		So(err, ShouldEqual, errCodecMarshalParam)

		for i := NoneType; i <= SnappyType; i++ {
			c = NewCodec(i)
			for _, frame := range getTestFrames() {
				mf, err0 := c.Marshal(context.Background(), frame)
				So(err0, ShouldBeNil)

				var uf []byte
				err = c.Unmarshal(context.Background(), mf, &uf)
				So(err, ShouldBeNil)

				So(frame, ShouldResemble, uf)
			}
		}
	})
}
