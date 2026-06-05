package xmap

import (
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
