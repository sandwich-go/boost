package smap

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

type ss string
type ii int

func TestAnyInt32String(t *testing.T) {
	Convey("test sync array should work ok", t, func() {
		tr0 := New[ss, ii]()
		tr1 := New[ii, ss]()

		So(tr0.Len(), ShouldEqual, 0)
		So(tr1.Len(), ShouldEqual, 0)
		tr0.Set("3", 1)
		tr1.Set(1, "3")

		So(tr0.Len(), ShouldEqual, 1)
		So(tr1.Len(), ShouldEqual, 1)

		So(tr0.SetNX("4", 2), ShouldBeTrue)
		So(tr0.SetNX("3", 2), ShouldBeFalse)

		So(tr1.SetNX(2, "4"), ShouldBeTrue)
		So(tr1.SetNX(1, "3"), ShouldBeFalse)

		So(tr0.GetAll(), ShouldContainKey, ss("3"))
		So(tr0.GetAll(), ShouldContainKey, ss("4"))

		So(tr1.GetAll(), ShouldContainKey, ii(1))
		So(tr1.GetAll(), ShouldContainKey, ii(2))
	})
	Convey("test sync array", t, func() {
		tr := New[int32, string]()
		So(tr.Len(), ShouldEqual, 0)
		So(tr.IsEmpty(), ShouldBeTrue)
		tr.Set(__formatKTypeToInt32String(1), __formatVTypeToInt32String(1))
		So(tr.Len(), ShouldEqual, 1)

		tr.Set(__formatKTypeToInt32String(1), __formatVTypeToInt32String(2))
		So(tr.Len(), ShouldEqual, 1)
		tr.Set(__formatKTypeToInt32String(2), __formatVTypeToInt32String(2))
		So(tr.Len(), ShouldEqual, 2)
		So(tr.Count(), ShouldEqual, 2)
		So(tr.Size(), ShouldEqual, 2)

		So(tr.Keys(), ShouldContain, __formatKTypeToInt32String(1))
		So(tr.Keys(), ShouldContain, __formatKTypeToInt32String(2))

		So(tr.GetAll(), ShouldContainKey, __formatKTypeToInt32String(1))
		So(tr.GetAll(), ShouldContainKey, __formatKTypeToInt32String(2))

		tr.Clear()
		So(tr.Len(), ShouldEqual, 0)

		tr.Set(__formatKTypeToInt32String(1), __formatVTypeToInt32String(2))
		tr.Set(__formatKTypeToInt32String(2), __formatVTypeToInt32String(2))
		So(func() {
			tr.ClearWithFuncLock(func(key int32, val string) {
				return
			})
		}, ShouldNotPanic)

		tr.Set(__formatKTypeToInt32String(1), __formatVTypeToInt32String(1))
		tr.Set(__formatKTypeToInt32String(2), __formatVTypeToInt32String(2))
		tr.Set(__formatKTypeToInt32String(3), __formatVTypeToInt32String(3))
		tr.Set(__formatKTypeToInt32String(4), __formatVTypeToInt32String(4))
		mk := []int32{__formatKTypeToInt32String(1), __formatKTypeToInt32String(2), __formatKTypeToInt32String(3)}
		m := tr.MGet(mk...)
		for _, k := range mk {
			So(m, ShouldContainKey, k)
		}

		tr2 := New[int32, string]()
		tr2.MSet(m)
		So(tr2.Len(), ShouldEqual, len(mk))

		So(tr2.SetNX(__formatKTypeToInt32String(5), __formatVTypeToInt32String(5)), ShouldBeTrue)
		So(tr2.SetNX(__formatKTypeToInt32String(1), __formatVTypeToInt32String(5)), ShouldBeFalse)

		So(func() {
			tr2.LockFuncWithKey(__formatKTypeToInt32String(5), func(shardData map[int32]string) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFuncWithKey(__formatKTypeToInt32String(5), func(shardData map[int32]string) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.LockFunc(func(shardData map[int32]string) {
				return
			})
		}, ShouldNotPanic)
		So(func() {
			tr2.RLockFunc(func(shardData map[int32]string) {
				return
			})
		}, ShouldNotPanic)

		dfv := __formatVTypeToInt32String(1)
		r, ret := tr2.GetOrSetFunc(__formatKTypeToInt32String(1), func(key int32) string {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSetFuncLock(__formatKTypeToInt32String(1), func(key int32) string {
			return dfv
		})
		So(r, ShouldEqual, dfv)
		So(ret, ShouldBeFalse)

		_, ret = tr2.GetOrSet(__formatKTypeToInt32String(1), __formatVTypeToInt32String(1))
		So(ret, ShouldBeFalse)
		r, ret = tr2.GetOrSet(__formatKTypeToInt32String(10), __formatVTypeToInt32String(10))
		So(r, ShouldEqual, __formatVTypeToInt32String(10))
		So(ret, ShouldBeTrue)

		So(tr.Has(__formatKTypeToInt32String(1)), ShouldBeTrue)

		tr2.Remove(__formatKTypeToInt32String(1))
		v, ret := tr2.GetAndRemove(__formatKTypeToInt32String(10))
		So(v, ShouldEqual, __formatVTypeToInt32String(10))
		So(ret, ShouldBeTrue)

		for _, f := range []func() <-chan Tuple[int32, string]{
			tr2.Iter, tr2.IterBuffered,
		} {
			cnt := 0
			for v := range f() {
				cnt++
				So(v.Key, ShouldBeIn, []int32{__formatKTypeToInt32String(2), __formatKTypeToInt32String(3), __formatKTypeToInt32String(5)})
				So(v.Val, ShouldBeIn, []string{__formatVTypeToInt32String(2), __formatVTypeToInt32String(3), __formatVTypeToInt32String(5)})
			}
			So(cnt, ShouldEqual, 3)
		}

	})
}
