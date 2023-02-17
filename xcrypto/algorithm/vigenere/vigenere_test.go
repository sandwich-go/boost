package vigenere

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVigenere(t *testing.T) {
	Convey("test vigenere encryption", t, func() {
		So(Sanitize([]byte("123456789")), ShouldBeEmpty)
		a := Sanitize([]byte("123456789ABC"))
		So(len(a), ShouldEqual, 3)
		b := "GHIJKLMNO"
		c := Encrypt([]byte(b), a)
		So(b, ShouldEqual, string(Decrypt(c, a)))

		b1 := []byte(b)
		EncryptAndInplace(b1, a)
		So(b, ShouldNotEqual, string(b1))
		DecryptAndInplace(b1, a)
		So(b, ShouldEqual, string(b1))
	})
}
