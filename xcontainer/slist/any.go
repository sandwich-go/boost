// slist 包提供了一个同步的链表实现
// 可以产生一个带读写锁的线程安全的SyncList，也可以产生一个非线程安全的SyncList
// New 产生非协程安全的版本
// NewSync 产生协程安全的版本
package slist

import (
	"container/list"

	"sync"
)

type List[T any] struct {
	mu   *localRWMutex
	list *list.List
}

func newWithSafe[T any](safe bool) *List[T] {
	return &List[T]{
		mu:   newLocalRWMutex(safe),
		list: list.New(),
	}
}

// New 创建非协程安全版本
func New[T any]() *List[T] { return newWithSafe[T](false) }

// NewSync 创建协程安全版本
func NewSync[T any]() *List[T] { return newWithSafe[T](true) }

// PushFront 队头添加
func (l *List[T]) PushFront(v T) (e *list.Element) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushFront(v)
	l.mu.Unlock()
	return
}

// PushBack 队尾添加
func (l *List[T]) PushBack(v T) (e *list.Element) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushBack(v)
	l.mu.Unlock()
	return
}

// PushFronts 队头添加多个元素
func (l *List[T]) PushFronts(values []T) {
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
func (l *List[T]) PushBacks(values []T) {
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
func (l *List[T]) PopBack() (value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Back(); e != nil {
		value = l.list.Remove(e).(T)
	}
	return
}

// PopFront 队头弹出元素
func (l *List[T]) PopFront() (value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Front(); e != nil {
		value = l.list.Remove(e).(T)
	}
	return
}

func (l *List[T]) pops(max int, front bool) (values []T) {
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
		values = make([]T, length)
		for i := 0; i < length; i++ {
			if front {
				values[i] = l.list.Remove(l.list.Front()).(T)
			} else {
				values[i] = l.list.Remove(l.list.Back()).(T)
			}
		}
	}
	return
}

// PopBacks 队尾弹出至多max个元素
func (l *List[T]) PopBacks(max int) (values []T) {
	return l.pops(max, false)
}

// PopFronts 队头弹出至多max个元素
func (l *List[T]) PopFronts(max int) (values []T) {
	return l.pops(max, true)
}

// PopBackAll 队尾弹出所有元素
func (l *List[T]) PopBackAll() []T {
	return l.PopBacks(-1)
}

// PopFrontAll 队头弹出所有元素
func (l *List[T]) PopFrontAll() []T {
	return l.PopFronts(-1)
}

// FrontAll 队头获取所有元素，拷贝操作
func (l *List[T]) FrontAll() (values []T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]T, length)
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value.(T)
		}
	}
	return
}

// BackAll 队尾获取所有元素，拷贝操作
func (l *List[T]) BackAll() (values []T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]T, length)
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value.(T)
		}
	}
	return
}

// FrontValue 获取队头元素
func (l *List[T]) FrontValue() (value T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Front(); e != nil {
		value = e.Value.(T)
	}
	return
}

// BackValue 获取队尾元素
func (l *List[T]) BackValue() (value T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Back(); e != nil {
		value = e.Value.(T)
	}
	return
}

// Front returns the first element of list l or nil if the list is empty.
func (l *List[T]) Front() (e *list.Element) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Front()
	return
}

// Back returns the last element of list l or nil if the list is empty.
func (l *List[T]) Back() (e *list.Element) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Back()
	return
}

// Len 获取长度，空返回0
func (l *List[T]) Len() (length int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length = l.list.Len()
	return
}

// Size Len的alias方法
func (l *List[T]) Size() int {
	return l.Len()
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *List[T]) MoveBefore(e, mark *list.Element) {
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
func (l *List[T]) MoveAfter(e, p *list.Element) {
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
func (l *List[T]) MoveToFront(e *list.Element) {
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
func (l *List[T]) MoveToBack(e *list.Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveToBack(e)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List[T]) PushBackList(other *List[T]) {
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
func (l *List[T]) PushFrontList(other *List[T]) {
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
func (l *List[T]) InsertAfter(p *list.Element, v T) (e *list.Element) {
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
func (l *List[T]) InsertBefore(p *list.Element, v T) (e *list.Element) {
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
func (l *List[T]) Remove(e *list.Element) (value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	value = l.list.Remove(e).(T)
	return
}

// Removes 删除多个元素，底层调用Remove
func (l *List[T]) Removes(es []*list.Element) {
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
func (l *List[T]) RemoveAll() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// Clear See RemoveAll().
func (l *List[T]) Clear() {
	l.RemoveAll()
}

// RLockFunc 读操作调用f方法
func (l *List[T]) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list != nil {
		f(l.list)
	}
}

// LockFunc 写操作调用f方法
func (l *List[T]) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	f(l.list)
}

// Iterator is alias of IteratorAsc.
func (l *List[T]) Iterator(f func(e *list.Element) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc 正序遍历，如果f返回false则停止遍历
func (l *List[T]) IteratorAsc(f func(e *list.Element) bool) {
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
func (l *List[T]) IteratorDesc(f func(e *list.Element) bool) {
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

type localRWMutex struct {
	*sync.RWMutex
}

func newLocalRWMutex(safe bool) *localRWMutex {
	mu := localRWMutex{}
	if safe {
		mu.RWMutex = new(sync.RWMutex)
	}
	return &mu
}

func (mu *localRWMutex) IsSafe() bool {
	return mu.RWMutex != nil
}

func (mu *localRWMutex) Lock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Lock()
	}
}

func (mu *localRWMutex) Unlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Unlock()
	}
}

func (mu *localRWMutex) RLock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RLock()
	}
}

func (mu *localRWMutex) RUnlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RUnlock()
	}
}
