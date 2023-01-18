package annotation

import "errors"

var ErrNoAnnotation = errors.New("no annotation")

// Annotation 注释
type Annotation interface {
	// Name 注释名
	// 如，'// annotation@X( a = "A" )' 的注释名为 'X'
	Name() string

	// Line 注释文本内容
	Line() string

	// Contains 是否包含指定 key
	// 如，'// annotation@X( a = "A" )' 的包含 'a'
	Contains(key string) bool

	// String 指定 key 对应的字符串类型 value
	// 如，'// annotation@X( a = "A" )' 中 key 为 'a' 的 value 为 'A'
	String(key string, defaultVal ...string) string

	// Int8 指定 key 对应的 int8 类型 value
	Int8(key string, defaultVal ...int8) (int8, error)

	// Int16 指定 key 对应的 int16 类型 value
	Int16(key string, defaultVal ...int16) (int16, error)

	// Int32 指定 key 对应的 int32 类型 value
	Int32(key string, defaultVal ...int32) (int32, error)

	// Int64 指定 key 对应的 int64 类型 value
	Int64(key string, defaultVal ...int64) (int64, error)

	// Uint8 指定 key 对应的 uint8 类型 value
	Uint8(key string, defaultVal ...uint8) (uint8, error)

	// Uint16 指定 key 对应的 uint16 类型 value
	Uint16(key string, defaultVal ...uint16) (uint16, error)

	// Uint32 指定 key 对应的 uint32 类型 value
	Uint32(key string, defaultVal ...uint32) (uint32, error)

	// Uint64 指定 key 对应的 uint64 类型 value
	Uint64(key string, defaultVal ...uint64) (uint64, error)

	// Int 指定 key 对应的 int 类型 value
	Int(key string, defaultVal ...int) (int, error)

	// Float32 指定 key 对应的 float32 类型 value
	Float32(key string, defaultVal ...float32) (float32, error)

	// Float64 指定 key 对应的 float64 类型 value
	Float64(key string, defaultVal ...float64) (float64, error)

	// Bool 指定 key 对应的 bool 类型 value
	Bool(key string, defaultVal ...bool) (bool, error)
}

// Descriptor 描述，可以规定萃取的注释名，以及判断萃取的注释是否合法
type Descriptor struct {
	Name      string
	Validator func(Annotation) bool
}

type Resolver interface {
	// Resolve 解析一行注释
	// 若未解析成功，则返回 ErrNoAnnotation 错误
	Resolve(line string) (Annotation, error)

	// ResolveWithName 解析一行注释，但要求 Annotation.Name 是指定的 name 参数
	// 否则返回 ErrNoAnnotation 错误
	ResolveWithName(lines []string, name string) (Annotation, error)

	// ResolveMultiple 解析多行注释
	ResolveMultiple(lines []string) ([]Annotation, error)

	// ResolveNoDuplicate 解析多行注释，但不允许有重复的 Annotation.Name
	// 否则返回错误
	ResolveNoDuplicate(lines []string) ([]Annotation, error)
}

// Default 默认的解析器，只有包含 'annotation@' 的行，才能萃取到注释
var Default = New()

// Resolve 使用默认的解析器 Default 解析一行注释
func Resolve(line string) (Annotation, error) { return Default.Resolve(line) }

// ResolveMultiple 使用默认的解析器 Default 解析多行注释
func ResolveMultiple(lines []string) ([]Annotation, error) { return Default.ResolveMultiple(lines) }

// ResolveWithName 使用默认的解析器 Default 解析一行注释，但要求 Annotation.Name 是指定的 name 参数
// 否则返回 ErrNoAnnotation 错误
func ResolveWithName(lines []string, name string) (Annotation, error) {
	return Default.ResolveWithName(lines, name)
}

// ResolveNoDuplicate 使用默认的解析器 Default 解析多行注释，但不允许有重复的 Annotation.Name
// 否则返回错误
func ResolveNoDuplicate(lines []string) ([]Annotation, error) {
	return Default.ResolveNoDuplicate(lines)
}
