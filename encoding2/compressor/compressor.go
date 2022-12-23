package compressor

import (
	"errors"
	"github.com/sandwich-go/boost/encoding2"
)

var (
	errCompressorCodecMarshalParam   = errors.New("compressor codec marshal must be []byte parameter")
	errCompressorCodecUnmarshalParam = errors.New("compressor codec unmarshal must be *[]byte parameter")
	errCompressorCodecNoFound        = errors.New("compressor codec not found")
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
