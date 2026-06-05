package fnv

import (
	"bytes"
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
		// 直接 buf.Write 即可：[]byte 的 binary.Write 等价于按字节顺序写入，
		// 不需要 LittleEndian 转换。bytes.Buffer.Write 永远不返回 err。
		buf.WriteString(v)
	default:
		panic("unsupported type")
	}

	return fnv1aHash(buf.Bytes())
}
