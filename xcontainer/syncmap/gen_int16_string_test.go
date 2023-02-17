// Code generated by gotemplate. DO NOT EDIT.

// syncmap 提供了一个同步的映射实现，允许安全并发的访问
package syncmap

import (
	"errors"

	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestInt16String(t *testing.T) {
	Convey("test sync map", t, func() {
		for _, tr := range []*Int16String{NewInt16String()} {
			So(tr.Len(), ShouldBeZeroValue)
			var k, v = __formatKTypeToInt16String(3), __formatVTypeToInt16String(4)
			So(tr.Len(), ShouldEqual, 0)
			tr.Store(k, v)
			v1, ok := tr.Load(k)
			So(ok, ShouldBeTrue)
			So(v1, ShouldEqual, v)

			So(tr.Keys(), ShouldResemble, []int16{__formatKTypeToInt16String(3)})
			So(tr.Get(__formatKTypeToInt16String(3)), ShouldEqual, __formatVTypeToInt16String(4))
			So(tr.Contains(__formatKTypeToInt16String(3)), ShouldBeTrue)

			tr.Store(__formatKTypeToInt16String(4), __formatVTypeToInt16String(5))
			tr.Store(__formatKTypeToInt16String(5), __formatVTypeToInt16String(6))
			ol := tr.Len()
			tr.DeleteMultiple(__formatKTypeToInt16String(4), __formatKTypeToInt16String(5))
			So(tr.Len(), ShouldEqual, ol-2)

			ol = tr.Len()
			tr.Store(__formatKTypeToInt16String(4), __formatVTypeToInt16String(5))
			tr.Store(__formatKTypeToInt16String(5), __formatVTypeToInt16String(6))
			vl, ok := tr.LoadAndDelete(__formatKTypeToInt16String(4))
			So(vl, ShouldEqual, __formatVTypeToInt16String(5))
			So(ok, ShouldBeTrue)
			So(tr.Len(), ShouldEqual, ol+1)

			tr.Store(__formatKTypeToInt16String(4), __formatVTypeToInt16String(5))
			fge := []func(key int16, cf func(key int16) (string, error)) (value string, loaded bool, err error){tr.GetOrSetFuncErrorLock}
			defv, defv2 := __formatVTypeToInt16String(6), __formatVTypeToInt16String(7)
			for _, f := range fge {
				v, l, e := f(__formatKTypeToInt16String(6), func(key int16) (string, error) {
					return defv, nil
				})
				So(v, ShouldEqual, defv)
				So(l, ShouldBeFalse)
				So(e, ShouldBeNil)

				v, l, e = f(__formatKTypeToInt16String(7), func(key int16) (string, error) {
					return defv2, errors.New("")
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
				So(e, ShouldNotBeNil)
			}
			fg := []func(key int16, cf func(key int16) string) (value string, loaded bool){tr.GetOrSetFuncLock}
			for _, f := range fg {
				v, l := f(__formatKTypeToInt16String(7), func(key int16) string {
					return defv2
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
			}

			v, ok = tr.LoadOrStore(__formatKTypeToInt16String(8), __formatVTypeToInt16String(9))
			So(v, ShouldEqual, __formatVTypeToInt16String(9))
			So(ok, ShouldBeFalse)

			So(func() {
				tr.Range(func(key int16, value string) bool {
					return true
				})
			}, ShouldNotPanic)

		}
	})
}