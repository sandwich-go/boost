// Code generated by gotemplate. DO NOT EDIT.

package sarray

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestAny(t *testing.T) {
	Convey("test sync array", t, func() {
		for _, tr := range []*Any{NewAny(), NewSyncAny()} {
			So(tr.Len(), ShouldBeZeroValue)
			_, exists := tr.Get(0)
			So(exists, ShouldBeFalse)
			So(tr.Empty(), ShouldBeTrue)
			var e0 = __formatToAny(3)
			tr.PushLeft(e0)
			So(tr.Len(), ShouldEqual, 1)

			e := tr.At(0)
			So(e, ShouldEqual, e0)
			v, f := tr.Get(0)
			So(f, ShouldBeTrue)
			So(v, ShouldEqual, e)
			_, f = tr.Get(1)
			So(f, ShouldBeFalse)

			tr.SetArray([]interface{}{__formatToAny(1), __formatToAny(1), __formatToAny(1), __formatToAny(1), __formatToAny(1)})
			So(tr.Len(), ShouldEqual, 5)
			tr.SetArray([]interface{}{__formatToAny(1), __formatToAny(2), __formatToAny(3), __formatToAny(4), __formatToAny(5)})
			So(tr.Len(), ShouldEqual, 5)

			rpv := __formatToAny(20)
			err := tr.Set(2, rpv)
			So(err, ShouldBeEmpty)
			So(tr.At(2), ShouldEqual, rpv)

			tr.Replace([]interface{}{__formatToAny(1), __formatToAny(1)})
			So(tr.Len(), ShouldEqual, 5)
			So(tr.At(0), ShouldEqual, tr.At(1))
			So(tr.At(2), ShouldNotEqual, tr.At(0))

			iv1 := __formatToAny(11)
			err = tr.InsertBefore(0, iv1)
			So(err, ShouldBeNil)
			So(tr.At(0), ShouldEqual, iv1)

			iv2 := __formatToAny(12)
			err = tr.InsertAfter(0, iv2)
			So(err, ShouldBeNil)
			So(tr.At(1), ShouldEqual, iv2)

			So(tr.Contains(iv1), ShouldBeTrue)
			So(tr.Search(iv1), ShouldNotEqual, -1)

			So(tr.DeleteValue(iv2), ShouldBeTrue)
			v, f = tr.LoadAndDelete(0)
			So(f, ShouldBeTrue)
			So(v, ShouldEqual, iv1)

			pl := __formatToAny(11)
			tr.PushLeft(pl)
			So(tr.At(0), ShouldEqual, pl)
			pr := __formatToAny(21)
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

			aps := []interface{}{__formatToAny(35), __formatToAny(40), __formatToAny(45), __formatToAny(50)}
			tr.Append(aps...)
			for i := len(aps); i > 0; i-- {
				So(aps[i-1], ShouldEqual, func() interface{} { v, f = tr.PopRight(); So(f, ShouldBeTrue); return v }())
			}

			tr.Clear()
			So(tr.Len(), ShouldEqual, 0)

			tr.Append(aps...)
			s := tr.Slice()
			So(len(s), ShouldEqual, tr.Len())

			k := 0
			tr.WalkAsc(func(key int, val interface{}) bool {
				So(key, ShouldEqual, k)
				So(val, ShouldEqual, s[k])
				k++
				return true
			})

			k = len(s) - 1
			tr.WalkDesc(func(key int, val interface{}) bool {
				So(key, ShouldEqual, k)
				So(val, ShouldEqual, s[k])
				k--
				return true
			})

			So(func() {
				tr.LockFunc(func(array []interface{}) {
					return
				})
			}, ShouldNotPanic)

			So(func() {
				tr.RLockFunc(func(array []interface{}) {
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
