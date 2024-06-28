// sset 包提供了多种类型的集合
// 可以产生一个带读写锁的线程安全的SyncSet，也可以产生一个非线程安全的SyncSet
// New 产生非协程安全的版本
// NewSync 产生协程安全的版本
package sset

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSyncSet(t *testing.T) {
	Convey("test sync set", t, func() {
		for _, tr := range []*Set[int16]{New[int16](), NewSync[int16]()} {
			So(tr.Size(), ShouldBeZeroValue)
			var e0 = int16(30)
			tr.Add(e0)
			So(tr.Size(), ShouldEqual, 1)
			tr.Add(e0)
			So(tr.Size(), ShouldEqual, 1)

			So(tr.AddIfNotExist(int16(2)), ShouldBeTrue)
			So(tr.Size(), ShouldEqual, 2)

			tr.AddIfNotExistFunc(int16(3), func() bool {
				return false
			})
			So(tr.Size(), ShouldEqual, 2)
			tr.AddIfNotExistFunc(int16(3), func() bool {
				return true
			})
			So(tr.Size(), ShouldEqual, 3)

			tr.AddIfNotExistFuncLock(int16(4), func() bool {
				return false
			})
			So(tr.Size(), ShouldEqual, 3)
			tr.AddIfNotExistFuncLock(int16(4), func() bool {
				return true
			})
			So(tr.Size(), ShouldEqual, 4)

			So(tr.Contains(int16(4)), ShouldBeTrue)

			tr.Remove(int16(4))
			So(tr.Size(), ShouldEqual, 3)

			So(tr.Slice(), ShouldContain, int16(30))
			So(tr.Slice(), ShouldContain, int16(2))
			So(tr.Slice(), ShouldContain, int16(3))

			tr.Clear()
			So(tr.Size(), ShouldEqual, 0)

			tr.Add(int16(3), int16(2))
			tr2 := newWithSafe[int16](false)
			tr2.Add(int16(3), int16(2))
			So(tr.Equal(tr2), ShouldBeTrue)

			tr3 := newWithSafe[int16](true)
			tr3.Add(int16(3), int16(2))
			So(tr.Equal(tr3), ShouldBeTrue)

			tr4, tr5 := newWithSafe[int16](true), newWithSafe[int16](true)
			tr4.Add(int16(1), int16(4))
			tr5.Add(int16(1), int16(4))
			s := tr.Size()
			tr.Merge(tr4)
			tr2.Merge(tr5)
			So(tr.Equal(tr2), ShouldBeTrue)
			So(tr.Size(), ShouldEqual, s+2)

			So(func() {
				tr5.Walk(func(item int16) int16 {
					return item
				})
			}, ShouldNotPanic)
		}
	})
}
