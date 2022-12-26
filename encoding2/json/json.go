package json

import (
	"encoding/json"

	"github.com/sandwich-go/boost/encoding2"
)

var Codec = &jsonCodec{}

func init() {
	encoding2.RegisterCodec(Codec)
}

// jsonCodec is a Codec implementation with json
type jsonCodec struct{}

func (jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (jsonCodec) Name() string {
	return "json"
}
