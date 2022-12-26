package curve25519

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPkgCurve25519(t *testing.T) {
	Convey("check alice & bob", t, func() {
		aliceSecretKey := GenerateSecretKey()
		t.Logf("aliceSecretKey: %v\n", aliceSecretKey)
		bobSecretKey := GenerateSecretKey()
		t.Logf("bobSecretKey: %v\n", bobSecretKey)

		alicePublicKey := GeneratePublicKey(aliceSecretKey)
		t.Logf("alicePublicKey: %v\n", alicePublicKey)
		bobPublicKey := GeneratePublicKey(bobSecretKey)
		t.Logf("bobPublicKey: %v\n", bobPublicKey)

		aliceSharedKey := GenerateSharedSecretKey(aliceSecretKey, bobPublicKey)
		t.Logf("aliceSharedKey: %v\n", aliceSharedKey)

		bobSharedKey := GenerateSharedSecretKey(bobSecretKey, alicePublicKey)
		t.Logf("bobSharedKey: %v\n", bobSharedKey)

		So(aliceSharedKey, ShouldResemble, bobSharedKey)
	})
}
