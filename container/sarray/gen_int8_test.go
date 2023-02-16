// Code generated by gotemplate. DO NOT EDIT.

package sarray

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestInt8(t *testing.T) {
	Convey("test sync array", t, func() {
		for _, tr := range []*Int8{NewInt8(), NewSyncInt8()} {
			So(tr.Len(), ShouldBeZeroValue)
			_, exists := tr.Get(0)
			So(exists, ShouldBeFalse)
			So(tr.Empty(), ShouldBeTrue)
			var e0 = __formatToInt8(3)
			tr.PushLeft(e0)
			So(tr.Len(), ShouldEqual, 1)

			e := tr.At(0)
			So(e, ShouldEqual, e0)
			v, f := tr.Get(0)
			So(f, ShouldBeTrue)
			So(v, ShouldEqual, e)
			_, f = tr.Get(1)
			So(f, ShouldBeFalse)

			tr.SetArray([]int8{__formatToInt8(1), __formatToInt8(1), __formatToInt8(1), __formatToInt8(1), __formatToInt8(1)})
			So(tr.Len(), ShouldEqual, 5)
			tr.SetArray([]int8{__formatToInt8(1), __formatToInt8(2), __formatToInt8(3), __formatToInt8(4), __formatToInt8(5)})
			So(tr.Len(), ShouldEqual, 5)

			rpv := __formatToInt8(20)
			err := tr.Set(2, rpv)
			So(err, ShouldBeEmpty)
			So(tr.At(2), ShouldEqual, rpv)

			tr.Replace([]int8{__formatToInt8(1), __formatToInt8(1)})
			So(tr.Len(), ShouldEqual, 5)
			So(tr.At(0), ShouldEqual, tr.At(1))
			So(tr.At(2), ShouldNotEqual, tr.At(0))

			iv1 := __formatToInt8(11)
			err = tr.InsertBefore(0, iv1)
			So(err, ShouldBeNil)
			So(tr.At(0), ShouldEqual, iv1)

			iv2 := __formatToInt8(12)
			err = tr.InsertAfter(0, iv2)
			So(err, ShouldBeNil)
			So(tr.At(1), ShouldEqual, iv2)

			So(tr.Contains(iv1), ShouldBeTrue)
			So(tr.Search(iv1), ShouldNotEqual, -1)

			So(tr.DeleteValue(iv2), ShouldBeTrue)
			v, f = tr.LoadAndDelete(0)
			So(f, ShouldBeTrue)
			So(v, ShouldEqual, iv1)

			pl := __formatToInt8(11)
			tr.PushLeft(pl)
			So(tr.At(0), ShouldEqual, pl)
			pr := __formatToInt8(21)
			tr.PushRight(pr)
			So(tr.At(tr.Len()-1), ShouldEqual, pr)

			v, f = tr.PopLeft()
			So(v, ShouldEqual, pl)
			v, f = tr.PopRight()
			So(v, ShouldEqual, pr)
			l := tr.Len()
			_, f = tr.PopRand()
			So(f, ShouldBeTrue)
			So(tr.Len()+1, ShouldEqual, l)
			l = tr.Len()
			poplen := 2
			pv := tr.PopRands(poplen)
			So(len(pv), ShouldEqual, poplen)
			So(tr.Len(), ShouldBeGreaterThanOrEqualTo, l-poplen)

			aps := []int8{__formatToInt8(35), __formatToInt8(40), __formatToInt8(45), __formatToInt8(50)}
			tr.Append(aps...)
			for i := len(aps); i > 0; i-- {
				So(aps[i-1], ShouldEqual, func() int8 { v, f = tr.PopRight(); So(f, ShouldBeTrue); return v }())
			}

			tr.Clear()
			So(tr.Len(), ShouldEqual, 0)

			tr.Append(aps...)
			s := tr.Slice()
			So(len(s), ShouldEqual, tr.Len())

			k := 0
			tr.WalkAsc(func(key int, val int8) bool {
				So(key, ShouldEqual, k)
				So(val, ShouldEqual, s[k])
				k++
				return true
			})

			k = len(s) - 1
			tr.WalkDesc(func(key int, val int8) bool {
				So(key, ShouldEqual, k)
				So(val, ShouldEqual, s[k])
				k--
				return true
			})

			So(func() {
				tr.LockFunc(func(array []int8) {
					return
				})
			}, ShouldNotPanic)

			So(func() {
				tr.RLockFunc(func(array []int8) {
					return
				})
			}, ShouldNotPanic)

			j, err := tr.MarshalJSON()
			So(err, ShouldBeNil)
			err = tr.UnmarshalJSON(j)
			So(err, ShouldBeNil)
		}
	})
}
