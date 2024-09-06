package xerror

import (
	"sync"
)

// SafeArray 错误数组，可以将多个 error 进行组装，并当成 error 进行函数传递或返回
// 协程安全
type SafeArray struct {
	mx  sync.RWMutex
	arr Array
}

// Error 实现 error 接口
func (e *SafeArray) Error() string {
	e.mx.RLock()
	defer e.mx.RUnlock()

	return e.arr.Error()
}

// Push 推入一个错误信息，err如果为nil则丢弃
func (e *SafeArray) Push(err error) {
	if err == nil {
		return
	}

	e.mx.Lock()
	defer e.mx.Unlock()

	e.arr.Push(err)
}

// LastErr 返回最后一个错误信息，如果没有错误则返回nil
func (e *SafeArray) LastErr() error {
	e.mx.RLock()
	defer e.mx.RUnlock()

	return e.arr.LastErr()
}

// Err 返回标准error对象，如果错误列表为空则返回nil
func (e *SafeArray) Err() error {
	e.mx.RLock()
	defer e.mx.RUnlock()

	return e.arr.Err()
}

// Is 对 errors.Is 的支持
func (e *SafeArray) Is(target error) bool {
	e.mx.RLock()
	defer e.mx.RUnlock()

	return e.arr.Is(target)
}

// String
func (e *SafeArray) String() string {
	e.mx.RLock()
	defer e.mx.RUnlock()

	return e.arr.String()
}

// WrappedErrors 返回内部所有的 error
func (e *SafeArray) WrappedErrors() []error {
	e.mx.RLock()
	defer e.mx.RUnlock()

	var out = make([]error, 0, len(e.arr.errors))
	for _, v := range e.arr.errors {
		out = append(out, v)
	}

	return out
}

// SetFormatFunc 设置格式化 error 数组函数，默认 ListFormatFunc
func (e *SafeArray) SetFormatFunc(f ErrorFormatFunc) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.arr.SetFormatFunc(f)
}
