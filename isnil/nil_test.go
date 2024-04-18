package isnil

import (
	"fmt"
	"testing"
)

func TestIsNil(t *testing.T) {
	tests := []struct {
		any interface{}
		e   bool
	}{
		{nil, true},

		{fmt.Stringer(nil), true},
	}

	for i, tt := range tests {
		func() {
			if res := Check(tt.any); res != tt.e {
				t.Errorf("#%d: got %t xpected %t", i, res, tt.e)
			}
		}()
	}
}
