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

	return reflect.ValueOf(any).Kind() == reflect.Ptr && reflect.ValueOf(any).IsNil()
}
