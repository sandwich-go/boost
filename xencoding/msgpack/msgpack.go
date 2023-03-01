package json

import (
	"context"
	"github.com/sandwich-go/boost/xencoding"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

var Codec = &codec{}

const (
	// CodecName msgpack 加解码名称，可以通过 encoding2.GetCodec(CodecName) 获取对应的 Codec
	CodecName = "msgpack"
)

func init() {
	xencoding.RegisterCodec(Codec)
}

// codec is a Codec implementation with json
type codec struct{}

// Name 返回 Codec 名
func (codec) Name() string { return CodecName }

// Marshal 编码
func (codec) Marshal(_ context.Context, v interface{}) ([]byte, error) { return msgpack.Marshal(v) }

// Unmarshal 解码
func (codec) Unmarshal(_ context.Context, data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
