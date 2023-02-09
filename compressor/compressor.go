package compressor

import "compress/flate"

var (
	NoCompression      = flate.NoCompression
	BestSpeed          = flate.BestSpeed
	BestCompression    = flate.BestCompression
	DefaultCompression = flate.DefaultCompression
	HuffmanOnly        = flate.HuffmanOnly
)

// Compressor 压缩器
type Compressor interface {
	// Flat 压缩二进制数据
	Flat([]byte) ([]byte, error)
	// Inflate 解压二进制数据
	Inflate([]byte) ([]byte, error)
}

// Default 默认的压缩器，使用 GZIP 方式，进行 DefaultCompression 等级的压缩与解压缩
var Default = MustNew(WithType(GZIP), WithLevel(DefaultCompression))
