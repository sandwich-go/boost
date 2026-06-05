package xrand

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// random_select.go 的功能测试。
//
// 由于这些函数依赖 math/rand 是非确定性的，测试策略：
//   1. 边界条件断言（empty / total<=0 走 fallback 路径）
//   2. 概率分布断言：用大样本 + 误差容忍验证抽样命中率
//   3. 范围断言：返回值落在期望区间

func TestRandomInt(t *testing.T) {
	Convey("RandomInt 落在 [min, max]", t, func() {
		min, max := 5, 10
		// 抽 1000 次，每次结果都应在 [min, max]
		for i := 0; i < 1000; i++ {
			v := RandomInt(min, max)
			So(v, ShouldBeBetweenOrEqual, min, max)
		}
	})

	Convey("RandomInt min == max 时确定返回 min", t, func() {
		So(RandomInt(7, 7), ShouldEqual, 7)
	})
}

func TestRandomSelectOneFromMap(t *testing.T) {
	Convey("空 map / total<=0 返回 (0, false)", t, func() {
		got, ok := RandomSelectOneFromMap(nil)
		So(ok, ShouldBeFalse)
		So(got, ShouldEqual, int32(0))

		got2, ok2 := RandomSelectOneFromMap(map[int32]int32{1: 0, 2: 0})
		So(ok2, ShouldBeFalse)
		So(got2, ShouldEqual, int32(0))
	})

	Convey("权重抽样：rate=100% 的 key 必中", t, func() {
		// 只有一个 key 有权重，必命中它
		m := map[int32]int32{42: 100, 99: 0}
		for i := 0; i < 100; i++ {
			got, ok := RandomSelectOneFromMap(m)
			So(ok, ShouldBeTrue)
			So(got, ShouldEqual, int32(42))
		}
	})

	Convey("权重抽样：分布大致符合权重比", t, func() {
		// 1 -> rate 10%, 2 -> rate 90%；抽 10000 次，2 应该出 ~9000
		m := map[int32]int32{1: 10, 2: 90}
		var c1, c2 int
		for i := 0; i < 10000; i++ {
			got, ok := RandomSelectOneFromMap(m)
			So(ok, ShouldBeTrue)
			switch got {
			case 1:
				c1++
			case 2:
				c2++
			}
		}
		So(c1+c2, ShouldEqual, 10000)
		// 期望比 90/10；±5% 容差（10000 样本下足够稳）
		So(c2, ShouldBeBetween, 8500, 9500)
	})
}

func TestRandomSelectOneKeyFromMap_Generic(t *testing.T) {
	Convey("string key + uint32 value", t, func() {
		m := map[string]uint32{"a": 0, "b": 100}
		// 只有 "b" 有权重，必命中 "b"
		for i := 0; i < 50; i++ {
			got, ok := RandomSelectOneKeyFromMap[string, uint32](m)
			So(ok, ShouldBeTrue)
			So(got, ShouldEqual, "b")
		}
	})

	Convey("int64 key + int64 value 空 map", t, func() {
		got, ok := RandomSelectOneKeyFromMap[int64, int64](map[int64]int64{})
		So(ok, ShouldBeFalse)
		So(got, ShouldEqual, int64(0))
	})
}

func TestRandomSelectOneFromArray(t *testing.T) {
	Convey("空 array / total<=0 返回 (0, false)", t, func() {
		got, ok := RandomSelectOneFromArray(nil)
		So(ok, ShouldBeFalse)
		So(got, ShouldEqual, 0)

		got2, ok2 := RandomSelectOneFromArray([]int32{0, 0, 0})
		So(ok2, ShouldBeFalse)
		So(got2, ShouldEqual, 0)
	})

	Convey("单元素必中 index 0", t, func() {
		got, ok := RandomSelectOneFromArray([]int32{50})
		So(ok, ShouldBeTrue)
		So(got, ShouldEqual, 0)
	})

	Convey("分布符合权重比 [20,30,50]", t, func() {
		arr := []int32{20, 30, 50}
		counts := make([]int, 3)
		for i := 0; i < 10000; i++ {
			idx, ok := RandomSelectOneFromArray(arr)
			So(ok, ShouldBeTrue)
			counts[idx]++
		}
		// 期望 2000 / 3000 / 5000；±10% 容差
		So(counts[0], ShouldBeBetween, 1700, 2300)
		So(counts[1], ShouldBeBetween, 2700, 3300)
		So(counts[2], ShouldBeBetween, 4700, 5300)
	})
}

func TestRandomSelectIndexFromArray_Generic(t *testing.T) {
	Convey("uint64 array 单元素", t, func() {
		got, ok := RandomSelectIndexFromArray([]uint64{42})
		So(ok, ShouldBeTrue)
		So(got, ShouldEqual, 0)
	})

	Convey("空 array 返回 (0, false)", t, func() {
		got, ok := RandomSelectIndexFromArray([]int8{})
		So(ok, ShouldBeFalse)
		So(got, ShouldEqual, 0)
	})
}

func TestIsSelected100n(t *testing.T) {
	Convey("number=0 永不命中", t, func() {
		for i := 0; i < 1000; i++ {
			So(IsSelected100n(0), ShouldBeFalse)
		}
	})

	Convey("number=100 必命中（覆盖整个 [1,100] 区间）", t, func() {
		for i := 0; i < 1000; i++ {
			So(IsSelected100n(100), ShouldBeTrue)
		}
	})

	Convey("number=50 命中率约 50%", t, func() {
		var hits int
		for i := 0; i < 10000; i++ {
			if IsSelected100n(50) {
				hits++
			}
		}
		// 期望 5000；±10% 容差
		So(hits, ShouldBeBetween, 4500, 5500)
	})

	Convey("number 超出 100 也总命中（rand.Int32N(100)+1 ≤ number）", t, func() {
		for i := 0; i < 100; i++ {
			So(IsSelected100n(200), ShouldBeTrue)
		}
	})
}

// BenchmarkRandomSelectOneKeyFromMap_Map10 / _Map100：测真实 hot path —
// 业务的"按权重抽奖"在玩家匹配 / 掉落系统等场景每秒可能万次调用。
// 这是泛型 + sort.Slice + map iteration 的复合开销，n 大时 sort 占大头。
func BenchmarkRandomSelectOneKeyFromMap_Map10(b *testing.B) {
	m := make(map[int32]int32, 10)
	for i := int32(0); i < 10; i++ {
		m[i] = 1 + i
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = RandomSelectOneKeyFromMap[int32, int32](m)
	}
}

func BenchmarkRandomSelectOneKeyFromMap_Map100(b *testing.B) {
	m := make(map[int32]int32, 100)
	for i := int32(0); i < 100; i++ {
		m[i] = 1 + i
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = RandomSelectOneKeyFromMap[int32, int32](m)
	}
}

// BenchmarkRandomSelectOneFromArray_Len100：array 版本不走 sort，作为
// map 版本的对照（同 size 下应该显著更快）。
func BenchmarkRandomSelectOneFromArray_Len100(b *testing.B) {
	arr := make([]int32, 100)
	for i := range arr {
		arr[i] = int32(i + 1)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = RandomSelectOneFromArray(arr)
	}
}

// BenchmarkIsSelected100n / BenchmarkRandomInt：极简函数，bench 主要看
// rand.IntN / Int32N 的调用开销（与 FastRandInt_Concurrent 的非并发对照）。
func BenchmarkIsSelected100n(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsSelected100n(50)
	}
}

func BenchmarkRandomInt_Single(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandomInt(1, 100)
	}
}
