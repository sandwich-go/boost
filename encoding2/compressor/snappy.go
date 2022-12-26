package compressor

import (
	"bytes"
	"github.com/golang/snappy"
	"github.com/sandwich-go/boost/encoding2"
	"io/ioutil"
)

var SnappyCodec = SnappyCompressorCodec{}

func init() {
	encoding2.RegisterCodec(SnappyCodec)
}

type SnappyCompressorCodec struct{}

func (c SnappyCompressorCodec) Name() string { return "snappy_compressor" }
func (c SnappyCompressorCodec) Marshal(v interface{}) ([]byte, error) {
	data, ok := v.([]byte)
	if !ok {
		return nil, errCompressorCodecMarshalParam
	}
	if len(data) == 0 {
		return data, nil
	}
	var buffer bytes.Buffer
	writer := snappy.NewBufferedWriter(&buffer)
	_, err := writer.Write(data)
	if err != nil {
		_ = writer.Close()
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
func (c SnappyCompressorCodec) Unmarshal(data []byte, v interface{}) error {
	v1, ok := v.(*[]byte)
	if !ok {
		return errCompressorCodecUnmarshalParam
	}
	if len(data) == 0 {
		*v1 = data
		return nil
	}
	reader := snappy.NewReader(bytes.NewReader(data))
	out, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	*v1 = out
	return nil
}
