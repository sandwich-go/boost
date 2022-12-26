package curve25519

import (
	"crypto/rand"

	"golang.org/x/crypto/curve25519"
)

func GenerateSecretKey() []byte {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
		return nil
	}
	return b
}

func GeneratePublicKey(secretKey []byte) []byte {
	pk, err := curve25519.X25519(secretKey, curve25519.Basepoint)
	if err != nil {
		panic(err)
		return nil
	}
	return pk
}

func GenerateSharedSecretKey(mySecretKey, otherPublicKey []byte) []byte {
	pk, err := curve25519.X25519(mySecretKey, otherPublicKey)
	if err != nil {
		panic(err)
		return nil
	}
	return pk
}
