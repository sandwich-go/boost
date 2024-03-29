// Code generated by gotemplate. DO NOT EDIT.

// smap 包提供了一个分片的协程安全的映射
// NewWithSharedCount 返回一个线程安全的映射实例
// New 返回一个线程安全的映射实例
package smap

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestSMapIntInt(t *testing.T) {
	Convey("test sync array", t, func() {
		tr := NewIntInt()
		So(tr.Len(), ShouldEqual, 0)
		So(tr.IsEmpty(), ShouldBeTrue)
		tr.Set(__formatKTypeToIntInt(1), __formatVTypeToIntInt(1))
		So(tr.Len(), ShouldEqual, 1)

		tr.Set(__formatKTypeToIntInt(1), __formatVTypeToIntInt(2))
		So(tr.Len(), ShouldEqual, 1)
		tr.Set(__formatKTypeToIntInt(2), __formatVTypeToIntInt(2))
		So(tr.Len(), ShouldEqual, 2)
		So(tr.Count(), ShouldEqual, 2)
		So(tr.Size(), ShouldEqual, 2)

		So(tr.Keys(), ShouldContain, __formatKTypeToIntInt(1))
		So(tr.Keys(), ShouldContain, __formatKTypeToIntInt(2))

		So(tr.GetAll(), ShouldContainKey, __formatKTypeToIntInt(1))
		So(tr.GetAll(), ShouldContainKey, __formatKTypeToIntInt(2))

		tr.Clear()
		So(tr.Len(), ShouldEqual, 0)

		tr.Set(__formatKTypeToIntInt(1), __formatVTypeToIntInt(2))
		tr.Set(__formatKTypeToIntInt(2), __formatVTypeToIntInt(2))
		So(func() {
			tr.ClearWithFuncLock(func(key int, val int) {
				return
			})
		}, ShouldNotPanic)

		tr.Set(__formatKTypeToIntInt(1), __formatVTypeToIntInt(1))
		tr.Set(__formatKTypeToIntInt(2), __formatVTypeToIntInt(2))
		tr.Set(__formatKTypeToIntInt(3), __formatVTypeToIntInt(3))
		tr.Set(__formatKTypeToIntInt(4), __formatVTypeToIntInt(4))
		mk := []int{__formatKTypeToIntInt(1), __formatKTypeToIntInt(2), __formatKTypeToIntInt(3)}
		m := tr.MGet(mk...)
		for _, k := range mk {
			So(m, ShouldContainKey, k)
		}

		tr2 := NewIntInt()
		tr2.MSet(m)
		So(tr2.Len(), ShouldEqual, len(mk))

		So(tr2.SetNX(__formatKTypeToIntInt(5), __formatVTypeToIntInt(5)), ShouldBeTrue)
		So(tr2.SetNX(__formatKTypeToIntInt(1), __formatVTypeToIntInt(5)), ShouldBeFalse)

		So(func() {
			tr2.LockFuncWithKey(__formatKTypeToIntInt(5), func(shardData map[int]int) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFuncWithKey(__formatKTypeToIntInt(5), func(shardData map[int]int) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.LockFunc(func(shardData map[int]int) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFunc(func(shardData map[int]int) {
				return
			})
		}, ShouldNotPanic)

		dfv := __formatVTypeToIntInt(1)
		r, ret := tr2.GetOrSetFunc(__formatKTypeToIntInt(1), func(key int) int {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSetFuncLock(__formatKTypeToIntInt(1), func(key int) int {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)

		_, ret = tr2.GetOrSet(__formatKTypeToIntInt(1), __formatVTypeToIntInt(1))
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSet(__formatKTypeToIntInt(10), __formatVTypeToIntInt(10))
		So(r, ShouldEqual, __formatVTypeToIntInt(10))
		So(ret, ShouldBeTrue)

		So(tr.Has(__formatKTypeToIntInt(1)), ShouldBeTrue)

		tr2.Remove(__formatKTypeToIntInt(1))
		v, ret := tr2.GetAndRemove(__formatKTypeToIntInt(10))
		So(v, ShouldEqual, __formatVTypeToIntInt(10))
		So(ret, ShouldBeTrue)

		for _, f := range []func() <-chan TupleIntInt{
			tr2.Iter, tr2.IterBuffered,
		} {
			cnt := 0
			for v := range f() {
				cnt++
				So(v.Key, ShouldBeIn, []int{__formatKTypeToIntInt(2), __formatKTypeToIntInt(3), __formatKTypeToIntInt(5)})
				So(v.Val, ShouldBeIn, []int{__formatVTypeToIntInt(2), __formatVTypeToIntInt(3), __formatVTypeToIntInt(5)})
			}
			So(cnt, ShouldEqual, 3)
		}

	})
}
