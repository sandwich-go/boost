// syncmap 提供了一个同步的映射实现，允许安全并发的访问
package syncmap

import (
	"errors"

	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestSyncMap(t *testing.T) {
	Convey("test sync map", t, func() {
		for _, tr := range []*SyncMap[int16, uint64]{New[int16, uint64]()} {
			So(tr.Len(), ShouldBeZeroValue)
			var k, v = int16(3), uint64(4)
			So(tr.Len(), ShouldEqual, 0)
			tr.Store(k, v)
			v1, ok := tr.Load(k)
			So(ok, ShouldBeTrue)
			So(v1, ShouldEqual, v)

			So(tr.Keys(), ShouldResemble, []int16{int16(3)})
			So(tr.Get(int16(3)), ShouldEqual, uint64(4))
			So(tr.Contains(int16(3)), ShouldBeTrue)

			tr.Store(int16(4), uint64(5))
			tr.Store(int16(5), uint64(6))
			ol := tr.Len()
			tr.DeleteMultiple(int16(4), int16(5))
			So(tr.Len(), ShouldEqual, ol-2)

			ol = tr.Len()
			tr.Store(int16(4), uint64(5))
			tr.Store(int16(5), uint64(6))
			vl, ok := tr.LoadAndDelete(int16(4))
			So(vl, ShouldEqual, uint64(5))
			So(ok, ShouldBeTrue)
			So(tr.Len(), ShouldEqual, ol+1)

			tr.Store(int16(4), uint64(5))
			fge := []func(key int16, cf func(key int16) (uint64, error)) (value uint64, loaded bool, err error){tr.GetOrSetFuncErrorLock}
			defv, defv2 := uint64(6), uint64(7)
			for _, f := range fge {
				v, l, e := f(int16(6), func(key int16) (uint64, error) {
					return defv, nil
				})
				So(v, ShouldEqual, defv)
				So(l, ShouldBeFalse)
				So(e, ShouldBeNil)

				v, l, e = f(int16(7), func(key int16) (uint64, error) {
					return defv2, errors.New("")
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
				So(e, ShouldNotBeNil)
			}
			fg := []func(key int16, cf func(key int16) uint64) (value uint64, loaded bool){tr.GetOrSetFuncLock}
			for _, f := range fg {
				v, l := f(int16(7), func(key int16) uint64 {
					return defv2
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
			}

			v, ok = tr.LoadOrStore(int16(8), uint64(9))
			So(v, ShouldEqual, uint64(9))
			So(ok, ShouldBeFalse)

			So(func() {
				tr.Range(func(key int16, value uint64) bool {
					return true
				})
			}, ShouldNotPanic)

		}
	})
}
