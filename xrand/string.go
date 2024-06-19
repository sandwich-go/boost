package xrand

import (
	"fmt"
	"strings"
	"time"
	"unsafe"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyz"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var nowFunc = time.Now

// StringWithTimestamp 随机 n 个字符的字符串，并以当前时间戳做后缀
func StringWithTimestamp(n int, letterList ...string) string {
	return fmt.Sprintf("%s_%d", String(n, letterList...), nowFunc().Unix())
}

// String 随机 n 个字符的字符串
func String(n int, letterList ...string) string {
	letterBytesUsing := strings.Join(letterList, "")
	if letterBytesUsing == "" {
		letterBytesUsing = letterBytes
	}
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, FastRand(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = FastRand(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytesUsing) {
			b[i] = letterBytesUsing[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
