// slist 包提供了一个同步的链表实现
// 可以产生一个带读写锁的线程安全的SyncList，也可以产生一个非线程安全的SyncList
// New 产生非协程安全的版本
// NewSync 产生协程安全的版本
package slist

import (
	"container/list"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
)

//template type SyncList(VType)

type VType interface{}
type Element = list.Element

// SyncList 包含一个读写锁和一个双向链表，根据不同需求可提供对切片协程安全版本或者非协程安全版本的实例
type SyncList struct {
	mu   *localRWMutexVType
	list *list.List
}

func newWithSafe(safe bool) *SyncList {
	return &SyncList{
		mu:   newLocalRWMutexVType(safe),
		list: list.New(),
	}
}

// New 创建非协程安全版本
func New() *SyncList { return newWithSafe(false) }

// NewSync 创建协程安全版本
func NewSync() *SyncList { return newWithSafe(true) }

// PushFront 队头添加
func (l *SyncList) PushFront(v VType) (e *Element) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushFront(v)
	l.mu.Unlock()
	return
}

// PushBack 队尾添加
func (l *SyncList) PushBack(v VType) (e *Element) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushBack(v)
	l.mu.Unlock()
	return
}

// PushFronts 队头添加多个元素
func (l *SyncList) PushFronts(values []VType) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	for _, v := range values {
		l.list.PushFront(v)
	}
	l.mu.Unlock()
}

// PushBacks 队尾添加多个元素
func (l *SyncList) PushBacks(values []VType) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	for _, v := range values {
		l.list.PushBack(v)
	}
	l.mu.Unlock()
}

// PopBack 队尾弹出元素
func (l *SyncList) PopBack() (value VType) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Back(); e != nil {
		value = (l.list.Remove(e)).(VType)
	}
	return
}

// PopFront 队头弹出元素
func (l *SyncList) PopFront() (value VType) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Front(); e != nil {
		value = l.list.Remove(e).(VType)
	}
	return
}

func (l *SyncList) pops(max int, front bool) (values []VType) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	length := l.list.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]VType, length)
		for i := 0; i < length; i++ {
			if front {
				values[i] = l.list.Remove(l.list.Front()).(VType)
			} else {
				values[i] = l.list.Remove(l.list.Back()).(VType)
			}
		}
	}
	return
}

// PopBacks 队尾弹出至多max个元素
func (l *SyncList) PopBacks(max int) (values []VType) {
	return l.pops(max, false)
}

// PopFronts 队头弹出至多max个元素
func (l *SyncList) PopFronts(max int) (values []VType) {
	return l.pops(max, true)
}

// PopBackAll 队尾弹出所有元素
func (l *SyncList) PopBackAll() []VType {
	return l.PopBacks(-1)
}

// PopFrontAll 队头弹出所有元素
func (l *SyncList) PopFrontAll() []VType {
	return l.PopFronts(-1)
}

// FrontAll 队头获取所有元素，拷贝操作
func (l *SyncList) FrontAll() (values []VType) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]VType, length)
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value.(VType)
		}
	}
	return
}

// BackAll 队尾获取所有元素，拷贝操作
func (l *SyncList) BackAll() (values []VType) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]VType, length)
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value.(VType)
		}
	}
	return
}

// FrontValue 获取队头元素
func (l *SyncList) FrontValue() (value VType) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Front(); e != nil {
		value = e.Value.(VType)
	}
	return
}

// BackValue 获取队尾元素
func (l *SyncList) BackValue() (value VType) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Back(); e != nil {
		value = e.Value.(VType)
	}
	return
}

// Front returns the first element of list l or nil if the list is empty.
func (l *SyncList) Front() (e *Element) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Front()
	return
}

// Back returns the last element of list l or nil if the list is empty.
func (l *SyncList) Back() (e *Element) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Back()
	return
}

// Len 获取长度，空返回0
func (l *SyncList) Len() (length int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length = l.list.Len()
	return
}

// Size Len的alias方法
func (l *SyncList) Size() int {
	return l.Len()
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *SyncList) MoveBefore(e, mark *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveBefore(e, mark)
}

// MoveAfter moves element <e> to its new position after <p>.
// If <e> or <p> is not an element of <l>, or <e> == <p>, the list is not modified.
// The element and <p> must not be nil.
func (l *SyncList) MoveAfter(e, p *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveAfter(e, p)
}

// MoveToFront moves element <e> to the front of list <l>.
// If <e> is not an element of <l>, the list is not modified.
// The element must not be nil.
func (l *SyncList) MoveToFront(e *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveToFront(e)
}

// MoveToBack moves element <e> to the back of list <l>.
// If <e> is not an element of <l>, the list is not modified.
// The element must not be nil.
func (l *SyncList) MoveToBack(e *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveToBack(e)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *SyncList) PushBackList(other *SyncList) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.PushBackList(other.list)
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *SyncList) PushFrontList(other *SyncList) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.PushFrontList(other.list)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *SyncList) InsertAfter(p *Element, v VType) (e *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.InsertAfter(v, p)
	return
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *SyncList) InsertBefore(p *Element, v VType) (e *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.InsertBefore(v, p)
	return
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *SyncList) Remove(e *Element) (value VType) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	value = l.list.Remove(e).(VType)
	return
}

// Removes 删除多个元素，底层调用Remove
func (l *SyncList) Removes(es []*Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	for _, e := range es {
		l.list.Remove(e)
	}
}

// RemoveAll 删除所有元素
func (l *SyncList) RemoveAll() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// Clear See RemoveAll().
func (l *SyncList) Clear() {
	l.RemoveAll()
}

// RLockFunc 读操作调用f方法
func (l *SyncList) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list != nil {
		f(l.list)
	}
}

