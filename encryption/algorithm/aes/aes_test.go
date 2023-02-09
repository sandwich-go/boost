package aes

import (
	"crypto/rand"
	"github.com/sandwich-go/boost/encryption/key/curve25519"
	"github.com/sandwich-go/boost/xrand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func getTestFrames() [][]byte {
	return [][]byte{
		nil,
		[]byte(""),
		[]byte("time.Duration,[]time.Duration,map[string]*Redis此类的非基础类型的slice或者map都需要辅助指明类型"),
		[]byte(xrand.String(100)),
	}
}

func TestAes(t *testing.T) {
	Convey("test aes encryption", t, func() {
		sk, err := curve25519.GenerateSecretKey()
		So(err, ShouldBeNil)

		for _, frame := range getTestFrames() {
			encryptBody, err0 := Encrypt(frame, sk)
			So(err0, ShouldBeNil)
			decryptDest, err1 := Decrypt(encryptBody, sk)
			So(err1, ShouldBeNil)
			So(decryptDest, ShouldResemble, frame)
		}
	})
}

// 80294 ns/op
func BenchmarkCBCEncrypt(b *testing.B) {
	sk, _ := curve25519.GenerateSecretKey()
	source := make([]byte, 65535)
	_, _ = rand.Read(source)
	for i := 0; i < b.N; i++ {
		_, _ = Encrypt(source, sk)
	}
}

// 70017 ns/op
func BenchmarkCBCDecrypt(b *testing.B) {
	sk, _ := curve25519.GenerateSecretKey()
	source := make([]byte, 65535)
	_, _ = rand.Read(source)
	dest, _ := Encrypt(source, sk)
	for i := 0; i < b.N; i++ {
		_, _ = Decrypt(dest, sk)
	}
}
