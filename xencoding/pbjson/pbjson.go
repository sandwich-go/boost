package pbjson

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/sandwich-go/boost/xencoding"
	"io"
)

var (
	errCodecParam = errors.New("pbjson codec marshal/unmarshal must be proto message")
)

const (
	// CodecName pbjson 加解码名称，可以通过 encoding2.GetCodec(CodecName) 获取对应的 Codec
	CodecName = "pbjson"
)

var Codec = codec{}

var (
	marshaler   = &jsonpb.Marshaler{EnumsAsInts: true}
	unmarshaler = &jsonpb.Unmarshaler{}
)

// EmitUnpopulated 指定是否使用零值渲染字段
func EmitUnpopulated(emit bool) { marshaler.EmitDefaults = emit }

// UseEnumNumbers 设置是否将 enum 序列化为数字，默认开启功能
func UseEnumNumbers(b bool) { marshaler.EnumsAsInts = b }

func init() {
	xencoding.RegisterCodec(Codec)
}

// codec is a Codec implementation with json
type codec struct{}

// Name 返回 Codec 名
func (codec) Name() string { return CodecName }

// Marshal 编码
func (codec) Marshal(obj interface{}) ([]byte, error) {
	if pm, ok := obj.(proto.Message); ok {
		var buf bytes.Buffer
		err := marshaler.Marshal(&buf, pm)
		return buf.Bytes(), err
	}
	return nil, errCodecParam
}

// Unmarshal 解码
func (codec) Unmarshal(data []byte, v interface{}) error {
	if pm, ok := v.(proto.Message); ok {
		err := unmarshaler.Unmarshal(bytes.NewBuffer(data), pm)
		if err == io.EOF {
			err = nil
		}
		return err
	}
	return errCodecParam
}
