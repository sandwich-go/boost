package aes

import (
	"github.com/coreos/pkg/cryptoutil"
)

func CBCEncrypt(src []byte, key []byte) ([]byte, error) {
	return cryptoutil.AESEncrypt(src, key)
}

func CBCDecrypt(src []byte, key []byte) ([]byte, error) {
	return cryptoutil.AESDecrypt(src, key)
}
