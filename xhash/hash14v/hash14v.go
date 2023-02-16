package hash14v

import "github.com/sandwich-go/boost/z"

// V uint64 转换后的字节数组
type V []byte

func (v V) String() string { return z.BytesToString(v) }

// Id 字节数组转换后 uint64
type Id = uint64

// Converter Id 转换器，Id 和 V 可以相互转换
type Converter interface {
	// Offset 获取 Id 的 offset
	Offset() Id
	// ToId 对 V 进行转换，转换为 Id
	ToId(V) Id
	// ToV 对 Id 进行转换，转换为 V
	ToV(Id) V
}

// global is global hash14v converter, must goroutine safe, do not change inner properties while running
var global Converter

func init() {
	global = New()
}

// ToV 对 Id 进行转换，转换为 V
func ToV(id Id) V { return global.ToV(id) }

// ToId 对 V 进行转换，转换为 Id
func ToId(v V) Id { return global.ToId(v) }

// Offset 获取 Id 的 offset
func Offset() Id { return global.Offset() }
