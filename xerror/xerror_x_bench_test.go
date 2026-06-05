package xerror_test

import (
	"errors"
	"testing"

	"github.com/sandwich-go/boost/xerror"
)

// xerror 的 hot path：错误产生（NewText / NewCode / Wrap / WrapCode），
// 下游业务每一次错误流转都会走。重点关注两个开销：
//   1. fmt.Sprintf 格式化：format 参数变动会让分配数翻倍
//   2. callersCheckIsErrorWithStack：runtime.Callers 抓栈是 §11.4 提到的
//      "错误产生时的真大头，比 fmt.Sprintf 还重"
//
// bench 全部入仓作为 baseline。修复历史问题：原 printIfError 在 b.Error
// 里调用 b.Error 让每次 NewText 都判 fail（NewText 永远返回非 nil 的
// *Error）—— 改为全局 sink 变量防止编译器把构造调用优化掉。

var (
	errBase = errors.New("test")

	// sinkErr 全局变量，逃逸到堆，让编译器认为 err 被使用过；防止 NewText /
	// Wrap 等纯构造调用被 dead-code-elimination。
	sinkErr error
	sinkStr string
)

func Benchmark_NewText(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkErr = xerror.NewText("test")
	}
}

func Benchmark_NewText_Format(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkErr = xerror.NewText("%s", "test")
	}
}

func Benchmark_Wrap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkErr = xerror.Wrap(errBase, "test")
	}
}

func Benchmark_Wrap_Format(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkErr = xerror.Wrap(errBase, "%s", "test")
	}
}

func Benchmark_NewCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkErr = xerror.NewCode(500, "test")
	}
}

func Benchmark_NewCode_Format(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkErr = xerror.NewCode(500, "%s", "test")
	}
}

func Benchmark_WrapCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkErr = xerror.WrapCode(500, errBase, "test")
	}
}

func Benchmark_WrapCode_Format(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkErr = xerror.WrapCode(500, errBase, "%s", "test")
	}
}

// Benchmark_Stack 测 *Error.Stack() 格式化输出栈帧的开销。生产里 log 一条
// error 时会调（zerolog/zap 的 .Stack() 等），下游高频日志场景可能 hot。
//
// 注意：构造 *Error 时已经抓了栈（callers），Stack() 只是把抓好的 stack 数组
// 格式化为 string。所以 Stack() 比 NewText 的整体开销小一个数量级，但 N 个
// stack 帧 × M 条 error 还是会显眼。
func Benchmark_Stack_NoStackFlag(b *testing.B) {
	// IsErrorWithStack=false（默认）：抓的栈被丢弃，Stack() 走 fallback
	err := xerror.NewText("test error")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkStr = xerror.Stack(err)
	}
}

func Benchmark_Stack_WithStackFlag(b *testing.B) {
	// IsErrorWithStack=true：完整栈被保留，Stack() 真做帧迭代+格式化
	old := xerror.IsErrorWithStack
	xerror.IsErrorWithStack = true
	defer func() { xerror.IsErrorWithStack = old }()
	err := xerror.NewText("test error")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkStr = xerror.Stack(err)
	}
}
