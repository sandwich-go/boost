package aes

import (
	"github.com/coreos/pkg/cryptoutil"
)

// Encrypt 使用 key 进行加密
func Encrypt(src []byte, key []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, nil
	}
	return cryptoutil.AESEncrypt(src, key)
}

// Decrypt 使用 key 进行解密
func Decrypt(src []byte, key []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, nil
	}
	return cryptoutil.AESDecrypt(src, key)
}
