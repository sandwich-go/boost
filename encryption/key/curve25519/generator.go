package curve25519

import (
	"crypto/rand"
	"golang.org/x/crypto/curve25519"
)

// GenerateSecretKey 生成密钥 SecretKey
func GenerateSecretKey() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GeneratePublicKey 通过密钥 SecretKey 生成公钥
func GeneratePublicKey(secretKey []byte) ([]byte, error) {
	return curve25519.X25519(secretKey, curve25519.Basepoint)
}

// GenerateSharedKey 通过密钥 SecretKey 以及公钥 PublicKey 生成共享公钥 SharedKey
func GenerateSharedKey(secretKey, publicKey []byte) ([]byte, error) {
	return curve25519.X25519(secretKey, publicKey)
}
