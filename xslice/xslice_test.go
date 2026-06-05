package xslice

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContain(t *testing.T) {
	Convey("Contain on int", t, func() {
		So(Contain([]int{1, 2, 3}, 2), ShouldBeTrue)
		So(Contain([]int{1, 2, 3}, 4), ShouldBeFalse)
		So(Contain([]int(nil), 1), ShouldBeFalse)
	})
	Convey("Contain on string", t, func() {
		So(Contain([]string{"abc", "b"}, "abc"), ShouldBeTrue)
		So(Contain([]string{"abc", "b"}, "a"), ShouldBeFalse)
	})
}

func TestSetAdd(t *testing.T) {
	Convey("SetAdd 幂等去重", t, func() {
		s := SetAdd[int](nil, 1)
		So(len(s), ShouldEqual, 1)
		So(len(SetAdd(s, 1)), ShouldEqual, 1)
		So(len(SetAdd(s, 2, 3, 2)), ShouldEqual, 3)
	})
}

func TestWalk(t *testing.T) {
	Convey("Walk 映射 + 过滤", t, func() {
		src := []int{1, 2, 3}
		dest := Walk(src, func(v int) (int, bool) {
			return v * 10, v != 2
		})
		So(dest, ShouldResemble, []int{10, 30})
	})
	Convey("Walk on any (T 放宽到 any)", t, func() {
		src := []any{1, "a", true}
		dest := Walk(src, func(v any) (any, bool) { return v, true })
		So(len(dest), ShouldEqual, 3)
	})
}

func TestRemoveRepeated(t *testing.T) {
	Convey("RemoveRepeated 全局保序去重", t, func() {
		// 双层 loop 路径
		So(RemoveRepeated([]int{1, 2, 2, 3, 1}), ShouldResemble, []int{1, 2, 3})
		So(RemoveRepeated([]string{"a", "b", "a"}), ShouldResemble, []string{"a", "b"})
		So(len(RemoveRepeated([]int{})), ShouldEqual, 0)

		// map 路径：构造 len >= tooManyElement 触发
		big := make([]int, 0, tooManyElement+10)
		for i := 0; i < tooManyElement+10; i++ {
			big = append(big, i%5) // 0,1,2,3,4 反复
		}
		out := RemoveRepeated(big)
		So(out, ShouldResemble, []int{0, 1, 2, 3, 4})
	})
}

func TestRemoveEmpty(t *testing.T) {
	Convey("RemoveEmpty 移除零值", t, func() {
		So(RemoveEmpty([]int{1, 0, 2}), ShouldResemble, []int{1, 2})
		So(RemoveEmpty([]string{"abc", "", "b"}), ShouldResemble, []string{"abc", "b"})
		So(RemoveEmpty([]int{0, 0, 0}), ShouldResemble, []int{})
	})
}

func TestShuffle(t *testing.T) {
	Convey("Shuffle 原地打乱（保留全部元素）", t, func() {
		s := []int{1, 2, 3, 4, 5}
		Shuffle(s)
		So(len(s), ShouldEqual, 5)
		// 元素仍然是 {1,2,3,4,5} 的某个排列
		seen := make(map[int]bool)
		for _, v := range s {
			seen[v] = true
		}
		So(len(seen), ShouldEqual, 5)
	})
}

func TestToAny(t *testing.T) {
	Convey("ToAny", t, func() {
		So(ToAny([]int{1, 2}), ShouldResemble, []any{1, 2})
	})
}

func TestLast(t *testing.T) {
	Convey("Last 空切片返回零值", t, func() {
		So(Last([]int{1, 2, 3}), ShouldEqual, 3)
		So(Last([]int{}), ShouldEqual, 0)
		So(Last([]string{}), ShouldEqual, "")
	})
}

func TestTo(t *testing.T) {
	Convey("To 类型映射", t, func() {
		So(To([]int{1, 2, 3}, func(v int) string { return string(rune('a' + v - 1)) }),
			ShouldResemble, []string{"a", "b", "c"})
		So(To[int, string](nil, func(int) string { return "" }), ShouldBeNil)
	})
}

func TestStringHelpers(t *testing.T) {
	Convey("ContainEqualFold", t, func() {
		So(ContainEqualFold([]string{"ABC", "b"}, "abc"), ShouldBeTrue)
		So(ContainEqualFold([]string{"abc", "b"}, "x"), ShouldBeFalse)
	})
	Convey("AddPrefix / AddSuffix", t, func() {
		So(AddPrefix([]string{"a", "b"}, ">"), ShouldResemble, []string{">a", ">b"})
		So(AddSuffix([]string{"a", "b"}, "<"), ShouldResemble, []string{"a<", "b<"})
	})
}
