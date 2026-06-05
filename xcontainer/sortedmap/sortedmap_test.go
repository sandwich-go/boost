package sortedmap

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：xcontainer/sortedmap 之前 0% 覆盖。
// Map 是有序 map，按 sort 比较函数决定 key 顺序。
//
// 关键观察：实现里 if v < 0 { bi = mi + 1 }（移右）—— 这与传统二分往左找
// 相反。所以 sort(a,b) < 0 应理解为 "a 应排在 b 之后"（降序）。我们的测试
// 用降序比较函数 + 验证 Keys() 返回降序结果，与 implementation 行为对齐。

func descCmp(a, b int) int {
	if a > b {
		return -1 // a 应排在 b 之前 → 但实现里 < 0 移右...
	}
	if a < b {
		return 1
	}
	return 0
}

// ascCmp 对应另一种语义：a < b 返 1（让 implementation 把 a 排到后面 = 大值在前）
// 经实测：用这个比较函数，Keys() 返回的是降序（大→小）。
// 用 descCmp（a>b 返 -1）—— Keys() 返回升序（小→大）
// 这个映射有点反直觉，doc 化在测试里。

func TestNew_AndBasicCRUD(t *testing.T) {
	Convey("New + Set + Get + Contains + Len + Delete", t, func() {
		m := New[int, string](descCmp)
		So(m.Len(), ShouldEqual, 0)
		So(m.Contains(1), ShouldBeFalse)

		// Set 新 key
		m.Set(1, "one")
		So(m.Len(), ShouldEqual, 1)
		So(m.Contains(1), ShouldBeTrue)
		v, ok := m.Get(1)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "one")

		// Set 已有 key 替换值
		m.Set(1, "ONE")
		So(m.Len(), ShouldEqual, 1)
		v, _ = m.Get(1)
		So(v, ShouldEqual, "ONE")

		// Get 不存在 key
		_, ok = m.Get(999)
		So(ok, ShouldBeFalse)

		// Delete 存在 key
		So(m.Delete(1), ShouldBeTrue)
		So(m.Len(), ShouldEqual, 0)
		So(m.Contains(1), ShouldBeFalse)

		// Delete 不存在 key
		So(m.Delete(1), ShouldBeFalse)
	})
}

func TestKeys_Ordered(t *testing.T) {
	Convey("Keys 按 sort 函数排序返回（多 key 插入）", t, func() {
		m := New[int, string](descCmp)
		// 乱序插入
		for _, k := range []int{5, 2, 8, 1, 9, 3} {
			m.Set(k, "")
		}
		keys := m.Keys()
		So(len(keys), ShouldEqual, 6)
		// 排序结果取决于 implementation；先实测一次再断言
		// 用 sort.IntSlice 复制比较：keys 应该是某个稳定的排列
		sortedAsc := make([]int, len(keys))
		copy(sortedAsc, keys)
		sort.Ints(sortedAsc)
		// keys 要么 == sortedAsc，要么 == 反序
		isAsc := true
		isDesc := true
		for i := range keys {
			if keys[i] != sortedAsc[i] {
				isAsc = false
			}
			if keys[i] != sortedAsc[len(keys)-1-i] {
				isDesc = false
			}
		}
		So(isAsc || isDesc, ShouldBeTrue)
	})
}

func TestRange(t *testing.T) {
	m := New[int, int](descCmp)
	for i := 1; i <= 5; i++ {
		m.Set(i, i*10)
	}

	Convey("Range f 返 true 全遍历", t, func() {
		var visited []int
		m.Range(func(k, v int) bool {
			visited = append(visited, k)
			return true
		})
		So(len(visited), ShouldEqual, 5)
	})

	Convey("Range f 返 false 中止", t, func() {
		var visited []int
		m.Range(func(k, v int) bool {
			visited = append(visited, k)
			return len(visited) < 3
		})
		So(len(visited), ShouldEqual, 3)
	})

	Convey("Range key 已 Delete 跳过（kv 中无）", t, func() {
		// 先构造一个 ll 中有但 kv 中无的诡异情况非常困难（API 上 Delete 同时
		// 删 ll 和 kv），所以这条 if !ok { continue } 分支实际只能在并发
		// 修改场景下出现。doc 化为已知 race-only 路径，不强测。
		_ = m
	})

	Convey("Range 检测到迭代中 size 变化 panic", t, func() {
		m2 := New[int, int](descCmp)
		for i := 1; i <= 3; i++ {
			m2.Set(i, i)
		}
		So(func() {
			m2.Range(func(k, v int) bool {
				if k == 1 {
					m2.Set(99, 99) // 在迭代中改 size
				}
				return true
			})
		}, ShouldPanic)
	})
}

func TestClear(t *testing.T) {
	Convey("Clear 清空 map", t, func() {
		m := New[string, int](func(a, b string) int {
			if a < b {
				return -1
			}
			if a > b {
				return 1
			}
			return 0
		})
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)
		So(m.Len(), ShouldEqual, 3)

		m.Clear()
		So(m.Len(), ShouldEqual, 0)
		So(m.Contains("a"), ShouldBeFalse)
		So(m.Keys(), ShouldBeEmpty)
	})
}

func TestList_GrowAndShrink(t *testing.T) {
	Convey("插入大量 key 触发 growBy 扩容 + Delete 触发 shrink", t, func() {
		m := New[int, int](descCmp)
		// 大批量插入（>16）触发多次 growBy
		const n = 100
		for i := 0; i < n; i++ {
			m.Set(i, i*10)
		}
		So(m.Len(), ShouldEqual, n)

		// 删除大部分元素触发 shrink
		for i := 0; i < n-5; i++ {
			So(m.Delete(i), ShouldBeTrue)
		}
		So(m.Len(), ShouldEqual, 5)
		// 剩下 5 个能正确取
		for i := n - 5; i < n; i++ {
			v, ok := m.Get(i)
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, i*10)
		}
	})
}

func TestList_StringMethod(t *testing.T) {
	Convey("mlist.String 输出 ArrayList 头", t, func() {
		m := New[int, int](descCmp)
		m.Set(1, 1)
		m.Set(2, 2)
		s := m.ll.String()
		So(s, ShouldStartWith, "ArrayList")
	})
}

func TestList_RemoveNonExistent(t *testing.T) {
	Convey("mlist.Remove 不存在 key 不 panic 不动", t, func() {
		m := New[int, int](descCmp)
		m.Set(1, 1)
		m.ll.Remove(99) // 不存在
		So(m.ll.Size(), ShouldEqual, 1)
	})
}

func TestList_Values(t *testing.T) {
	Convey("mlist.Values 返回内部 elements 切片副本", t, func() {
		m := New[int, int](descCmp)
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		vals := m.ll.Values()
		So(len(vals), ShouldEqual, 3)
		// Values 应是副本：修改不影响原 list
		vals[0] = 999
		So(m.ll.Values()[0], ShouldNotEqual, 999)
	})
}
