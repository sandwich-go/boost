// Code generated by gotemplate. DO NOT EDIT.

// smap 包提供了一个分片的协程安全的映射
// NewWithSharedCount 返回一个线程安全的映射实例
// New 返回一个线程安全的映射实例
package smap

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestSMapStringAny(t *testing.T) {
	Convey("test sync array", t, func() {
		tr := NewStringAny()
		So(tr.Len(), ShouldEqual, 0)
		So(tr.IsEmpty(), ShouldBeTrue)
		tr.Set(__formatKTypeToStringAny(1), __formatVTypeToStringAny(1))
		So(tr.Len(), ShouldEqual, 1)

		tr.Set(__formatKTypeToStringAny(1), __formatVTypeToStringAny(2))
		So(tr.Len(), ShouldEqual, 1)
		tr.Set(__formatKTypeToStringAny(2), __formatVTypeToStringAny(2))
		So(tr.Len(), ShouldEqual, 2)
		So(tr.Count(), ShouldEqual, 2)
		So(tr.Size(), ShouldEqual, 2)

		So(tr.Keys(), ShouldContain, __formatKTypeToStringAny(1))
		So(tr.Keys(), ShouldContain, __formatKTypeToStringAny(2))

		So(tr.GetAll(), ShouldContainKey, __formatKTypeToStringAny(1))
		So(tr.GetAll(), ShouldContainKey, __formatKTypeToStringAny(2))

		tr.Clear()
		So(tr.Len(), ShouldEqual, 0)

		tr.Set(__formatKTypeToStringAny(1), __formatVTypeToStringAny(2))
		tr.Set(__formatKTypeToStringAny(2), __formatVTypeToStringAny(2))
		So(func() {
			tr.ClearWithFuncLock(func(key string, val interface{}) {
				return
			})
		}, ShouldNotPanic)

		tr.Set(__formatKTypeToStringAny(1), __formatVTypeToStringAny(1))
		tr.Set(__formatKTypeToStringAny(2), __formatVTypeToStringAny(2))
		tr.Set(__formatKTypeToStringAny(3), __formatVTypeToStringAny(3))
		tr.Set(__formatKTypeToStringAny(4), __formatVTypeToStringAny(4))
		mk := []string{__formatKTypeToStringAny(1), __formatKTypeToStringAny(2), __formatKTypeToStringAny(3)}
		m := tr.MGet(mk...)
		for _, k := range mk {
			So(m, ShouldContainKey, k)
		}

		tr2 := NewStringAny()
		tr2.MSet(m)
		So(tr2.Len(), ShouldEqual, len(mk))

		So(tr2.SetNX(__formatKTypeToStringAny(5), __formatVTypeToStringAny(5)), ShouldBeTrue)
		So(tr2.SetNX(__formatKTypeToStringAny(1), __formatVTypeToStringAny(5)), ShouldBeFalse)

		So(func() {
			tr2.LockFuncWithKey(__formatKTypeToStringAny(5), func(shardData map[string]interface{}) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFuncWithKey(__formatKTypeToStringAny(5), func(shardData map[string]interface{}) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.LockFunc(func(shardData map[string]interface{}) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFunc(func(shardData map[string]interface{}) {
				return
			})
		}, ShouldNotPanic)

		dfv := __formatVTypeToStringAny(1)
		r, ret := tr2.GetOrSetFunc(__formatKTypeToStringAny(1), func(key string) interface{} {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSetFuncLock(__formatKTypeToStringAny(1), func(key string) interface{} {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)

		_, ret = tr2.GetOrSet(__formatKTypeToStringAny(1), __formatVTypeToStringAny(1))
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSet(__formatKTypeToStringAny(10), __formatVTypeToStringAny(10))
		So(r, ShouldEqual, __formatVTypeToStringAny(10))
		So(ret, ShouldBeTrue)

		So(tr.Has(__formatKTypeToStringAny(1)), ShouldBeTrue)

		tr2.Remove(__formatKTypeToStringAny(1))
		v, ret := tr2.GetAndRemove(__formatKTypeToStringAny(10))
		So(v, ShouldEqual, __formatVTypeToStringAny(10))
		So(ret, ShouldBeTrue)

		for _, f := range []func() <-chan TupleStringAny{
			tr2.Iter, tr2.IterBuffered,
		} {
			cnt := 0
			for v := range f() {
				cnt++
				So(v.Key, ShouldBeIn, []string{__formatKTypeToStringAny(2), __formatKTypeToStringAny(3), __formatKTypeToStringAny(5)})
				So(v.Val, ShouldBeIn, []interface{}{__formatVTypeToStringAny(2), __formatVTypeToStringAny(3), __formatVTypeToStringAny(5)})
			}
			So(cnt, ShouldEqual, 3)
		}

	})
}
