package compressor

import (
	"bytes"
	"github.com/golang/snappy"
	"io/ioutil"
)

type snappyCompressor struct{}

func newSnappyCompressor() (Compressor, error) {
	return &snappyCompressor{}, nil
}

func (c *snappyCompressor) Flat(data []byte) ([]byte, error) {
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

func (c *snappyCompressor) Inflate(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	return ioutil.ReadAll(snappy.NewReader(bytes.NewReader(data)))
}
