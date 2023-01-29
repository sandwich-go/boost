package encrypt

import (
	"github.com/sandwich-go/boost/encoding2"
	"github.com/sandwich-go/boost/encryption/algorithm/aes"
)

type aesEncryptCodec struct {
	key []byte
}

func (c aesEncryptCodec) Name() string { return "aes_encrypt" }
func (c aesEncryptCodec) SetKey(k []byte) encoding2.Codec {
	c.key = k
	return c
}
func (c aesEncryptCodec) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); !ok {
		return nil, errEncryptCodecMarshalParam
	} else {
		return aes.Encrypt(data, c.key)
	}
}

func (c aesEncryptCodec) Unmarshal(bytes []byte, v interface{}) error {
	v1, ok := v.(*[]byte)
	if !ok {
		return errEncryptCodecUnmarshalParam
	}
	data, err := aes.Decrypt(bytes, c.key)
	if err != nil {
		return err
	}
	*v1 = data
	return nil
}
