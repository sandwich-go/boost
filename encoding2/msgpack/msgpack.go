package json

import (
	"github.com/sandwich-go/boost/encoding2"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

var Codec = &codec{}

const (
	// CodecName msgpack 加解码名称，可以通过 encoding2.GetCodec(CodecName) 获取对应的 Codec
	CodecName = "msgpack"
)

func init() {
	encoding2.RegisterCodec(Codec)
}

// codec is a Codec implementation with json
type codec struct{}

// Name 返回 Codec 名
func (codec) Name() string { return CodecName }

// Marshal 编码
func (codec) Marshal(v interface{}) ([]byte, error) { return msgpack.Marshal(v) }

// Unmarshal 解码
func (codec) Unmarshal(data []byte, v interface{}) error { return msgpack.Unmarshal(data, v) }
