package isnil

import (
	"unsafe"
)

// eface 与 runtime.eface 内存布局一致，用于直接读取 interface{} 的数据字。
type eface struct {
	rtype unsafe.Pointer
	data  unsafe.Pointer
}

// Check checks if any is nil.
// 判定规则：接口本身为 nil，或接口的数据字为 nil（即 typed nil，如 (*T)(nil)）。
// 注意 nil slice/map/chan 装箱后数据字非 nil（runtime.convTslice 等会指向 zeroVal），
// 因此对它们返回 false，与 == nil 的直接比较语义不同，调用方需自行区分。
//
// 实现约束：unsafe.Pointer(&any) 派生的指针必须在本函数内解引用，不得作为返回值或
// 参数传出。一旦流出，逃逸分析只能把 any 判为逃逸，导致每次调用都堆分配一个 eface
// （16B/次，即便调用方传的已经是指针或已装箱的 interface{}）。
func Check(any interface{}) bool {
	if any == nil {
		return true
	}
	return (*eface)(unsafe.Pointer(&any)).data == nil
}
