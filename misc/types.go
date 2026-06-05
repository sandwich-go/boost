package misc

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/modern-go/reflect2"
)

// UnsafeShift0 通过指针移动 offset，来获取指定类型数据
func UnsafeShift0(value any, offset uintptr, sign int, t reflect.Type) any {
	p := pointerShift(reflect2.PtrOf(value), offset, sign)
	return reflect.NewAt(t.Elem(), p).Interface()
}

// UnsafeShift1 通过 *F 指针移动 offset，来获取指定类型数据
func UnsafeShift1[F any](value *F, offset uintptr, sign int, t reflect.Type) any {
	p := pointerShift(unsafe.Pointer(value), offset, sign)
	return reflect.NewAt(t.Elem(), p).Interface()
}

// UnsafeShift2 通过 *F 指针移动 offset，来获取指定类型数据 *T
func UnsafeShift2[F any, T any](value *F, offset uintptr, sign int) *T {
	p := pointerShift(unsafe.Pointer(value), offset, sign)
	return (*T)(p)
}

func pointerShift(p unsafe.Pointer, offset uintptr, sign int) unsafe.Pointer {
	if sign == 0 {
		return p
	}

	// 使用 unsafe.Add 以满足 unsafe.Pointer 的安全规则：指针算术必须在
	// 单个表达式内完成，避免产生裸 uintptr 中间值。否则在开启 -race
	// （同时启用 -d=checkptr）时，会触发
	// "checkptr: pointer arithmetic result points to invalid allocation"。
	if sign > 0 {
		return unsafe.Add(p, offset)
	}
	return unsafe.Add(p, -int64(offset))
}

// SliceCast slice 转换
func SliceCast[F any, T any](from []F) []T {
	if from == nil {
		return nil
	}

	ret := make([]T, len(from))
	for i, f := range from {
		var a any = f
		ret[i] = a.(T)
	}

	return ret
}

// FuncName 获取函数名
func FuncName(f any) string {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		panic("not a function")
	}
	return runtime.FuncForPC(v.Pointer()).Name()
}

// Zero 零值
func Zero[T any]() T {
	var v T
	return v
}

func ExportField(owner reflect.Value, field reflect.StructField) reflect.Value {
	f := owner.Field(field.Index[0])
	return reflect.NewAt(field.Type, unsafe.Pointer(f.UnsafeAddr())).Elem()
}
