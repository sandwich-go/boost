package compressor

import (
	"github.com/sandwich-go/boost/compressor"
	"github.com/sandwich-go/boost/encoding2"
)

var GzipCodec = GzipCompressorCodec{}

func init() {
	encoding2.RegisterCodec(GzipCodec)
}

type GzipCompressorCodec struct{}

func (c GzipCompressorCodec) Name() string { return "gzip_compressor" }
func (c GzipCompressorCodec) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errCompressorCodecMarshalParam
	} else {
		return compressor.Zip(data)
	}
}

func (c GzipCompressorCodec) Unmarshal(bytes []byte, v interface{}) error {
	v1, ok := v.(*[]byte)
	if !ok {
		return errCompressorCodecUnmarshalParam
	}
	data, err := compressor.Unzip(bytes)
	if err != nil {
		return err
	}
	*v1 = data
	return nil
}
