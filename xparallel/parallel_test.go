package xparallel

import (
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：xparallel 之前 0% 覆盖。
// 测试 4 个并发 helper（MapK/MapV/SliceK/SliceV）+ closeThenParallel +
// worker 的 panic recover 路径。

func TestSliceV_NoLimit(t *testing.T) {
	Convey("SliceV 不限并发：所有元素都被处理一次", t, func() {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		var seen sync.Map
		var counter atomic.Int32

		SliceV(NoLimit, input, func(v int) {
			seen.Store(v, true)
			counter.Add(1)
		})

		So(counter.Load(), ShouldEqual, int32(len(input)))
		for _, v := range input {
			_, ok := seen.Load(v)
			So(ok, ShouldBeTrue)
		}
	})
}

func TestSliceV_WithLimit(t *testing.T) {
	Convey("SliceV maxp=2 限制并发数", t, func() {
		input := make([]int, 100)
		for i := range input {
			input[i] = i
		}

		var concurrent atomic.Int32
		var maxConcurrent atomic.Int32
		var wg sync.WaitGroup
		_ = wg

		SliceV(2, input, func(v int) {
			n := concurrent.Add(1)
			defer concurrent.Add(-1)
			// 更新 maxConcurrent
			for {
				old := maxConcurrent.Load()
				if n <= old || maxConcurrent.CompareAndSwap(old, n) {
					break
				}
			}
		})

		// maxp=2 时 maxConcurrent 应 ≤ 2
		So(maxConcurrent.Load(), ShouldBeLessThanOrEqualTo, int32(2))
	})
}

func TestSliceK(t *testing.T) {
	Convey("SliceK 传入 index 而非 value", t, func() {
		input := []string{"a", "b", "c", "d", "e"}
		var indices sync.Map

		SliceK(NoLimit, input, func(idx int) {
			indices.Store(idx, true)
		})

		// 应该看到 0..4 全部 index
		for i := 0; i < len(input); i++ {
			_, ok := indices.Load(i)
			So(ok, ShouldBeTrue)
		}
	})
}

func TestMapK_AndMapV(t *testing.T) {
	Convey("MapK 遍历 key", t, func() {
		input := map[string]int{"a": 1, "b": 2, "c": 3}
		var keys sync.Map
		MapK(NoLimit, input, func(k string) {
			keys.Store(k, true)
		})
		for k := range input {
			_, ok := keys.Load(k)
			So(ok, ShouldBeTrue)
		}
	})

	Convey("MapV 遍历 value", t, func() {
		input := map[string]int{"a": 1, "b": 2, "c": 3}
		var values sync.Map
		MapV(NoLimit, input, func(v int) {
			values.Store(v, true)
		})
		for _, v := range input {
			_, ok := values.Load(v)
			So(ok, ShouldBeTrue)
		}
	})
}

func TestEmptyInput(t *testing.T) {
	Convey("空输入不死循环不 panic", t, func() {
		So(func() { SliceV(NoLimit, []int{}, func(v int) {}) }, ShouldNotPanic)
		So(func() { SliceK(NoLimit, []string(nil), func(int) {}) }, ShouldNotPanic)
		So(func() { MapK(NoLimit, map[int]int{}, func(int) {}) }, ShouldNotPanic)
		So(func() { MapV(NoLimit, map[int]int(nil), func(int) {}) }, ShouldNotPanic)
	})
}

func TestWorker_PanicRecover(t *testing.T) {
	Convey("worker 内 fn panic 被 xpanic.Try 捕获，不影响其它元素处理", t, func() {
		input := []int{1, 2, 3, 4, 5}
		var processed atomic.Int32

		SliceV(NoLimit, input, func(v int) {
			if v == 3 {
				panic("oops on 3")
			}
			processed.Add(1)
		})

		// 4 个非 3 的元素应该都被处理（panic 被 recover 不阻断其它）
		So(processed.Load(), ShouldEqual, int32(4))
	})
}

// BenchmarkSliceV hot path：业务高频用 xparallel 并发处理 batch。
func BenchmarkSliceV_NoLimit(b *testing.B) {
	input := make([]int, 100)
	for i := range input {
		input[i] = i
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SliceV(NoLimit, input, func(v int) {
			_ = v
		})
	}
}
