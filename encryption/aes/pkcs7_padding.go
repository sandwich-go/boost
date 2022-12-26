package aes

import "bytes"

func Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func UnPadding(src []byte) []byte {
	length := len(src)
	if length == 0 {
		return nil
	}
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}
