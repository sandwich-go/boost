package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func CBCEncrypt(src []byte, key []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("CBCEncrypt empty src")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = Padding(src, block.BlockSize())
	iv := key[:block.BlockSize()]
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(src, src)
	return src, nil
}

func CBCDecrypt(src []byte, key []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := key[:block.BlockSize()]
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(src, src)
	src = UnPadding(src)
	return src, nil
}
