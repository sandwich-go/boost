package protobuf

import (
	"bytes"
	"errors"
	"github.com/sandwich-go/boost/xencoding"
	"github.com/sandwich-go/boost/xerror"
	"math"
	"reflect"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
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
func (p codec) Marshal(v interface{}) ([]byte, error) {
	if pm, ok := v.(proto.Marshaler); ok {
		// object can marshal itself, no need for buffer
		return pm.Marshal()
	}
	if pm, ok := v.(proto.Message); ok {
		if p.usingPool {
			cb := protoBufferPool.Get().(*cachedProtoBuffer)
			out, err := marshal(pm, cb)
			// put back buffer and lose the ref to the slice
			cb.SetBuf(nil)
			protoBufferPool.Put(cb)
			return out, err
		}
		return proto.Marshal(pm)
	}
	return nil, xerror.NewText("%T is not a proto.Marshaler", v)
}

// Uri 获取 Message Name
func (codec) Uri(t interface{}) string { return proto.MessageName(t.(proto.Message)) }

// Type 获取 Message Type
func (codec) Type(uri string) reflect.Type { return proto.MessageType(uri) }

// Unmarshal 解码
func (p codec) Unmarshal(data []byte, v interface{}) error {
	if pu, ok := v.(proto.Unmarshaler); ok {
		// object can unmarshal itself, no need for buffer
		return pu.Unmarshal(data)
	}

	if m, ok := v.(proto.Message); ok {
		m.Reset()
		if p.usingPool {
			cb := protoBufferPool.Get().(*cachedProtoBuffer)
			cb.SetBuf(data)
			err := cb.Unmarshal(m)
			cb.SetBuf(nil)
			protoBufferPool.Put(cb)
			return err
		}
		return proto.Unmarshal(data, m)
	}

	return xerror.NewText("%T is not a proto.Unmarshaler", v)
}

func (codec) JSONMarshal(obj interface{}) ([]byte, error) {
	if pm, ok := obj.(proto.Message); ok {
		m := jsonpb.Marshaler{EmitDefaults: false}
		var buf bytes.Buffer
		return buf.Bytes(), m.Marshal(&buf, pm)
	}
	return nil, errors.New("not proto message")
}

func marshal(pm proto.Message, cb *cachedProtoBuffer) ([]byte, error) {
	newSlice := make([]byte, 0, cb.lastMarshaledSize)

	cb.SetBuf(newSlice)
	cb.Reset()
	if err := cb.Marshal(pm); err != nil {
		return nil, err
	}
	out := cb.Bytes()
	cb.lastMarshaledSize = capToMaxInt32(len(out))
	return out, nil
}

func capToMaxInt32(val int) uint32 {
	if val > math.MaxInt32 {
		return uint32(math.MaxInt32)
	}
	return uint32(val)
}

type cachedProtoBuffer struct {
	lastMarshaledSize uint32
	proto.Buffer
}

var protoBufferPool = &sync.Pool{
	New: func() interface{} {
		return &cachedProtoBuffer{
			Buffer:            proto.Buffer{},
			lastMarshaledSize: 16,
		}
	},
}
