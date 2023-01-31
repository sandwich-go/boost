package compressor

//go:generate stringer -type=Type
type Type byte

const (
	NoneType   Type = iota // 无压缩效果
	GZIPType               // gzip 压缩
	SnappyType             // snappy压缩
)
