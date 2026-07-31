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
// 判定规则：接口本身为 nil，或接口的数据字为 nil。
//
// pointer-shaped 类型（指针、map、chan、func、unsafe.Pointer）是 direct interface，
// 值本身就存在数据字里，因此 (*T)(nil)、nil map、nil chan、nil func 一律返回 true。
//
// slice 是三字结构，不是 pointer-shaped，装箱必须间接存：runtime.convTslice 对 nil
// slice 取 &zeroVal[0]，数据字非 nil，因此 ([]T)(nil) 返回 false。要判空 slice 请用
// len(s) == 0，别指望这里。TestCheckTruthTable 钉住了完整真值表。
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
