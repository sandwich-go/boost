package xmap

import (
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDiff(t *testing.T) {
	Convey("Diff a == nil 直接返回 nil（不分配）", t, func() {
		So(Diff[string, int](nil, nil), ShouldBeNil)
		So(Diff[string, int](nil, map[string]int{"x": 1}), ShouldBeNil)
	})

	Convey("Diff b == nil 直接返回 a（不分配）", t, func() {
		a := map[string]int{"x": 1, "y": 2}
		// 注意：实现是 return a 不是 return copy(a) —— 调用方修改返回值会影响 a
		got := Diff(a, nil)
		So(got, ShouldEqual, a) // 同一个 map 引用
	})

	Convey("Diff a 与 b 完全相同返回空 map（不是 nil）", t, func() {
		a := map[string]int{"x": 1, "y": 2}
		b := map[string]int{"x": 1, "y": 2}
		got := Diff(a, b)
		So(got, ShouldNotBeNil) // 是 make() 出来的空 map
		So(got, ShouldBeEmpty)
	})

	Convey("Diff 找出 a 中 b 没有或值不同的 entry", t, func() {
		a := map[string]int{"x": 1, "y": 2, "z": 3}
		b := map[string]int{"x": 1, "y": 99} // y 值变了，z 缺失
		got := Diff(a, b)
		So(got, ShouldResemble, map[string]int{"y": 2, "z": 3})
	})

	Convey("Diff 不对称：只看 a 中的差异，b 中独有的 key 被忽略", t, func() {
		a := map[string]int{"x": 1}
		b := map[string]int{"x": 1, "y": 2}
		got := Diff(a, b)
		So(got, ShouldBeEmpty)
		// 反过来才会看到 y
		got2 := Diff(b, a)
		So(got2, ShouldResemble, map[string]int{"y": 2})
	})

	Convey("Diff 不修改入参 a / b", t, func() {
		a := map[int]int{1: 1, 2: 2}
		b := map[int]int{1: 1}
		_ = Diff(a, b)
		So(a, ShouldResemble, map[int]int{1: 1, 2: 2})
		So(b, ShouldResemble, map[int]int{1: 1})
	})
}

func TestToMap(t *testing.T) {
	Convey("ToMap from == nil 返回 nil", t, func() {
		got := ToMap[string, int, string, string](nil, func(k string, v int) (string, string) {
			return k, strconv.Itoa(v)
		})
		So(got, ShouldBeNil)
	})

	Convey("ToMap 空 map 返回非 nil 的空 map", t, func() {
		from := map[string]int{}
		got := ToMap(from, func(k string, v int) (string, string) {
			return k, strconv.Itoa(v)
		})
		So(got, ShouldNotBeNil)
		So(got, ShouldBeEmpty)
	})

	Convey("ToMap 同类型 key 转值类型", t, func() {
		from := map[string]int{"a": 1, "b": 2, "c": 3}
		got := ToMap(from, func(k string, v int) (string, string) {
			return k, strconv.Itoa(v)
		})
		So(got, ShouldResemble, map[string]string{"a": "1", "b": "2", "c": "3"})
	})

	Convey("ToMap 同时变 key 和 value 类型", t, func() {
		from := map[int]string{1: "a", 2: "b"}
		got := ToMap(from, func(k int, v string) (string, int) {
			return v, k
		})
		So(got, ShouldResemble, map[string]int{"a": 1, "b": 2})
	})

	Convey("ToMap callback 让多个原 key 映射到同一新 key（覆盖）", t, func() {
		from := map[int]int{1: 100, 2: 200, 3: 300}
		// 全部映射到 "all"，最后一个写入胜出
		got := ToMap(from, func(k int, v int) (string, int) {
			return "all", v
		})
		So(len(got), ShouldEqual, 1)
		So(got["all"], ShouldBeIn, 100, 200, 300)
	})
}

func TestEqual(t *testing.T) {
	Convey("Equal nil/empty 区分", t, func() {
		So(Equal[string, string](nil, nil), ShouldBeTrue)
		So(Equal(map[string]string{}, nil), ShouldBeFalse) // nil != empty
		So(Equal[string, string](nil, map[string]string{}), ShouldBeFalse)
		So(Equal(map[string]string{}, map[string]string{}), ShouldBeTrue)
	})
	Convey("Equal 长度/键值对", t, func() {
		a := map[int]string{1: "a", 2: "b"}
		b := map[int]string{1: "a", 2: "b"}
		c := map[int]string{1: "a", 2: "c"}
		d := map[int]string{1: "a"}
		So(Equal(a, b), ShouldBeTrue)
		So(Equal(a, c), ShouldBeFalse)
		So(Equal(a, d), ShouldBeFalse)
	})
	Convey("Equal 大量数据", t, func() {
		a := make(map[int]string, 100)
		b := make(map[int]string, 100)
		for i := 0; i < 100; i++ {
			a[i] = strconv.Itoa(i)
			b[i] = strconv.Itoa(i)
		}
		So(Equal(a, b), ShouldBeTrue)
		delete(a, 50)
		So(Equal(a, b), ShouldBeFalse)
	})
}

func TestWalkMapDeterministic(t *testing.T) {
	Convey("WalkMapDeterministic 按 key 升序", t, func() {
		in := map[int]string{3: "c", 1: "a", 2: "b"}
		var keys []int
		WalkMapDeterministic(in, func(k int, v string) bool {
			keys = append(keys, k)
			return true
		})
		So(keys, ShouldResemble, []int{1, 2, 3})
	})
	Convey("WalkMapDeterministic V=any (放宽约束的关键场景)", t, func() {
		in := map[string]any{"b": 2, "a": 1, "c": "three"}
		var keys []string
		WalkMapDeterministic(in, func(k string, v any) bool {
			keys = append(keys, k)
			return true
		})
		So(keys, ShouldResemble, []string{"a", "b", "c"})
	})
	Convey("WalkMapDeterministic walkFunc 返回 false 停止", t, func() {
		in := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		var visited []int
		WalkMapDeterministic(in, func(k, v int) bool {
			visited = append(visited, k)
			return k < 2 // 访问到 k=2 时返回 false 停止
		})
		So(visited, ShouldResemble, []int{1, 2})
	})
	Convey("WalkMapDeterministic V=func (不可比较类型)", t, func() {
		// 关键：V=any 让我们能放函数等不可比较类型
		in := map[int]func() int{2: func() int { return 2 }, 1: func() int { return 1 }}
		var keys []int
		WalkMapDeterministic(in, func(k int, v func() int) bool {
			keys = append(keys, k)
			return true
		})
		So(keys, ShouldResemble, []int{1, 2})
	})
}
