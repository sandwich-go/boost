package vigenere

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 vigenere 包的边界 + nil/empty key 路径。

func TestSanitize(t *testing.T) {
	Convey("Sanitize 仅保留 A-Z / a-z（小写转大写）", t, func() {
		So(string(Sanitize([]byte("Hello World!"))), ShouldEqual, "HELLOWORLD")
		So(string(Sanitize([]byte("ABC123def"))), ShouldEqual, "ABCDEF")
		So(string(Sanitize([]byte("!@#$%^"))), ShouldEqual, "")
	})

	Convey("Sanitize 空 key 返空", t, func() {
		So(Sanitize([]byte{}), ShouldBeEmpty)
		So(Sanitize(nil), ShouldBeNil)
	})
}

func TestEncryptDecrypt_EmptyKey(t *testing.T) {
	Convey("Encrypt 空 key 直接返 src", t, func() {
		src := []byte("HELLO")
		got := Encrypt(src, []byte{})
		So(got, ShouldResemble, src)
		got = Encrypt(src, nil)
		So(got, ShouldResemble, src)
	})

	Convey("Decrypt 空 key 直接返 src", t, func() {
		src := []byte("HELLO")
		got := Decrypt(src, []byte{})
		So(got, ShouldResemble, src)
	})

	Convey("EncryptAndInplace / DecryptAndInplace 空 key 不动 src", t, func() {
		src := []byte("HELLO")
		original := make([]byte, len(src))
		copy(original, src)

		EncryptAndInplace(src, []byte{})
		So(src, ShouldResemble, original)

		DecryptAndInplace(src, nil)
		So(src, ShouldResemble, original)
	})
}

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	Convey("Encrypt + Decrypt 是逆运算", t, func() {
		key := []byte("KEY")
		src := []byte("HELLOWORLD")
		encrypted := Encrypt(src, key)
		So(encrypted, ShouldNotResemble, src) // 真加密了
		decrypted := Decrypt(encrypted, key)
		So(decrypted, ShouldResemble, src)
	})

	Convey("EncryptAndInplace + DecryptAndInplace 互逆", t, func() {
		key := []byte("KEY")
		original := []byte("HELLOWORLD")

		// 加密
		encrypted := make([]byte, len(original))
		copy(encrypted, original)
		EncryptAndInplace(encrypted, key)
		So(encrypted, ShouldNotResemble, original)

		// 解密回来
		DecryptAndInplace(encrypted, key)
		So(encrypted, ShouldResemble, original)
	})
}
