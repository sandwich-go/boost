package isnil

import (
	"github.com/modern-go/reflect2"
)

// Check checks if any is nil.
func Check(any interface{}) bool {
	return reflect2.IsNil(any)
}
