package xmath

import (
	"testing"
)

func TestDisturb(t *testing.T) {
	n := 10000
	for i := 0; i < 100; i++ {
		n1 := Disturb(n, 10)
		t.Log(n1)
		if n1 < 9000 || n1 > 11000 {
			t.Fail()
		}
	}
}
