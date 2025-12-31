package pbjson

import (
	"context"
	"errors"

	"github.com/sandwich-go/boost/xencoding"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

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
	marshaler   = &protojson.MarshalOptions{UseEnumNumbers: true}
	unmarshaler = &protojson.UnmarshalOptions{}
)

// EmitUnpopulated 指定是否使用零值渲染字段
func EmitUnpopulated(emit bool) { marshaler.EmitUnpopulated = emit }

// UseEnumNumbers 设置是否将 enum 序列化为数字，默认开启功能
func UseEnumNumbers(b bool) { marshaler.UseEnumNumbers = b }

func init() {
	xencoding.RegisterCodec(Codec)
}

// codec is a Codec implementation with json
type codec struct{}

// Name 返回 Codec 名
func (codec) Name() string { return CodecName }

// Marshal 编码
func (codec) Marshal(_ context.Context, obj interface{}) ([]byte, error) {
	if pm, ok := obj.(proto.Message); ok {
		buf, err := marshaler.Marshal(pm)
		return buf, err
	}
	return nil, errCodecParam
}

// Unmarshal 解码
func (codec) Unmarshal(_ context.Context, data []byte, v interface{}) error {
	if pm, ok := v.(proto.Message); ok {
		err := unmarshaler.Unmarshal(data, pm)
		if err == io.EOF {
			err = nil
		}
		return err
	}
	return errCodecParam
}
