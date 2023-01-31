package encrypt

import (
	"errors"
	"github.com/sandwich-go/boost/encoding2"
	"github.com/sandwich-go/boost/encryption/algorithm/aes"
)

var (
	errCodecMarshalParam   = errors.New("encrypt codec marshal must be []byte parameter")
	errCodecUnmarshalParam = errors.New("encrypt codec unmarshal must be *[]byte parameter")
	errCodecNoFound        = errors.New("encrypt codec not found")
)

const (
	// NoneCodecName 无加解密效果名称，可以通过 encoding2.GetCodec(NoneCodecName) 获取对应的 Codec
	NoneCodecName = "none_encrypt"
	// AESCodecName aes 加解密名称，可以通过 encoding2.GetCodec(AESCodecName) 获取对应的 Codec
	AESCodecName = "aes_encrypt"
)

var (
	// NoneCodec 无加解密效果
	NoneCodec = noneCodec{}
	// AESCodec aes 加解密
	AESCodec = aesCodec{}
)

// SetKey 设置加解密 key
func SetKey(key []byte) { AESCodec.key = key }

var codecs = map[Type]encoding2.Codec{
	NoneType: NoneCodec,
	AESType:  AESCodec,
}

func init() {
	for _, v := range codecs {
		encoding2.RegisterCodec(v)
	}
}

type noneCodec struct{}

func (c noneCodec) Name() string { return NoneCodecName }
func (c noneCodec) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errCodecMarshalParam
	} else {
		return data, nil
	}
}
func (c noneCodec) Unmarshal(bytes []byte, v interface{}) error {
	v1, ok := v.(*[]byte)
	if !ok {
		return errCodecUnmarshalParam
	}
	*v1 = bytes
	return nil
}

type aesCodec struct {
	key []byte
}

func (c aesCodec) Name() string { return AESCodecName }
func (c aesCodec) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errCodecMarshalParam
	} else {
		return aes.Encrypt(data, c.key)
	}
}

func (c aesCodec) Unmarshal(bytes []byte, v interface{}) error {
	v1, ok := v.(*[]byte)
	if !ok {
		return errCodecUnmarshalParam
	}
	data, err := aes.Decrypt(bytes, c.key)
	if err != nil {
		return err
	}
	*v1 = data
	return nil
}

type codec struct {
	encoding2.Codec
	encryptType Type
}

// NewCodec 通过类型创建解压缩 Codec
func NewCodec(encryptType Type, key []byte) encoding2.Codec {
	c := codec{encryptType: encryptType}
	switch encryptType {
	case AESType:
		c.Codec = aesCodec{key: key}
	case NoneType:
		c.Codec = noneCodec{}
	}
	return c
}

// Name 返回 Codec 名
func (c codec) Name() string { return "encrypt" }

// Marshal 编码
func (c codec) Marshal(v interface{}) ([]byte, error) {
	if c.Codec == nil {
		return nil, errCodecNoFound
	}
	return c.Codec.Marshal(v)
}

// Unmarshal 解码
func (c codec) Unmarshal(bytes []byte, v interface{}) error {
	if c.Codec == nil {
		return errCodecNoFound
	}
	return c.Codec.Unmarshal(bytes, v)
}
