package encrypt

import (
	"github.com/sandwich-go/boost/xcrypto/key/curve25519"
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

func getKey() []byte {
	key, err := curve25519.GenerateSecretKey()
	So(err, ShouldBeNil)
	return key
}

func TestEncrypt(t *testing.T) {
	Convey("encrypt", t, func() {
		c := NewCodec(AESType+1, getKey())
		_, err := c.Marshal("")
		So(err, ShouldEqual, errCodecNoFound)

		c = NewCodec(AESType, getKey())
		_, err = c.Marshal("")
		So(err, ShouldEqual, errCodecMarshalParam)

		for i := NoneType; i <= AESType; i++ {
			c = NewCodec(i, getKey())
			for _, frame := range getTestFrames() {
				mf, err0 := c.Marshal(frame)
				So(err0, ShouldBeNil)

				var uf []byte
				err = c.Unmarshal(mf, &uf)
				So(err, ShouldBeNil)

				So(frame, ShouldResemble, uf)
			}
		}
	})
}
