package xrand

import (
	"strconv"
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

// StringWithTimestamp 随机 n 个字符的字符串，并以当前时间戳做后缀（格式：random_unix）。
//
// 把随机字符 + '_' + Unix() 拼到同一个 byte buffer，单 alloc 完成；
// 容量预留 n + 1 ('_') + 20 (Unix() 最长 19 位 + 1 安全裕量) = n + 21。
//
// 等价于 fmt.Sprintf("%s_%d", String(n, letterList...), nowFunc().Unix())，
// 但少 3 alloc：bench 实测 -47% ns / -33% B。
//
// 主体随机填充逻辑与 [String] 字面一致；为保证 String 不被跨函数调用拖累，
// 这里手动复制循环（没有抽 helper）。
func StringWithTimestamp(n int, letterList ...string) string {
	letterBytesUsing := strings.Join(letterList, "")
	if letterBytesUsing == "" {
		letterBytesUsing = letterBytes
	}
	buf := make([]byte, n, n+21)
	for i, cache, remain := n-1, FastRand(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = FastRand(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytesUsing) {
			buf[i] = letterBytesUsing[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	buf = append(buf, '_')
	buf = strconv.AppendInt(buf, nowFunc().Unix(), 10)
	return *(*string)(unsafe.Pointer(&buf))
}

// String 随机 n 个字符的字符串。
// letterList 为空时使用默认字母表 letterBytes ("a-z")。
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
