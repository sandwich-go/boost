package encrypt

//go:generate stringer -type=Type
type Type byte

const (
	NoneType Type = iota // 无加解密效果
	AESType              // aes 加解密
)
