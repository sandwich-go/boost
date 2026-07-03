package smap

import "github.com/sandwich-go/boost/z"

// KeyToHashAny smap只是利用KeyToHash来确定单次运行过程中shared计算，并不在意hash冲突，如需持久化固定的hash则选择其他算法
var KeyToHashAny = z.KeyToHash
var KeyToHashBytes = z.MemHash
var KeyToHashString = z.MemHashString

func KeyToHashInt(k int) uint64     { return uint64(k) }
func KeyToHashInt32(k int32) uint64 { return uint64(k) }
func KeyToHashInt64(k int64) uint64 { return uint64(k) }

var _ = KeyToHashAny
var _ = KeyToHashBytes
var _ = KeyToHashString
var _ = KeyToHashInt
var _ = KeyToHashInt32
var _ = KeyToHashInt64