// LockFunc 写操作调用f方法
func (l *SyncList) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	f(l.list)
}

// Iterator is alias of IteratorAsc.
func (l *SyncList) Iterator(f func(e *Element) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc 正序遍历，如果f返回false则停止遍历
func (l *SyncList) IteratorAsc(f func(e *Element) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			if !f(e) {
				break
			}
		}
	}
}

// IteratorDesc 逆序遍历，如果f返回false则停止遍历
func (l *SyncList) IteratorDesc(f func(e *Element) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			if !f(e) {
				break
			}
		}
	}
}

type localRWMutexVType struct {
	*sync.RWMutex
}

func newLocalRWMutexVType(safe bool) *localRWMutexVType {
	mu := localRWMutexVType{}
	if safe {
		mu.RWMutex = new(sync.RWMutex)
	}
	return &mu
}

func (mu *localRWMutexVType) IsSafe() bool {
	return mu.RWMutex != nil
}

func (mu *localRWMutexVType) Lock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Lock()
	}
}

func (mu *localRWMutexVType) Unlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Unlock()
	}
}

func (mu *localRWMutexVType) RLock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RLock()
	}
}

func (mu *localRWMutexVType) RUnlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RUnlock()
	}
}

//template format
var __formatTo func(interface{}) VType

func TestSyncList(t *testing.T) {
	Convey("test sync list", t, func() {
		for _, tr := range []*SyncList{New(), NewSync()} {
			So(tr.Len(), ShouldBeZeroValue)
			var e0 = __formatTo(3)
			tr.PushBack(e0)
			So(tr.Len(), ShouldEqual, 1)
			e := tr.PopBack()
			So(e, ShouldEqual, e0)

			tr.PushFront(e0)
			So(tr.Len(), ShouldEqual, 1)
			e = tr.PopFront()
			So(e, ShouldEqual, e0)

			ps := []VType{__formatTo(1), __formatTo(1), __formatTo(3), __formatTo(4), __formatTo(5)}
			tr.PushFronts(ps)
			So(tr.Len(), ShouldEqual, len(ps))
			pops := tr.PopFronts(2)
			So(len(pops), ShouldEqual, 2)
			pops = tr.PopFrontAll()
			So(len(pops), ShouldEqual, 3)

			ps = []VType{__formatTo(1), __formatTo(2), __formatTo(3), __formatTo(4)}
			tr.PushBacks(ps)
			So(tr.Len(), ShouldEqual, len(ps))
			pops = tr.PopBacks(2)
			So(len(pops), ShouldEqual, 2)
			pops = tr.PopBackAll()
			So(len(pops), ShouldEqual, 2)

			ps = []VType{__formatTo(1), __formatTo(2), __formatTo(3), __formatTo(4)}
			tr.Clear()
			tr.PushBacks(ps)
			So(tr.FrontAll(), ShouldResemble, ps)
			tr.Clear()

			psrev := []VType{__formatTo(4), __formatTo(3), __formatTo(2), __formatTo(1)}
			tr.PushBacks(ps)
			So(tr.BackAll(), ShouldResemble, psrev)

			So(tr.FrontValue(), ShouldEqual, __formatTo(1))
			So(tr.Front().Value, ShouldEqual, __formatTo(1))
			So(tr.BackValue(), ShouldEqual, __formatTo(4))
			So(tr.Back().Value, ShouldEqual, __formatTo(4))

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

			n, ns, ol := New(), NewSync(), tr.Len()
			n.PushFronts([]VType{__formatTo(1), __formatTo(2)})
			ns.PushFronts([]VType{__formatTo(1), __formatTo(2)})
			tr.PushFrontList(n)
			So(tr.Len(), ShouldEqual, ol+2)
			tr.PushFrontList(ns)
			So(tr.Len(), ShouldEqual, ol+2+2)

			f0, trl := tr.Front(), tr.Len()
			tr.InsertBefore(tr.Front(), __formatTo(10))
			So(tr.Front().Next(), ShouldEqual, f0)
			So(tr.Front().Value, ShouldEqual, __formatTo(10))
			So(tr.Len(), ShouldEqual, trl+1)

			b, trl = tr.Back(), tr.Len()
			tr.InsertAfter(tr.Back(), __formatTo(10))
			So(tr.Back().Prev(), ShouldEqual, b)
			So(tr.Back().Value, ShouldEqual, __formatTo(10))
			So(tr.Len(), ShouldEqual, trl+1)

			bv := tr.Back().Value
			So(tr.Remove(tr.Back()), ShouldEqual, bv)

			So(func() { tr.Removes([]*Element{tr.Front(), tr.Front().Next()}) }, ShouldNotPanic)
			So(func() { tr.RemoveAll() }, ShouldNotPanic)
			So(tr.Len(), ShouldEqual, 0)

			tr.PushFrontList(n)
			tr.Clear()
			So(tr.Len(), ShouldEqual, 0)

			tr.PushFronts([]VType{__formatTo(10), __formatTo(20), __formatTo(30), __formatTo(40)})

			So(func() {
				tr.RLockFunc(func(list *list.List) {
					So(list.Front().Value, ShouldEqual, __formatTo(40))
				})
			}, ShouldNotPanic)

			So(func() {
				tr.LockFunc(func(list *list.List) {
					So(list.Front().Value, ShouldEqual, __formatTo(40))
				})
			}, ShouldNotPanic)

			So(func() {
				tr.Iterator(func(e *Element) bool {
					return true
				})
			}, ShouldNotPanic)

			So(func() {
				tr.IteratorAsc(func(e *Element) bool {
					return true
				})
			}, ShouldNotPanic)

			So(func() {
				tr.IteratorDesc(func(e *Element) bool {
					return true
				})
			}, ShouldNotPanic)

		}
	})
}
