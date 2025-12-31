package protobuf

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/sandwich-go/boost/xencoding"
	"github.com/sandwich-go/boost/xerror"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// vtprotoMessage 定义了 vtproto 生成的序列化方法接口
type vtprotoMessage interface {
	MarshalVT() ([]byte, error)
	SizeVT() int
	MarshalToSizedBufferVT(dAtA []byte) (int, error)
}

// vtprotoUnmarshaler 定义了 vtproto 生成的反序列化方法接口
type vtprotoUnmarshaler interface {
	UnmarshalVT(dAtA []byte) error
}

var (
	Codec          = codec{usingPool: false, name: CodecName}
	CodecUsingPool = codec{usingPool: true, name: UsingPoolCodecName}
)

const (
	// CodecName proto 压缩效果名称，可以通过 encoding2.GetCodec(CodecName) 获取对应的 Codec
	CodecName = "proto"
	// UsingPoolCodecName 带对象池的 proto 压缩效果名称，可以通过 encoding2.GetCodec(UsingPoolCodecName) 获取对应的 Codec
	// bench 测试显示 提升有限
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
		// 优先使用 vtproto 的 MarshalVT 方法
		if vt, ok := v.(vtprotoMessage); ok {
			if p.usingPool {
				return marshalVTWithPool(vt)
			}
			return vt.MarshalVT()
		}

		if p.usingPool {
			return standardMarshalWithPool(pm)
		}
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
		// 优先使用 vtproto 的 UnmarshalVT 方法
		if vt, ok := v.(vtprotoUnmarshaler); ok {
			return vt.UnmarshalVT(data)
		}

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

func standardMarshalWithPool(pm proto.Message) ([]byte, error) {
	op := proto.MarshalOptions{UseCachedSize: true}
	size := op.Size(pm)
	bufPtr := bufPool.Get().(*[]byte)
	buf := *bufPtr

	// 确保容量足够
	if cap(buf) < size {
		buf = make([]byte, 0, size)
	}

	buf, err := op.MarshalAppend(buf[:0], pm)
	if err != nil {
		return nil, err
	}
	ret := make([]byte, len(buf))
	copy(ret, buf)
	*bufPtr = buf
	bufPool.Put(bufPtr)
	return ret, nil
}

func marshalVTWithPool(vt vtprotoMessage) ([]byte, error) {
	size := vt.SizeVT()
	bufPtr := bufPool.Get().(*[]byte)
	buf := *bufPtr

	// 确保容量足够
	if cap(buf) < size {
		buf = make([]byte, size)
	} else {
		buf = buf[:size]
	}

	n, err := vt.MarshalToSizedBufferVT(buf)
	if err != nil {
		*bufPtr = buf
		bufPool.Put(bufPtr)
		return nil, err
	}

	// MarshalToSizedBufferVT 从后往前写，所以数据在 buf[len(buf)-n:]
	ret := make([]byte, n)
	copy(ret, buf[len(buf)-n:])
	*bufPtr = buf
	bufPool.Put(bufPtr)
	return ret, nil
}

var bufPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 0, 512)
		return &buf
	},
}
