package encrypt

import "github.com/sandwich-go/boost/encoding2"

type dummyEncryptCodec struct{}

func (c dummyEncryptCodec) Name() string                  { return "dummy_encrypt" }
func (c dummyEncryptCodec) SetKey([]byte) encoding2.Codec { return c }
func (c dummyEncryptCodec) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errEncryptCodecMarshalParam
	} else {
		return data, nil
	}
}
func (c dummyEncryptCodec) Unmarshal(bytes []byte, v interface{}) error {
	v1, ok := v.(*[]byte)
	if !ok {
		return errEncryptCodecUnmarshalParam
	}
	*v1 = bytes
	return nil
}
