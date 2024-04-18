package isnil

import (
	"fmt"
	"testing"
)

type foo struct{}

func TestIsNil(t *testing.T) {
	cases := []struct {
		in  interface{}
		exp bool
	}{
		{1, false},
		{"nope", false},
		{foo{}, false},
		{&foo{}, false},
		{nil, true},
		{(*foo)(nil), true},
		{(interface{})(nil), true},
		{fmt.Stringer(nil), true},
	}

	for _, tcase := range cases {
		if out := Check(tcase.in); out != tcase.exp {
			if tcase.exp {
				t.Errorf("Expected %++v to be nil", tcase.in)
			} else {
				t.Errorf("Expected %++v to not be nil", tcase.in)
			}
		}
	}
}

func BenchmarkEqNilBasic(b *testing.B) {
	var v *int
	for i := 0; i < b.N; i++ {
		_ = v != nil
	}
}

func BenchmarkEqNilInterface(b *testing.B) {
	var v interface{}
	v = (*foo)(nil)
	for i := 0; i < b.N; i++ {
		_ = v != nil
	}
}

func BenchmarkIsNilBasic(b *testing.B) {
	var v *int
	for i := 0; i < b.N; i++ {
		Check(v)
	}
}

func BenchmarkIsNilInterface(b *testing.B) {
	var v interface{}
	v = (*foo)(nil)
	for i := 0; i < b.N; i++ {
		Check(v)
	}
}

func BenchmarkIsNilNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Check(nil)
	}
}
