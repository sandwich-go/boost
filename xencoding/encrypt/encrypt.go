package encrypt

import (
	"context"
	"errors"
	"github.com/sandwich-go/boost/xcrypto/algorithm/aes"
	"github.com/sandwich-go/boost/xencoding"
	"github.com/sandwich-go/boost/xpanic"
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

// KeySetter key 的设置
// 如果自定义的 Codec 实现了 KeySetter 接口，那么在 NewCodec 的时候，会将 key 设置进去
type KeySetter interface {
	SetKey(key []byte)
}

var (
	// NoneCodec 无加解密效果
	NoneCodec = noneCodec{}
	// AESCodec aes 加解密
	AESCodec = aesCodec{}
)

// SetKey 设置加解密 key
func SetKey(key []byte) { AESCodec.key = key }

var codecs = map[Type]xencoding.Codec{
	NoneType: NoneCodec,
	AESType:  AESCodec,
}

func init() {
	for _, v := range codecs {
		xencoding.RegisterCodec(v)
	}
}

// Register 注册自定义的加解密 Codec ，该方法非协程安全
func Register(t Type, codec xencoding.Codec) {
	_, exists := codecs[t]
	xpanic.WhenTrue(exists, "register called twice for codec, %d", t)
	codecs[t] = codec
	xencoding.RegisterCodec(codec)
}

type noneCodec struct{}

func (c noneCodec) Name() string { return NoneCodecName }
func (c noneCodec) Marshal(_ context.Context, v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errCodecMarshalParam
	} else {
		return data, nil
	}
}
func (c noneCodec) Unmarshal(_ context.Context, bytes []byte, v interface{}) error {
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
func (c aesCodec) Marshal(_ context.Context, v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errCodecMarshalParam
	} else {
		return aes.Encrypt(data, c.key)
	}
}

func (c aesCodec) Unmarshal(_ context.Context, bytes []byte, v interface{}) error {
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
	xencoding.Codec
	encryptType Type
}

// NewCodec 通过类型创建解压缩 Codec
func NewCodec(encryptType Type, key []byte) xencoding.Codec {
	c := codec{encryptType: encryptType}
	switch encryptType {
	case AESType:
		c.Codec = aesCodec{key: key}
	case NoneType:
		c.Codec = noneCodec{}
	default:
		c.Codec = codecs[encryptType]
	}
	if c.Codec == nil {
		xpanic.WhenError(errCodecNoFound)
	}
	if cc, ok := c.Codec.(KeySetter); ok {
		cc.SetKey(key)
	}
	return c
}

// Name 返回 Codec 名
func (c codec) Name() string { return "encrypt" }

// Marshal 编码
func (c codec) Marshal(ctx context.Context, v interface{}) ([]byte, error) {
	if c.Codec == nil {
		return nil, errCodecNoFound
	}
	return c.Codec.Marshal(ctx, v)
}

// Unmarshal 解码
func (c codec) Unmarshal(ctx context.Context, bytes []byte, v interface{}) error {
	if c.Codec == nil {
		return errCodecNoFound
	}
	return c.Codec.Unmarshal(ctx, bytes, v)
}
