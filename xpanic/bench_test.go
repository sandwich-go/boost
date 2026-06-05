package xpanic

import (
	"errors"
	"testing"
)

// hot path bench：仓内 xtime/dispatcher / xparallel/parallel / xcompress/snappy
// 都在每次回调 / 每个 task / 每次 encode 里用 xpanic.Try.Catch 包裹业务代码。
// 这是真实 hot path：每次进入 Try 都要 defer + 闭包装箱 + recover()。
//
// bench 关注两个稳态指标：
// 1. 无 panic 路径（happy case）—— 业务里 99% 是这条；defer + recover() 调用
//    本身的开销
// 2. 有 panic 路径（error case）—— recover 拿到 value、构造 exception 结构体、
//    抓栈 (debug.Stack)。栈抓取比 defer 慢一个数量级，看 alloc 数能否给后续
//    优化决策提供 baseline
//
// 不加 bench 的函数（按 §3.4 不做占位 bench）：
//   - WhenError / WhenTrue / WhenFalse / WhenNil / WhenNotNil / WhenHereNotNil /
//     WhenErrorAsFmtFirst / Throw / AutoRecover / Catch / Do
//   均为错误路径或 panic 路径，调用频次低或仅在出错时承担开销，加 bench
//   只能量到分支判断本身（~ns 级），无优化价值。

// BenchmarkTry_NoPanic 业务热点：每次 Try 都付的 defer + closure 装箱开销。
func BenchmarkTry_NoPanic(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Try(func() {
			// 模拟业务工作：原子操作避免被编译器优化掉
			_ = i
		}).Catch(func(_ E) {
			b.Fatal("Catch should not run on no-panic path")
		})
	}
}

// BenchmarkTry_NoPanic_WithFinally Finally 路径：defer chain 多一层。
func BenchmarkTry_NoPanic_WithFinally(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Try(func() {
			_ = i
		}).Finally(func() {
			_ = i
		}).Catch(func(_ E) {
			b.Fatal("Catch should not run on no-panic path")
		})
	}
}

// BenchmarkTry_PanicString 错误路径 baseline：panic + recover + debug.Stack。
// 业务很少跑这条，但 dispatcher / parallel 在错误风暴时会触发，量出来作为
// "panic 风暴单次开销" 的 baseline。
func BenchmarkTry_PanicString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Try(func() {
			panic("benchmark panic")
		}).Catch(func(_ E) {
			// 不做事，纯量 Try/Catch 框架开销
		})
	}
}

// BenchmarkTry_PanicError 错误路径，panic value 是 error（更常见的业务路径
// 比如 xpanic.WhenErrorAsFmtFirst 内部 Sprintf 出 string 后 panic 的等价模拟，
// 但这里直接 panic error 让 catch handler 拿到 typed value）。
func BenchmarkTry_PanicError(b *testing.B) {
	err := errors.New("benchmark error")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Try(func() {
			panic(err)
		}).Catch(func(_ E) {
			// 不做事
		})
	}
}
