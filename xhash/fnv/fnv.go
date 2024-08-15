package fnv

import (
	"bytes"
	"encoding/binary"
)

const (
	fnvOffsetBasis uint64 = 14695981039346656037
	fnvPrime       uint64 = 1099511628211
)

func fnv1aHash(data []byte) uint64 {
	hash := fnvOffsetBasis
	for _, b := range data {
		hash ^= uint64(b)
		hash *= fnvPrime
	}
	return hash
}

func Hash(value interface{}) uint64 {
	var buf bytes.Buffer

	switch v := value.(type) {
	case int:
		if v < 0 {
			return uint64(-v)
		}
		return (uint64)(v)
	case uint:
		return (uint64)(v)
	case int32:
		if v < 0 {
			return uint64(-v)
		}
		return (uint64)(v)
	case uint32:
		return (uint64)(v)
	case int64:
		if v < 0 {
			return uint64(-v)
		}
		return (uint64)(v)
	case uint64:
		return (uint64)(v)
	case float32:
		if v < 0 {
			return uint64(-v)
		}
		return (uint64)(v)
	case float64:
		if v < 0 {
			return uint64(-v)
		}
		return (uint64)(v)
	case string:
		data := []byte(v)
		binary.Write(&buf, binary.LittleEndian, data)
	default:
		panic("unsupported type")
	}

	return fnv1aHash(buf.Bytes())
}
