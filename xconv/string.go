package xconv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/sandwich-go/boost/z"
)

// String [影响性能] converts `any` to string.
// String [影响性能] converts `any` to string.
func String(any interface{}) string {
	if any == nil {
		return ""
	}
	// 如果是 reflect.Value 类型
	if rv, ok := any.(reflect.Value); ok {
		// 如果 reflect.Value 是零值，返回空字符串
		if !rv.IsValid() {
			return ""
		}
		// 使用 Interface() 获取原始值并递归调用 String 方法
		return String(rv.Interface())
	}

	// 其他类型处理
	switch value := any.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	default:
		// Empty checks. govet 静态分析认为 line 16 已排除 nil，但 type switch
		// 的 default 可能匹配到 typed nil（如 (*MyType)(nil) 包成 interface{}
		// 后既不是 untyped nil 也不进任何 case），这里仍是合法防御。
		//nolint:govet // typed nil 检查，非 false alarm
		if value == nil {
			return ""
		}
		if f, ok := value.(iString); ok {
			return f.String()
		}
		if f, ok := value.(error); ok {
			return f.Error()
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return String(rv.Elem().Interface())
		}
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return z.BytesToString(jsonContent)
		}
	}
}
