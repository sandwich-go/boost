package pbjson

import (
	"bytes"
	"errors"
	"io"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/sandwich-go/boost/encoding2"
)

var Codec = &jsonCodec{}

var marshaler = &jsonpb.Marshaler{
	EnumsAsInts: true,
}
var unMarshaler = &jsonpb.Unmarshaler{}

func SetMarshalerEmitDefaults(emit bool) {
	marshaler.EmitDefaults = emit
}

//SetMarshalerEnumsAsInts 设置marsher是否将enum序列化为数字，默认开启功能
func SetMarshalerEnumsAsInts(b bool) {
	marshaler.EnumsAsInts = b
}

func init() {
	encoding2.RegisterCodec(Codec)
}

// jsonCodec is a Codec implementation with json
type jsonCodec struct{}

func (jsonCodec) Marshal(obj interface{}) ([]byte, error) {
	if pm, ok := obj.(proto.Message); ok {
		var buf bytes.Buffer
		err := marshaler.Marshal(&buf, pm)
		return buf.Bytes(), err
	}
	return nil, errors.New("not proto message")
}

func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	if pm, ok := v.(proto.Message); ok {
		err := unMarshaler.Unmarshal(bytes.NewBuffer(data), pm)
		if err == io.EOF {
			err = nil
		}
		return err
	}
	return errors.New("not proto message")
}

func (jsonCodec) Name() string {
	return "pbjson"
}
