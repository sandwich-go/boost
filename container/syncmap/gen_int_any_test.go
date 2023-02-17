// Code generated by gotemplate. DO NOT EDIT.

// syncmap 提供了一个同步的映射实现，允许安全并发的访问
package syncmap

import (
	"errors"

	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestIntAny(t *testing.T) {
	Convey("test sync map", t, func() {
		for _, tr := range []*IntAny{NewIntAny()} {
			So(tr.Len(), ShouldBeZeroValue)
			var k, v = __formatKTypeToIntAny(3), __formatVTypeToIntAny(4)
			So(tr.Len(), ShouldEqual, 0)
			tr.Store(k, v)
			v1, ok := tr.Load(k)
			So(ok, ShouldBeTrue)
			So(v1, ShouldEqual, v)

			So(tr.Keys(), ShouldResemble, []int{__formatKTypeToIntAny(3)})
			So(tr.Get(__formatKTypeToIntAny(3)), ShouldEqual, __formatVTypeToIntAny(4))
			So(tr.Contains(__formatKTypeToIntAny(3)), ShouldBeTrue)

			tr.Store(__formatKTypeToIntAny(4), __formatVTypeToIntAny(5))
			tr.Store(__formatKTypeToIntAny(5), __formatVTypeToIntAny(6))
			ol := tr.Len()
			tr.DeleteMultiple(__formatKTypeToIntAny(4), __formatKTypeToIntAny(5))
			So(tr.Len(), ShouldEqual, ol-2)

			ol = tr.Len()
			tr.Store(__formatKTypeToIntAny(4), __formatVTypeToIntAny(5))
			tr.Store(__formatKTypeToIntAny(5), __formatVTypeToIntAny(6))
			vl, ok := tr.LoadAndDelete(__formatKTypeToIntAny(4))
			So(vl, ShouldEqual, __formatVTypeToIntAny(5))
			So(ok, ShouldBeTrue)
			So(tr.Len(), ShouldEqual, ol+1)

			tr.Store(__formatKTypeToIntAny(4), __formatVTypeToIntAny(5))
			fge := []func(key int, cf func(key int) (interface{}, error)) (value interface{}, loaded bool, err error){tr.GetOrSetFuncErrorLock}
			defv, defv2 := __formatVTypeToIntAny(6), __formatVTypeToIntAny(7)
			for _, f := range fge {
				v, l, e := f(__formatKTypeToIntAny(6), func(key int) (interface{}, error) {
					return defv, nil
				})
				So(v, ShouldEqual, defv)
				So(l, ShouldBeFalse)
				So(e, ShouldBeNil)

				v, l, e = f(__formatKTypeToIntAny(7), func(key int) (interface{}, error) {
					return defv2, errors.New("")
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
				So(e, ShouldNotBeNil)
			}
			fg := []func(key int, cf func(key int) interface{}) (value interface{}, loaded bool){tr.GetOrSetFuncLock}
			for _, f := range fg {
				v, l := f(__formatKTypeToIntAny(7), func(key int) interface{} {
					return defv2
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
			}

			v, ok = tr.LoadOrStore(__formatKTypeToIntAny(8), __formatVTypeToIntAny(9))
			So(v, ShouldEqual, __formatVTypeToIntAny(9))
			So(ok, ShouldBeFalse)

			So(func() {
				tr.Range(func(key int, value interface{}) bool {
					return true
				})
			}, ShouldNotPanic)

		}
	})
}
