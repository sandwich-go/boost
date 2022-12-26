package aes

import (
	"crypto/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAes(t *testing.T) {
	Convey("test aes encryption", t, func() {
		sharedSecretKey := []byte{77, 159, 192, 198, 37, 164, 129, 197, 223, 156, 117, 46, 34, 71, 50, 161, 203, 9, 180, 255, 108, 158, 37, 46, 128, 236, 180, 168, 3, 147, 203, 41}
		source := []byte("hello world hello world")
		encryptBody, err := CBCEncrypt(source, sharedSecretKey)
		So(err, ShouldBeNil)
		decryptDest, err := CBCDecrypt(encryptBody, sharedSecretKey)
		So(err, ShouldBeNil)
		So(decryptDest, ShouldResemble, source)
	})
}

// 80294 ns/op
func BenchmarkCBCEncrypt(b *testing.B) {
	sharedSecretKey := []byte{77, 159, 192, 198, 37, 164, 129, 197, 223, 156, 117, 46, 34, 71, 50, 161, 203, 9, 180, 255, 108, 158, 37, 46, 128, 236, 180, 168, 3, 147, 203, 41}
	source := make([]byte, 65535)
	_, _ = rand.Read(source)
	for i := 0; i < b.N; i++ {
		_, _ = CBCEncrypt(source, sharedSecretKey)
	}
}

// 70017 ns/op
func BenchmarkCBCDecrypt(b *testing.B) {
	sharedSecretKey := []byte{77, 159, 192, 198, 37, 164, 129, 197, 223, 156, 117, 46, 34, 71, 50, 161, 203, 9, 180, 255, 108, 158, 37, 46, 128, 236, 180, 168, 3, 147, 203, 41}
	source := make([]byte, 65535)
	_, _ = rand.Read(source)
	dest, _ := CBCEncrypt(source, sharedSecretKey)
	for i := 0; i < b.N; i++ {
		_, _ = CBCDecrypt(dest, sharedSecretKey)
	}
}
