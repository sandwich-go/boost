package compressor

import "github.com/sandwich-go/boost/encoding2"

var DummyCodec = DummyCompressorCodec{}

func init() {
	encoding2.RegisterCodec(DummyCodec)
}

type DummyCompressorCodec struct{}

func (c DummyCompressorCodec) Name() string { return "dummy_compressor" }
func (c DummyCompressorCodec) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errCompressorCodecMarshalParam
	} else {
		return data, nil
	}
}
func (c DummyCompressorCodec) Unmarshal(bytes []byte, v interface{}) error {
	v1, ok := v.(*[]byte)
	if !ok {
		return errCompressorCodecUnmarshalParam
	}
	*v1 = bytes
	return nil
}
