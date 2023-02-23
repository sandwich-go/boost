package xcompress

import (
	"fmt"
	"github.com/golang/snappy"
	"github.com/sandwich-go/boost/xpanic"
)

type snappyCompressor struct{}

func newSnappyCompressor() (Compressor, error) {
	return &snappyCompressor{}, nil
}

func (c *snappyCompressor) Flat(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	var err error
	var out []byte
	xpanic.Try(func() {
		out = snappy.Encode(nil, data)
	}).Catch(func(e xpanic.E) {
		err = fmt.Errorf("snappy flat error, %v", e)
	})
	return out, err
}

func (c *snappyCompressor) Inflate(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	return snappy.Decode(nil, data)
}
