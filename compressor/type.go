package compressor

//go:generate stringer -type=Type
type Type int

const (
	Dummy  Type = iota // 不使用任何解压缩
	GZIP               // gzip 解压缩
	Snappy             // snappy 解压缩
)
