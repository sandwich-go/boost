package z

import "unsafe"

// StringToBytes converts string to byte slice without a memory allocation.
//
// 历史实现用 reflect.StringHeader/SliceHeader，Go 1.20 起 deprecated 且
// vet 会报 "possible misuse"。Go 1.21+ 推荐用 unsafe.StringData / unsafe.Slice
// 显式表达"借用 string 底层 bytes，不拷贝"的语义。
//
// 注意：返回的 []byte 与 s 共享底层数据，**禁止写入**——string 的字节是
// immutable 内存，写入会触发 SIGSEGV 或破坏其他 string。
func StringToBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts byte slice to string without a memory allocation.
//
// 注意：返回的 string 与 b 共享底层数据，调用方在持有该 string 期间
// **不得修改 b**——否则 string immutable 不变量被打破。
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}
