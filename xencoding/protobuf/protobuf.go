package protobuf

import (
	"context"
	"errors"
	"reflect"

	"github.com/sandwich-go/boost/xencoding"
	"github.com/sandwich-go/boost/xerror"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	Codec          = codec{usingPool: false, name: CodecName}
	CodecUsingPool = codec{usingPool: true, name: UsingPoolCodecName}
)

const (
	// CodecName proto 压缩效果名称，可以通过 encoding2.GetCodec(CodecName) 获取对应的 Codec
	CodecName = "proto"
	// UsingPoolCodecName 带对象池的 proto 压缩效果名称，可以通过 encoding2.GetCodec(UsingPoolCodecName) 获取对应的 Codec
	UsingPoolCodecName = "proto_using_pool"
)

func init() {
	xencoding.RegisterCodec(Codec)
	xencoding.RegisterCodec(CodecUsingPool)
}

// codec is a Codec implementation with protobuf. It is the default codec.
type codec struct {
	usingPool bool
	name      string
}

// Name 返回 Codec 名
func (p codec) Name() string { return p.name }

// Marshal 编码
func (p codec) Marshal(_ context.Context, v interface{}) ([]byte, error) {
	if pm, ok := v.(proto.Message); ok {
		return proto.Marshal(pm)
	}
	return nil, xerror.NewText("%T is not a proto.Message", v)
}

// Uri 获取 Message Name
func (codec) Uri(t interface{}) string {
	return string(t.(proto.Message).ProtoReflect().Descriptor().FullName())
}

// Type 获取 Message Type
func (codec) Type(uri string) reflect.Type {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(uri))
	if err != nil {
		return nil
	}
	return reflect.TypeOf(mt.Zero().Interface())
}

// Unmarshal 解码
func (p codec) Unmarshal(ctx context.Context, data []byte, v interface{}) error {
	if m, ok := v.(proto.Message); ok {
		proto.Reset(m)
		return proto.Unmarshal(data, m)
	}

	return xerror.NewText("%T is not a proto.Message", v)
}

func (codec) JSONMarshal(obj interface{}) ([]byte, error) {
	if pm, ok := obj.(proto.Message); ok {
		return protojson.MarshalOptions{
			EmitUnpopulated: false,
		}.Marshal(pm)
	}
	return nil, errors.New("not proto message")
}

