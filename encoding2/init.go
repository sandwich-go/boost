package encoding2

import (
	"github.com/sandwich-go/boost/encoding2/compressor"
)

func GetCompressor(name string) compressor.Compressor {
	return compressor.GetCompressor(name)
}
