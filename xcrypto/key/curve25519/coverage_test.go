package curve25519

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 curve25519 包的 GenerateSecretKey 错误路径 / 通用流程。
// 注：GenerateSecretKey 错误路径需要 crypto/rand 出错，正常机器上不会发生；
// GeneratePublicKey / GenerateSharedKey 是 stdlib 包装，主要测调用流程。

func TestGenerateSecretKey(t *testing.T) {
	Convey("GenerateSecretKey 返回 32 字节随机 key", t, func() {
		key, err := GenerateSecretKey()
		So(err, ShouldBeNil)
		So(len(key), ShouldEqual, 32)

		// 多次调用得到不同 key（极大概率）
		key2, _ := GenerateSecretKey()
		So(key, ShouldNotResemble, key2)
	})
}

func TestGeneratePublicKey_AndSharedKey(t *testing.T) {
	Convey("GeneratePublicKey 从 secretKey 派生 publicKey", t, func() {
		sk, _ := GenerateSecretKey()
		pk, err := GeneratePublicKey(sk)
		So(err, ShouldBeNil)
		So(len(pk), ShouldEqual, 32)
	})

	Convey("GenerateSharedKey 双方 ECDH 协商出相同共享 key", t, func() {
		// Alice 与 Bob 各自生成 secret + public
		aliceSK, _ := GenerateSecretKey()
		alicePK, _ := GeneratePublicKey(aliceSK)
		bobSK, _ := GenerateSecretKey()
		bobPK, _ := GeneratePublicKey(bobSK)

		// Alice 用自己 SK + Bob 的 PK 生成共享
		aliceShared, err := GenerateSharedKey(aliceSK, bobPK)
		So(err, ShouldBeNil)
		// Bob 用自己 SK + Alice 的 PK 生成共享
		bobShared, err := GenerateSharedKey(bobSK, alicePK)
		So(err, ShouldBeNil)
		// ECDH：双方应得到相同共享 key
		So(aliceShared, ShouldResemble, bobShared)
	})

	Convey("GeneratePublicKey 在非法 secretKey 长度时返错误", t, func() {
		_, err := GeneratePublicKey([]byte{1, 2, 3}) // 长度不足 32
		So(err, ShouldNotBeNil)
	})
}
