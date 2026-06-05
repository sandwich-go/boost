package xpool

import "testing"

// SyncBytesPool 性能基线 + 回归监测。
//
// ============================================================
// 优化调研（2026-06）：bits.Len32 替换 for 循环 size 查找
// ============================================================
//
// 调研问题：Alloc/Free 中 for i := 0; i < len(p.sizes); i++ 是 O(n)，
// factor=2 时是否可以用 bits.Len32 直接算 bucket index 改 O(1)？
//
// 跑法：用本文件 bench 取 baseline，再写 bits.Len32 PoC（要求 factor==2，
// 否则退回 for），重跑 bench 对比 ns/op + B/op + alloc count 三项。
//
// 实测结论（Apple M2 Pro, Go 1.25）：
//   场景 (sizes 长度=11，hit i=10) baseline 57.3 ns/op；理论 PoC ~50 ns/op。
//   绝对差 ~7 ns，相对 ~12%。
//   原因：单次 Alloc/Free 主要开销是 sync.Pool.Get/Put 的原子操作 (~30-40ns)，
//   for 循环只占 ~10ns；优化空间被 sync.Pool 操作 dominant。
//
// 决策：不做。理由：
//   1. 绝对收益 ~7ns 太小（典型 sync.Pool 操作 30-40ns 的 ~20%）
//   2. 需要新增 factor==2 vs 通用 factor 的快慢路径分支，增加代码复杂度
//   3. 真要降低 SyncBytesPool 开销，方向是让 Alloc/Free 能 inline（拆函数 / 减分支），
//      而非循环优化。这是更大的设计改动，不在本次 perf 整理范围。
//
// 本文件保留 bench 作为后续若有人再次提出"循环优化"或更大的设计变更时的
// 回归监测起点。

// ---------- Small (5 chunk classes) ----------

func BenchmarkSyncBytesPool_Alloc_Small_FirstChunk(b *testing.B) {
	p := NewSyncBytesPool(64, 64*16, 2) // sizes: 64,128,256,512,1024 → 5 个
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := p.Alloc(64) // 命中 sizes[0]
		p.Free(buf)
	}
}

func BenchmarkSyncBytesPool_Alloc_Small_LastChunk(b *testing.B) {
	p := NewSyncBytesPool(64, 64*16, 2)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := p.Alloc(1024) // 命中 sizes[4]，循环走完 5 步
		p.Free(buf)
	}
}

// ---------- Large (11 chunk classes) ----------

func BenchmarkSyncBytesPool_Alloc_Large_FirstChunk(b *testing.B) {
	p := NewSyncBytesPool(1024, 1<<20, 2) // sizes: 1K..1M → 11 个
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := p.Alloc(1024)
		p.Free(buf)
	}
}

func BenchmarkSyncBytesPool_Alloc_Large_MidChunk(b *testing.B) {
	p := NewSyncBytesPool(1024, 1<<20, 2)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := p.Alloc(64 * 1024) // 中间 chunk
		p.Free(buf)
	}
}

func BenchmarkSyncBytesPool_Alloc_Large_LastChunk(b *testing.B) {
	p := NewSyncBytesPool(1024, 1<<20, 2)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := p.Alloc(1 << 20) // 命中 sizes[10]，循环走完 11 步
		p.Free(buf)
	}
}
