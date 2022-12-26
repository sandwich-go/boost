package encrypt

import (
	"errors"
	"github.com/sandwich-go/boost/encoding2"
)

var (
	errEncryptCodecMarshalParam   = errors.New("encrypt codec marshal must be []byte parameter")
	errEncryptCodecUnmarshalParam = errors.New("encrypt codec unmarshal must be *[]byte parameter")
	errEncryptCodecNoFound        = errors.New("encrypt codec not found")
)

type Type byte

const (
	None Type = iota
	AES
)

var codecs = map[Type]encoding2.Codec{
	None: dummyEncryptCodec{},
	AES:  aesEncryptCodec{},
}

type keySetter interface {
	SetKey([]byte) encoding2.Codec
}

type Codec struct {
	key         []byte
	encryptType Type
}

func NewCodec(encryptType Type, key []byte) encoding2.Codec {
	return Codec{encryptType: encryptType, key: key}
}

func (c Codec) Name() string { return "encrypt" }
func (c Codec) Marshal(v interface{}) ([]byte, error) {
	cc, ok1 := codecs[c.encryptType]
	if !ok1 {
		return nil, errEncryptCodecNoFound
	}
	return cc.(keySetter).SetKey(c.key).Marshal(v)
}
func (c Codec) Unmarshal(bytes []byte, v interface{}) error {
	cc, ok1 := codecs[c.encryptType]
	if !ok1 {
		return errEncryptCodecNoFound
	}
	return cc.(keySetter).SetKey(c.key).Unmarshal(bytes, v)
}
