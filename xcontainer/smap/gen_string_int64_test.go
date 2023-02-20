// Code generated by gotemplate. DO NOT EDIT.

// smap 包提供了一个分片的协程安全的映射
// NewWithSharedCount 返回一个线程安全的映射实例
// New 返回一个线程安全的映射实例
package smap

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestSMapStringInt64(t *testing.T) {
	Convey("test sync array", t, func() {
		tr := NewStringInt64()
		So(tr.Len(), ShouldEqual, 0)
		So(tr.IsEmpty(), ShouldBeTrue)
		tr.Set(__formatKTypeToStringInt64(1), __formatVTypeToStringInt64(1))
		So(tr.Len(), ShouldEqual, 1)

		tr.Set(__formatKTypeToStringInt64(1), __formatVTypeToStringInt64(2))
		So(tr.Len(), ShouldEqual, 1)
		tr.Set(__formatKTypeToStringInt64(2), __formatVTypeToStringInt64(2))
		So(tr.Len(), ShouldEqual, 2)
		So(tr.Count(), ShouldEqual, 2)
		So(tr.Size(), ShouldEqual, 2)

		So(tr.Keys(), ShouldContain, __formatKTypeToStringInt64(1))
		So(tr.Keys(), ShouldContain, __formatKTypeToStringInt64(2))

		So(tr.GetAll(), ShouldContainKey, __formatKTypeToStringInt64(1))
		So(tr.GetAll(), ShouldContainKey, __formatKTypeToStringInt64(2))

		tr.Clear()
		So(tr.Len(), ShouldEqual, 0)

		tr.Set(__formatKTypeToStringInt64(1), __formatVTypeToStringInt64(2))
		tr.Set(__formatKTypeToStringInt64(2), __formatVTypeToStringInt64(2))
		So(func() {
			tr.ClearWithFuncLock(func(key string, val int64) {
				return
			})
		}, ShouldNotPanic)

		tr.Set(__formatKTypeToStringInt64(1), __formatVTypeToStringInt64(1))
		tr.Set(__formatKTypeToStringInt64(2), __formatVTypeToStringInt64(2))
		tr.Set(__formatKTypeToStringInt64(3), __formatVTypeToStringInt64(3))
		tr.Set(__formatKTypeToStringInt64(4), __formatVTypeToStringInt64(4))
		mk := []string{__formatKTypeToStringInt64(1), __formatKTypeToStringInt64(2), __formatKTypeToStringInt64(3)}
		m := tr.MGet(mk...)
		for _, k := range mk {
			So(m, ShouldContainKey, k)
		}

		tr2 := NewStringInt64()
		tr2.MSet(m)
		So(tr2.Len(), ShouldEqual, len(mk))

		So(tr2.SetNX(__formatKTypeToStringInt64(5), __formatVTypeToStringInt64(5)), ShouldBeTrue)
		So(tr2.SetNX(__formatKTypeToStringInt64(1), __formatVTypeToStringInt64(5)), ShouldBeFalse)

		So(func() {
			tr2.LockFuncWithKey(__formatKTypeToStringInt64(5), func(shardData map[string]int64) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFuncWithKey(__formatKTypeToStringInt64(5), func(shardData map[string]int64) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.LockFunc(func(shardData map[string]int64) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFunc(func(shardData map[string]int64) {
				return
			})
		}, ShouldNotPanic)

		dfv := __formatVTypeToStringInt64(1)
		r, ret := tr2.GetOrSetFunc(__formatKTypeToStringInt64(1), func(key string) int64 {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSetFuncLock(__formatKTypeToStringInt64(1), func(key string) int64 {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)

		_, ret = tr2.GetOrSet(__formatKTypeToStringInt64(1), __formatVTypeToStringInt64(1))
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSet(__formatKTypeToStringInt64(10), __formatVTypeToStringInt64(10))
		So(r, ShouldEqual, __formatVTypeToStringInt64(10))
		So(ret, ShouldBeTrue)

		So(tr.Has(__formatKTypeToStringInt64(1)), ShouldBeTrue)

		tr2.Remove(__formatKTypeToStringInt64(1))
		v, ret := tr2.GetAndRemove(__formatKTypeToStringInt64(10))
		So(v, ShouldEqual, __formatVTypeToStringInt64(10))
		So(ret, ShouldBeTrue)

		for _, f := range []func() <-chan TupleStringInt64{
			tr2.Iter, tr2.IterBuffered,
		} {
			cnt := 0
			for v := range f() {
				cnt++
				So(v.Key, ShouldBeIn, []string{__formatKTypeToStringInt64(2), __formatKTypeToStringInt64(3), __formatKTypeToStringInt64(5)})
				So(v.Val, ShouldBeIn, []int64{__formatVTypeToStringInt64(2), __formatVTypeToStringInt64(3), __formatVTypeToStringInt64(5)})
			}
			So(cnt, ShouldEqual, 3)
		}

	})
}
