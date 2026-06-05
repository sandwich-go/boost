package redblacktree

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 redblacktree 包未覆盖的：
//   - Successor / Predecessor（节点前驱/后继）
//   - Clear / Release / releaseSubtree
//   - WalkTailNodeKeys / WalkTailNodes（按 Ceiling 后向走）
//   - WalkHeadNodeKeys / WalkHeadNodes（按 Floor 后逆走，已覆盖部分）

func TestNode_SuccessorAndPredecessor(t *testing.T) {
	Convey("Successor 返回中序遍历下一个节点", t, func() {
		tree := New[int, string]()
		for _, k := range []int{5, 3, 7, 1, 4, 6, 8} {
			tree.Put(k, "")
		}
		// 中序：1,3,4,5,6,7,8
		// 1 的 Successor = 3
		n1 := tree.GetNode(1)
		So(n1, ShouldNotBeNil)
		s := n1.Successor()
		So(s, ShouldNotBeNil)
		So(s.Key, ShouldEqual, 3)

		// 5 的 Successor = 6
		n5 := tree.GetNode(5)
		So(n5.Successor().Key, ShouldEqual, 6)

		// 8 是最大值，Successor = nil
		n8 := tree.GetNode(8)
		So(n8.Successor(), ShouldBeNil)
	})

	Convey("Predecessor 返回中序遍历前一个节点", t, func() {
		tree := New[int, string]()
		for _, k := range []int{5, 3, 7, 1, 4, 6, 8} {
			tree.Put(k, "")
		}
		// 8 的 Predecessor = 7
		n8 := tree.GetNode(8)
		p := n8.Predecessor()
		So(p, ShouldNotBeNil)
		So(p.Key, ShouldEqual, 7)

		// 4 的 Predecessor = 3
		n4 := tree.GetNode(4)
		So(n4.Predecessor().Key, ShouldEqual, 3)

		// 1 是最小值，Predecessor = nil
		n1 := tree.GetNode(1)
		So(n1.Predecessor(), ShouldBeNil)
	})
}

func TestTree_ClearAndRelease(t *testing.T) {
	Convey("Clear 清空 tree（保留 Comparator）", t, func() {
		tree := New[int, string]()
		for i := 1; i <= 10; i++ {
			tree.Put(i, "")
		}
		So(tree.Size(), ShouldEqual, 10)

		tree.Clear()
		So(tree.Size(), ShouldEqual, 0)
		So(tree.Root, ShouldBeNil)
		// Comparator 应仍可用：再插入不 panic
		tree.Put(1, "after-clear")
		v, found := tree.Get(1)
		So(found, ShouldBeTrue)
		So(v, ShouldEqual, "after-clear")
	})

	Convey("Release 归还 tree 到 pool（清空 + Comparator nil）", t, func() {
		tree := New[int, string]()
		for i := 1; i <= 5; i++ {
			tree.Put(i, "")
		}
		tree.Release()
		// Release 之后 tree 已归还到 pool，不应继续使用
		// 但 *Tree 字段已被清空：Root nil, size 0, Comparator nil
		So(tree.Root, ShouldBeNil)
		So(tree.Size(), ShouldEqual, 0)
		// Comparator 已置 nil，再 Put 会 panic
		// 这是 Release 后对象已死的 doc 化契约
	})

	Convey("Release nil tree 安全", t, func() {
		var tree *Tree[int, string]
		So(func() { tree.Release() }, ShouldNotPanic)
	})
}

func TestTree_WalkTailNodes(t *testing.T) {
	tree := New[int, string]()
	for _, k := range []int{1, 3, 5, 7, 9} {
		tree.Put(k, "")
	}

	Convey("WalkTailNodeKeys 从 Ceiling 走到尾", t, func() {
		var keys []int
		tree.WalkTailNodeKeys(4, func(k int) bool {
			keys = append(keys, k)
			return false // 不中止
		})
		// Ceiling(4) = 5；后续 5,7,9
		So(keys, ShouldResemble, []int{5, 7, 9})
	})

	Convey("WalkTailNodes 提前中止", t, func() {
		var keys []int
		tree.WalkTailNodes(0, func(k int, v string) bool {
			keys = append(keys, k)
			return len(keys) >= 2 // 走两个就停
		})
		// Ceiling(0) = 1；前两个 1,3
		So(keys, ShouldResemble, []int{1, 3})
	})

	Convey("WalkTailNodeKeys 在不存在 Ceiling（k 大于最大值）时不调 handler", t, func() {
		var keys []int
		tree.WalkTailNodeKeys(100, func(k int) bool {
			keys = append(keys, k)
			return false
		})
		So(keys, ShouldBeEmpty)
	})
}

func TestTree_WalkHeadNodes(t *testing.T) {
	tree := New[int, string]()
	for _, k := range []int{1, 3, 5, 7, 9} {
		tree.Put(k, "")
	}

	Convey("WalkHeadNodeKeys 从 Floor 反向走到头", t, func() {
		var keys []int
		tree.WalkHeadNodeKeys(6, func(k int) bool {
			keys = append(keys, k)
			return false
		})
		// Floor(6) = 5；逆走 5,3,1
		So(keys, ShouldResemble, []int{5, 3, 1})
	})

	Convey("WalkHeadNodes 提前中止", t, func() {
		var keys []int
		tree.WalkHeadNodes(100, func(k int, v string) bool {
			keys = append(keys, k)
			return len(keys) >= 2
		})
		// Floor(100) = 9；逆走 9,7
		So(keys, ShouldResemble, []int{9, 7})
	})
}

// TestTree_FullSortRoundtrip 通过插入随机顺序 + 中序遍历得到升序，
// 强化对 Successor 链的正确性验证。
func TestTree_FullSortRoundtrip(t *testing.T) {
	Convey("插入大量随机 key，中序遍历得到升序", t, func() {
		tree := New[int, string]()
		input := []int{50, 30, 70, 20, 40, 60, 80, 10, 35, 65, 90, 5, 25, 55, 75, 95}
		for _, k := range input {
			tree.Put(k, "")
		}

		// 用 Successor 链从最小值（leftmost）走到尾
		node := tree.Left()
		var got []int
		for node != nil {
			got = append(got, node.Key)
			node = node.Successor()
		}

		expected := make([]int, len(input))
		copy(expected, input)
		sort.Ints(expected)
		So(got, ShouldResemble, expected)
	})
}
