// Code generated by gotemplate. DO NOT EDIT.

// sset 包提供了多种类型的集合
// 可以产生一个带读写锁的线程安全的SyncSet，也可以产生一个非线程安全的SyncSet
// New 产生非协程安全的版本
// NewSync 产生协程安全的版本
package sset

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestInt(t *testing.T) {
	Convey("test sync set", t, func() {
		for _, tr := range []*Int{NewInt(), NewSyncInt()} {
			So(tr.Size(), ShouldBeZeroValue)
			var e0 = __formatToInt(30)
			tr.Add(e0)
			So(tr.Size(), ShouldEqual, 1)
			tr.Add(e0)
			So(tr.Size(), ShouldEqual, 1)

			So(tr.AddIfNotExist(__formatToInt(2)), ShouldBeTrue)
			So(tr.Size(), ShouldEqual, 2)

			tr.AddIfNotExistFunc(__formatToInt(3), func() bool {
				return false
			})
			So(tr.Size(), ShouldEqual, 2)
			tr.AddIfNotExistFunc(__formatToInt(3), func() bool {
				return true
			})
			So(tr.Size(), ShouldEqual, 3)

			tr.AddIfNotExistFuncLock(__formatToInt(4), func() bool {
				return false
			})
			So(tr.Size(), ShouldEqual, 3)
			tr.AddIfNotExistFuncLock(__formatToInt(4), func() bool {
				return true
			})
			So(tr.Size(), ShouldEqual, 4)

			So(tr.Contains(__formatToInt(4)), ShouldBeTrue)

			tr.Remove(__formatToInt(4))
			So(tr.Size(), ShouldEqual, 3)

			So(tr.Slice(), ShouldContain, __formatToInt(30))
			So(tr.Slice(), ShouldContain, __formatToInt(2))
			So(tr.Slice(), ShouldContain, __formatToInt(3))

			tr.Clear()
			So(tr.Size(), ShouldEqual, 0)

			tr.Add(__formatToInt(3), __formatToInt(2))
			tr2 := newWithSafeInt(false)
			tr2.Add(__formatToInt(3), __formatToInt(2))
			So(tr.Equal(tr2), ShouldBeTrue)

			tr3 := newWithSafeInt(true)
			tr3.Add(__formatToInt(3), __formatToInt(2))
			So(tr.Equal(tr3), ShouldBeTrue)

			tr4, tr5 := newWithSafeInt(true), newWithSafeInt(true)
			tr4.Add(__formatToInt(1), __formatToInt(4))
			tr5.Add(__formatToInt(1), __formatToInt(4))
			s := tr.Size()
			tr.Merge(tr4)
			tr2.Merge(tr5)
			So(tr.Equal(tr2), ShouldBeTrue)
			So(tr.Size(), ShouldEqual, s+2)

			So(func() {
				tr5.Walk(func(item int) int {
					return item
				})
			}, ShouldNotPanic)
		}
	})
}
