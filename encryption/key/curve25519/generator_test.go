package curve25519

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerator(t *testing.T) {
	Convey("generate SecretKey/PublicKey/SharedKey", t, func() {
		aliceSecretKey, err0 := GenerateSecretKey()
		So(err0, ShouldBeNil)
		t.Logf("aliceSecretKey: %v\n", aliceSecretKey)

		bobSecretKey, err1 := GenerateSecretKey()
		So(err1, ShouldBeNil)
		t.Logf("bobSecretKey: %v\n", bobSecretKey)

		alicePublicKey, err2 := GeneratePublicKey(aliceSecretKey)
		So(err2, ShouldBeNil)
		t.Logf("alicePublicKey: %v\n", alicePublicKey)

		bobPublicKey, err3 := GeneratePublicKey(bobSecretKey)
		So(err3, ShouldBeNil)
		t.Logf("bobPublicKey: %v\n", bobPublicKey)

		aliceSharedKey, err4 := GenerateSharedKey(aliceSecretKey, bobPublicKey)
		So(err4, ShouldBeNil)
		t.Logf("aliceSharedKey: %v\n", aliceSharedKey)

		bobSharedKey, err5 := GenerateSharedKey(bobSecretKey, alicePublicKey)
		So(err5, ShouldBeNil)
		t.Logf("bobSharedKey: %v\n", bobSharedKey)

		So(aliceSharedKey, ShouldResemble, bobSharedKey)
	})
}
