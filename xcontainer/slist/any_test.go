// slist 包提供了一个同步的链表实现
// 可以产生一个带读写锁的线程安全的SyncList，也可以产生一个非线程安全的SyncList
// New 产生非协程安全的版本
// NewSync 产生协程安全的版本
package slist

import (
	"container/list"

	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestSyncList(t *testing.T) {
	Convey("test sync list", t, func() {
		for _, tr := range []*List[int8]{New[int8](), NewSync[int8]()} {
			So(tr.Len(), ShouldBeZeroValue)
			var e0 = int8(3)
			tr.PushBack(e0)
			So(tr.Len(), ShouldEqual, 1)
			e := tr.PopBack()
			So(e, ShouldEqual, e0)

			tr.PushFront(e0)
			So(tr.Len(), ShouldEqual, 1)
			e = tr.PopFront()
			So(e, ShouldEqual, e0)

			ps := []int8{int8(1), int8(1), int8(3), int8(4), int8(5)}
			tr.PushFronts(ps)
			So(tr.Len(), ShouldEqual, len(ps))
			pops := tr.PopFronts(2)
			So(len(pops), ShouldEqual, 2)
			pops = tr.PopFrontAll()
			So(len(pops), ShouldEqual, 3)

			ps = []int8{int8(1), int8(2), int8(3), int8(4)}
			tr.PushBacks(ps)
			So(tr.Len(), ShouldEqual, len(ps))
			pops = tr.PopBacks(2)
			So(len(pops), ShouldEqual, 2)
			pops = tr.PopBackAll()
			So(len(pops), ShouldEqual, 2)

			ps = []int8{int8(1), int8(2), int8(3), int8(4)}
			tr.Clear()
			tr.PushBacks(ps)
			So(tr.FrontAll(), ShouldResemble, ps)
			tr.Clear()

			psrev := []int8{int8(4), int8(3), int8(2), int8(1)}
			tr.PushBacks(ps)
			So(tr.BackAll(), ShouldResemble, psrev)

			So(tr.FrontValue(), ShouldEqual, int8(1))
			So(tr.Front().Value, ShouldEqual, int8(1))
			So(tr.BackValue(), ShouldEqual, int8(4))
			So(tr.Back().Value, ShouldEqual, int8(4))

			b, b1 := tr.Back(), tr.Back().Prev()
			tr.MoveBefore(tr.Back(), tr.Front())
			So(tr.Front(), ShouldEqual, b)
			So(tr.Back(), ShouldEqual, b1)

			f0, f1 := tr.Front(), tr.Front().Next()
			tr.MoveAfter(tr.Front(), tr.Back())
			So(tr.Back(), ShouldEqual, f0)
			So(tr.Front(), ShouldEqual, f1)

			b, b1 = tr.Back(), tr.Back().Prev()
			tr.MoveToFront(tr.Back())
			So(tr.Front(), ShouldEqual, b)
			So(tr.Back(), ShouldEqual, b1)

			f0, f1 = tr.Front(), tr.Front().Next()
			tr.MoveToBack(tr.Front())
			So(tr.Back(), ShouldEqual, f0)
			So(tr.Front(), ShouldEqual, f1)

			n, ns, ol := New[int8](), NewSync[int8](), tr.Len()
			n.PushFronts([]int8{int8(1), int8(2)})
			ns.PushFronts([]int8{int8(1), int8(2)})
			tr.PushFrontList(n)
			So(tr.Len(), ShouldEqual, ol+2)
			tr.PushFrontList(ns)
			So(tr.Len(), ShouldEqual, ol+2+2)

			f0, trl := tr.Front(), tr.Len()
			tr.InsertBefore(tr.Front(), int8(10))
			So(tr.Front().Next(), ShouldEqual, f0)
			So(tr.Front().Value, ShouldEqual, int8(10))
			So(tr.Len(), ShouldEqual, trl+1)

			b, trl = tr.Back(), tr.Len()
			tr.InsertAfter(tr.Back(), int8(10))
			So(tr.Back().Prev(), ShouldEqual, b)
			So(tr.Back().Value, ShouldEqual, int8(10))
			So(tr.Len(), ShouldEqual, trl+1)

			bv := tr.Back().Value
			So(tr.Remove(tr.Back()), ShouldEqual, bv)

			So(func() { tr.Removes([]*list.Element{tr.Front(), tr.Front().Next()}) }, ShouldNotPanic)
			So(func() { tr.RemoveAll() }, ShouldNotPanic)
			So(tr.Len(), ShouldEqual, 0)

			tr.PushFrontList(n)
			tr.Clear()
			So(tr.Len(), ShouldEqual, 0)

			tr.PushFronts([]int8{int8(10), int8(20), int8(30), int8(40)})

			So(func() {
				tr.RLockFunc(func(list *list.List) {
					So(list.Front().Value, ShouldEqual, int8(40))
				})
			}, ShouldNotPanic)

			So(func() {
				tr.LockFunc(func(list *list.List) {
					So(list.Front().Value, ShouldEqual, int8(40))
				})
			}, ShouldNotPanic)

			So(func() {
				tr.Iterator(func(e *list.Element) bool {
					return true
				})
			}, ShouldNotPanic)

			So(func() {
				tr.IteratorAsc(func(e *list.Element) bool {
					return true
				})
			}, ShouldNotPanic)

			So(func() {
				tr.IteratorDesc(func(e *list.Element) bool {
					return true
				})
			}, ShouldNotPanic)
		}
	})
}
