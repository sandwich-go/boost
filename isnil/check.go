package isnil

import "reflect"

type isNiler interface {
	IsNil() bool
}

// Check checks if any is nil.
func Check(any interface{}) bool {
	if any == nil {
		return true
	}

	if checker, ok := any.(isNiler); ok {
		return checker.IsNil()
	}

	v := reflect.ValueOf(any)
	return v.Kind() == reflect.Ptr && v.IsNil()
}
