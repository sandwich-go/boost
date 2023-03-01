package compressor

import (
	"errors"
	"github.com/sandwich-go/boost/xcompress"
	"github.com/sandwich-go/boost/xencoding"
	"github.com/sandwich-go/boost/xpanic"
)

var (
	errCodecMarshalParam   = errors.New("compressor codec marshal must be []byte parameter")
	errCodecUnmarshalParam = errors.New("compressor codec unmarshal must be *[]byte parameter")
	errCodecNoFound        = errors.New("compressor codec not found")
)

const (
	// NoneCodecName 无压缩效果名称，可以通过 encoding2.GetCodec(NoneCodecName) 获取对应的 Codec
	NoneCodecName = "none_compressor"
	// GZIPCodecName gzip 压缩名称，可以通过 encoding2.GetCodec(GZIPCodecName) 获取对应的 Codec
	GZIPCodecName = "gzip_compressor"
	// SnappyCodecName snappy 压缩名称，可以通过 encoding2.GetCodec(SnappyCodecName) 获取对应的 Codec
	SnappyCodecName = "snappy_compressor"
)

var (
	// NoneCodec 无压缩效果
	NoneCodec = noneCodec{newBaseCodec(xcompress.WithType(xcompress.Dummy))}
	// GZIPCodec gzip 压缩，使用 compressor.DefaultCompression 等级
	GZIPCodec = gzipCodec{newBaseCodec(xcompress.WithType(xcompress.GZIP), xcompress.WithLevel(xcompress.DefaultCompression))}
	// SnappyCodec snappy 压缩
	SnappyCodec = snappyCodec{newBaseCodec(xcompress.WithType(xcompress.Snappy))}
)

var codecs = map[Type]xencoding.Codec{
	NoneType:   NoneCodec,
	GZIPType:   GZIPCodec,
	SnappyType: SnappyCodec,
}

func init() {
	for _, v := range codecs {
		xencoding.RegisterCodec(v)
	}
}

// Register 注册自定义的解压缩 Codec ，该方法非协程安全
func Register(t Type, codec xencoding.Codec) {
	_, exists := codecs[t]
	xpanic.WhenTrue(exists, "register called twice for codec, %d", t)
	codecs[t] = codec
}

type baseCodec struct {
	xcompress.Compressor
}

func newBaseCodec(opts ...xcompress.Option) baseCodec {
	return baseCodec{Compressor: xcompress.MustNew(opts...)}
}

func (c baseCodec) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errCodecMarshalParam
	} else {
		return c.Compressor.Flat(data)
	}
}

func (c baseCodec) Unmarshal(bytes []byte, v interface{}) error {
	v1, ok := v.(*[]byte)
	if !ok {
		return errCodecUnmarshalParam
	}
	data, err := c.Compressor.Inflate(bytes)
	if err != nil {
		return err
	}
	*v1 = data
	return nil
}

type (
	noneCodec   struct{ baseCodec }
	gzipCodec   struct{ baseCodec }
	snappyCodec struct{ baseCodec }
)

func (c noneCodec) Name() string   { return NoneCodecName }
func (c gzipCodec) Name() string   { return GZIPCodecName }
func (c snappyCodec) Name() string { return SnappyCodecName }

type codec struct {
	_type Type
}

// NewCodec 通过压缩类型创建压缩 Codec
func NewCodec(compressType Type) xencoding.Codec {
	return codec{_type: compressType}
}

// Name 返回 Codec 名
func (c codec) Name() string { return "compressor" }

// Marshal 编码
func (c codec) Marshal(v interface{}) ([]byte, error) {
	cc, ok1 := codecs[c._type]
	if !ok1 {
		return nil, errCodecNoFound
	}
	return cc.Marshal(v)
}

// Unmarshal 解码
func (c codec) Unmarshal(bytes []byte, v interface{}) error {
	cc, ok1 := codecs[c._type]
	if !ok1 {
		return errCodecNoFound
	}
	return cc.Unmarshal(bytes, v)
}
