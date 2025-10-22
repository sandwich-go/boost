package misc

import "unsafe"

// Extend is a helper struct to extend a struct with a pointer to itself.
// 该结构体必须放在继承的结构体的第一个字段
type Extend[S any, PS *S] struct{}

func (e *Extend[S, PS]) Self() PS {
	return (*S)(unsafe.Pointer(e))
}
