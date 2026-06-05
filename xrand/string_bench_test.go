package xrand

import "testing"

// 性能基线：用于回归监测。可用 -bench=. -benchmem 对比 before/after。
//
// 当前基线（Apple M2 Pro, Go 1.25, n=16）：
//   String_NoLetterList                ~ 91 ns / 16 B / 1 alloc
//   String_WithLetterList              ~ 118 ns / 24 B / 2 alloc
//   StringWithTimestamp_NoLetterList   ~ 143 ns / 48 B / 1 alloc
//   StringWithTimestamp_WithLetterList ~ 164 ns / 56 B / 2 alloc

func BenchmarkString_NoLetterList(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = String(16)
	}
}

func BenchmarkString_WithLetterList(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = String(16, "ABC", "xyz")
	}
}

func BenchmarkStringWithTimestamp_NoLetterList(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = StringWithTimestamp(16)
	}
}

func BenchmarkStringWithTimestamp_WithLetterList(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = StringWithTimestamp(16, "ABC", "xyz")
	}
}
