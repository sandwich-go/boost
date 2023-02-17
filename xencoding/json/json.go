package json

import (
	"encoding/json"
	"github.com/sandwich-go/boost/xencoding"
)

var Codec = &codec{}

const (
	// CodecName json 加解码名称，可以通过 encoding2.GetCodec(CodecName) 获取对应的 Codec
	CodecName = "json"
)

func init() {
	xencoding.RegisterCodec(Codec)
}

// codec is a Codec implementation with json
type codec struct{}

// Name 返回 Codec 名
func (codec) Name() string { return CodecName }

// Marshal 编码
func (codec) Marshal(v interface{}) ([]byte, error) { return json.Marshal(v) }

// Unmarshal 解码
func (codec) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }
