package xerror_test

import (
	"errors"
	"testing"

	"github.com/sandwich-go/boost/xerror"
)

var (
	errBase = errors.New("test")
)

func printIfError(b *testing.B, err error) {
	if err != nil {
		b.Error(err)
	}
}

func Benchmark_NewText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printIfError(b, xerror.NewText("test"))
	}
}

func Benchmark_NewText_Format(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printIfError(b, xerror.NewText("%s", "test"))
	}
}

func Benchmark_Wrap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printIfError(b, xerror.Wrap(errBase, "test"))
	}
}

func Benchmark_Wrap_Format(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printIfError(b, xerror.Wrap(errBase, "%s", "test"))
	}
}
func Benchmark_NewCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printIfError(b, xerror.NewCode(500, "test"))
	}
}

func Benchmark_NewCode_Format(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printIfError(b, xerror.NewCode(500, "%s", "test"))
	}
}

func Benchmark_WrapCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printIfError(b, xerror.WrapCode(500, errBase, "test"))
	}
}

func Benchmark_WrapCode_Format(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printIfError(b, xerror.WrapCode(500, errBase, "test"))
	}
}
