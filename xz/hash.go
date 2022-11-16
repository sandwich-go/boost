package xz

import (
	"fmt"
	"unsafe"
)

type stringStruct struct {
	str unsafe.Pointer
	len int
}

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

// MemHash 不同进程内同一个值获取的可能不同，不能用作持久化的hash
func MemHash(data []byte) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&data))
	return uint64(memhash(ss.str, 0, uintptr(ss.len)))
}

// MemHashString 不同进程内同一个值获取的可能不同，不能用作持久化的hash
func MemHashString(str string) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&str))
	return uint64(memhash(ss.str, 0, uintptr(ss.len)))
}

// KeyToHash
// TODO: Figure out a way to re-use memhash for the second uint64 hash, we
//       already know that appending bytes isn't reliable for generating a
//       second hash (see Ristretto PR #88).
//
//       We also know that while the Go runtime has a runtime memhash128
//       function, it's not possible to use it to generate [2]uint64 or
//       anything resembling a 128bit hash, even though that's exactly what
//       we need in this situation.

type HashUint64 interface{ HashUint64() uint64 }

func KeyToHash(key interface{}) uint64 {
	if key == nil {
		return 0
	}
	switch k := key.(type) {
	case uint64:
		return k
	case HashUint64:
		return k.HashUint64()
	case fmt.Stringer:
		return MemHashString(k.String())
	case string:
		return MemHashString(k)
	case []byte:
		return MemHash(k)
	case byte:
		return uint64(k)
	case int:
		return uint64(k)
	case int8:
		return uint64(k)
	case int16:
		return uint64(k)
	case uint16:
		return uint64(k)
	case int32:
		return uint64(k)
	case uint32:
		return uint64(k)
	case int64:
		return uint64(k)
	default:
		panic("Key type not supported")
	}
}
