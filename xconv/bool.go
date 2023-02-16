package xconv

import (
	"github.com/sandwich-go/boost/xstrings"
	"reflect"
	"strings"
)

// Bool [影响性能] converts `any` to bool.
// It returns false if `any` is: false, "", 0, "0", "false", "off", "n", "no", empty slice/map.
func Bool(any interface{}) bool {
	if any == nil {
		return false
	}
	switch value := any.(type) {
	case bool:
		return value
	case []byte:
		return xstrings.IsTrue(strings.ToLower(string(value)))
	case string:
		return xstrings.IsTrue(strings.ToLower(value))
	case int:
		return value != 0
	case int8:
		return value != 0
	case int16:
		return value != 0
	case int32:
		return value != 0
	case int64:
		return value != 0
	case uint:
		return value != 0
	case uint8:
		return value != 0
	case uint16:
		return value != 0
	case uint32:
		return value != 0
	case uint64:
		return value != 0
	default:
		if f, ok := value.(iBool); ok {
			return f.Bool()
		}
		rv := reflect.ValueOf(any)
		switch rv.Kind() {
		case reflect.Ptr:
			return !rv.IsNil()
		case reflect.Map:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			return rv.Len() != 0
		case reflect.Struct:
			return true
		default:
			return xstrings.IsTrue(strings.ToLower(String(any)))
		}
	}
}
