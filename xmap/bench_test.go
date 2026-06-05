package xmap

import (
	"strconv"
	"testing"
)

// xmap 的 hot path 评估（按 §3.4）：
//   - WalkMapDeterministic：业务里"按确定性顺序遍历 map"用得最多
//     （日志输出 / 序列化 / 哈希计算等），每次走 sort.Slice O(n log n) +
//     遍历 O(n)，map 大时性能曲线显著
//   - Equal：cache 一致性比较 / 配置 diff 等场景，关注早返路径（长度不一
//     致就立刻 false）vs 全扫
//
// 不加 bench 的函数（按 §3.4 不做占位 bench）：
//   - Diff / ToMap：通常一次性比较 / 类型转换，无明确优化诉求；当前没有
//     性能投诉，加 bench 仅作 baseline 占位无价值
//
// bench 入仓作为后续若有人提"WalkMapDeterministic 排序优化"等决策的起点。

var (
	sinkBool bool
	sinkInt  int
)

// makeIntStringMap 构造稳定 size 的测试 map（每次 bench 入参一致，避免
// map 实现的随机迭代顺序污染基线）。
func makeIntStringMap(size int) map[int]string {
	m := make(map[int]string, size)
	for i := 0; i < size; i++ {
		m[i] = strconv.Itoa(i)
	}
	return m
}

// BenchmarkWalkMapDeterministic_N10 / _N100 / _N1000：n 不同时反映 sort
// 占比变化（n=10 时 sort 是噪声，n=1000 时是大头）。
func BenchmarkWalkMapDeterministic_N10(b *testing.B) {
	m := makeIntStringMap(10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		WalkMapDeterministic(m, func(k int, v string) bool {
			sinkInt = k
			return true
		})
	}
}

func BenchmarkWalkMapDeterministic_N100(b *testing.B) {
	m := makeIntStringMap(100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		WalkMapDeterministic(m, func(k int, v string) bool {
			sinkInt = k
			return true
		})
	}
}

func BenchmarkWalkMapDeterministic_N1000(b *testing.B) {
	m := makeIntStringMap(1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		WalkMapDeterministic(m, func(k int, v string) bool {
			sinkInt = k
			return true
		})
	}
}

// BenchmarkEqual_FastPath_LenMismatch：a / b 长度不同走早返路径。
// 期望 ns/op 极低、0 alloc —— 这是 Equal 的 happy fast path。
func BenchmarkEqual_FastPath_LenMismatch(b *testing.B) {
	a := makeIntStringMap(100)
	bb := makeIntStringMap(50) // 长度不同
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkBool = Equal(a, bb)
	}
}

// BenchmarkEqual_FullScan_AllEqual：a == b 必须扫完全部 entry。
// 这是 Equal 真正的最坏开销（O(n) value 比较）。
func BenchmarkEqual_FullScan_AllEqual(b *testing.B) {
	a := makeIntStringMap(100)
	bb := makeIntStringMap(100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkBool = Equal(a, bb)
	}
}

// BenchmarkEqual_FullScan_OneDiff：扫到一半发现差异 break。
// 介于 fast path 和 full scan 之间的"普通错路"。
func BenchmarkEqual_FullScan_OneDiff(b *testing.B) {
	a := makeIntStringMap(100)
	bb := makeIntStringMap(100)
	bb[50] = "different"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkBool = Equal(a, bb)
	}
}
