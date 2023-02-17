// Code generated by gotemplate. DO NOT EDIT.

package smap

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestSMapInt64String(t *testing.T) {
	Convey("test sync array", t, func() {
		tr := NewInt64String()
		So(tr.Len(), ShouldEqual, 0)
		So(tr.IsEmpty(), ShouldBeTrue)
		tr.Set(__formatKTypeToInt64String(1), __formatVTypeToInt64String(1))
		So(tr.Len(), ShouldEqual, 1)

		tr.Set(__formatKTypeToInt64String(1), __formatVTypeToInt64String(2))
		So(tr.Len(), ShouldEqual, 1)
		tr.Set(__formatKTypeToInt64String(2), __formatVTypeToInt64String(2))
		So(tr.Len(), ShouldEqual, 2)

		So(tr.Keys(), ShouldContain, __formatKTypeToInt64String(1))
		So(tr.Keys(), ShouldContain, __formatKTypeToInt64String(2))

		So(tr.GetAll(), ShouldContainKey, __formatKTypeToInt64String(1))
		So(tr.GetAll(), ShouldContainKey, __formatKTypeToInt64String(2))

		tr.Clear()
		So(tr.Len(), ShouldEqual, 0)

		tr.Set(__formatKTypeToInt64String(1), __formatVTypeToInt64String(2))
		tr.Set(__formatKTypeToInt64String(2), __formatVTypeToInt64String(2))
		So(func() {
			tr.ClearWithFuncLock(func(key int64, val string) {
				return
			})
		}, ShouldNotPanic)

		tr.Set(__formatKTypeToInt64String(1), __formatVTypeToInt64String(1))
		tr.Set(__formatKTypeToInt64String(2), __formatVTypeToInt64String(2))
		tr.Set(__formatKTypeToInt64String(3), __formatVTypeToInt64String(3))
		tr.Set(__formatKTypeToInt64String(4), __formatVTypeToInt64String(4))
		mk := []int64{__formatKTypeToInt64String(1), __formatKTypeToInt64String(2), __formatKTypeToInt64String(3)}
		m := tr.MGet(mk...)
		for _, k := range mk {
			So(m, ShouldContainKey, k)
		}

		tr2 := NewInt64String()
		tr2.MSet(m)
		So(tr2.Len(), ShouldEqual, len(mk))

		So(tr2.SetNX(__formatKTypeToInt64String(5), __formatVTypeToInt64String(5)), ShouldBeTrue)
		So(tr2.SetNX(__formatKTypeToInt64String(1), __formatVTypeToInt64String(5)), ShouldBeFalse)

		So(func() {
			tr2.LockFuncWithKey(__formatKTypeToInt64String(5), func(shardData map[int64]string) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFuncWithKey(__formatKTypeToInt64String(5), func(shardData map[int64]string) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.LockFunc(func(shardData map[int64]string) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFunc(func(shardData map[int64]string) {
				return
			})
		}, ShouldNotPanic)

		dfv := __formatVTypeToInt64String(1)
		r, ret := tr2.GetOrSetFunc(__formatKTypeToInt64String(1), func(key int64) string {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSetFuncLock(__formatKTypeToInt64String(1), func(key int64) string {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)

		_, ret = tr2.GetOrSet(__formatKTypeToInt64String(1), __formatVTypeToInt64String(1))
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSet(__formatKTypeToInt64String(10), __formatVTypeToInt64String(10))

		So(ret, ShouldBeTrue)
	})
}
