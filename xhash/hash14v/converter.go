package hash14v

import (
	"bytes"
	"github.com/sandwich-go/boost/xcrypto/algorithm/vigenere"
	"github.com/sandwich-go/boost/z"
	"math"
)

const (
	base          = 26
	maxLength     = 14
	minHashOffset = "A"
)

type converter struct {
	reservedBuff [16]byte
	idOffset     uint64
	hashOffset   []byte
	hashKey      []byte
	minLength    int
	maxLength    int
	spec         *Options
}

// New return new hash14v operator
func New(opts ...Option) Converter {
	spec := NewOptions(opts...)
	oo := &converter{spec: spec}
	oo.hashKey = vigenere.Sanitize(spec.GetHashKey())
	oo.hashOffset = spec.GetHashOffset()
	if len(oo.hashOffset) == 0 {
		oo.hashOffset = z.StringToBytes(minHashOffset)
	}
	oo.minLength = len(oo.hashOffset)
	oo.maxLength = maxLength
	oo.idOffset = oo.decode(oo.hashOffset)
	return oo
}

// ToV 对 Id 进行转换，转换为 V
func (o *converter) ToV(id Id) V {
	if math.MaxUint64-id < o.idOffset {
		return nil
	}
	ret := o.encode(id + o.idOffset)
	if o.hashKey != nil {
		vigenere.EncryptAndInplace(ret, o.hashKey)
	}
	return ret
}

// ToId 对 Id 进行转换，转换为 V
func (o *converter) ToId(v V) Id {
	if len(v) < o.minLength || len(v) > o.maxLength {
		return 0
	}
	if o.hashKey == nil {
		return o.decode(v) - o.idOffset
	}
	return o.decode(vigenere.Decrypt(v, o.hashKey)) - o.idOffset
}

// Offset 获取 Id 的 offset
func (o *converter) Offset() Id { return o.idOffset }

func (o *converter) encode(num uint64) V {
	if o.spec.GetUsingReservedBuff() {
		var i int
		for i = len(o.reservedBuff) - 1; num != 0; i-- {
			o.reservedBuff[i] = 'A' + byte(num%base) - 1
			num /= 26
		}
		return o.reservedBuff[i+1:]
	} else {
		var rb [16]byte
		var i int
		for i = len(rb) - 1; num != 0; i-- {
			rb[i] = 'A' + byte(num%base) - 1
			num /= 26
		}
		return rb[i+1:]
	}
}

func (o *converter) decode(hi V) (ret uint64) {
	hi = bytes.ToUpper(hi)
	lenID := len(hi)
	for index, c := range hi {
		ret += uint64(math.Pow(base, float64(lenID-1-index))) * uint64(c-'A'+1)
	}
	return ret
}
