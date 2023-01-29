package compressor

import (
	"errors"
	"github.com/sandwich-go/boost/compressor"
	"github.com/sandwich-go/boost/encoding2"
)

var (
	errCompressorCodecMarshalParam   = errors.New("compressor codec marshal must be []byte parameter")
	errCompressorCodecUnmarshalParam = errors.New("compressor codec unmarshal must be *[]byte parameter")
	errCompressorCodecNoFound        = errors.New("compressor codec not found")
)

var (
	DummyCodec  = DummyCompressorCodec{newBaseCompressorCodec(compressor.WithType(compressor.Dummy))}
	SnappyCodec = SnappyCompressorCodec{newBaseCompressorCodec(compressor.WithType(compressor.Snappy))}
	GzipCodec   = GzipCompressorCodec{newBaseCompressorCodec(compressor.WithType(compressor.GZIP), compressor.WithLevel(compressor.DefaultCompression))}
)

type CompressType byte

const (
	CompressNone CompressType = iota
	CompressGzip
	CompressSnappy
)

var compressorCodecs = map[CompressType]encoding2.Codec{
	CompressNone:   DummyCodec,
	CompressGzip:   GzipCodec,
	CompressSnappy: SnappyCodec,
}

func init() {
	for _, v := range compressorCodecs {
		encoding2.RegisterCodec(v)
	}
}

type baseCompressorCodec struct {
	compressor.Compressor
}

func newBaseCompressorCodec(opts ...compressor.Option) baseCompressorCodec {
	return baseCompressorCodec{Compressor: compressor.MustNew(opts...)}
}

func (c baseCompressorCodec) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errCompressorCodecMarshalParam
	} else {
		return c.Compressor.Flat(data)
	}
}

func (c baseCompressorCodec) Unmarshal(bytes []byte, v interface{}) error {
	v1, ok := v.(*[]byte)
	if !ok {
		return errCompressorCodecUnmarshalParam
	}
	data, err := c.Compressor.Inflate(bytes)
	if err != nil {
		return err
	}
	*v1 = data
	return nil
}

type (
	DummyCompressorCodec  struct{ baseCompressorCodec }
	GzipCompressorCodec   struct{ baseCompressorCodec }
	SnappyCompressorCodec struct{ baseCompressorCodec }
)

func (c DummyCompressorCodec) Name() string  { return "dummy_compressor" }
func (c GzipCompressorCodec) Name() string   { return "gzip_compressor" }
func (c SnappyCompressorCodec) Name() string { return "snappy_compressor" }

type Codec struct {
	compressType CompressType
}

func NewCodec(compressType CompressType) encoding2.Codec {
	return Codec{compressType: compressType}
}

func (c Codec) Name() string { return "compressor" }
func (c Codec) Marshal(v interface{}) ([]byte, error) {
	cc, ok1 := compressorCodecs[c.compressType]
	if !ok1 {
		return nil, errCompressorCodecNoFound
	}
	return cc.Marshal(v)
}
func (c Codec) Unmarshal(bytes []byte, v interface{}) error {
	cc, ok1 := compressorCodecs[c.compressType]
	if !ok1 {
		return errCompressorCodecNoFound
	}
	return cc.Unmarshal(bytes, v)
}
