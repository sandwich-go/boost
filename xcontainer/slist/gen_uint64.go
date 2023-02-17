// Code generated by gotemplate. DO NOT EDIT.

// slist 包提供了一个同步的链表实现
// 可以产生一个带读写锁的线程安全的SyncList，也可以产生一个非线程安全的SyncList
// New 产生非协程安全的版本
// NewSync 产生协程安全的版本
package slist

import (
	"container/list"
	"strconv"

	"sync"
)

//template type SyncList(VType)

type ElementUint64 = list.Element

type Uint64 struct {
	mu   *localRWMutexVTypeUint64
	list *list.List
}

func newWithSafeUint64(safe bool) *Uint64 {
	return &Uint64{
		mu:   newLocalRWMutexVTypeUint64(safe),
		list: list.New(),
	}
}

// New 创建非协程安全版本
func NewUint64() *Uint64 { return newWithSafeUint64(false) }

// NewSync 创建协程安全版本
func NewSyncUint64() *Uint64 { return newWithSafeUint64(true) }

// PushFront 队头添加
func (l *Uint64) PushFront(v uint64) (e *ElementUint64) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushFront(v)
	l.mu.Unlock()
	return
}

// PushBack 队尾添加
func (l *Uint64) PushBack(v uint64) (e *ElementUint64) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushBack(v)
	l.mu.Unlock()
	return
}

// PushFronts 队头添加多个元素
func (l *Uint64) PushFronts(values []uint64) {
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
func (l *Uint64) PushBacks(values []uint64) {
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
func (l *Uint64) PopBack() (value uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Back(); e != nil {
		value = (l.list.Remove(e)).(uint64)
	}
	return
}

// PopFront 队头弹出元素
func (l *Uint64) PopFront() (value uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Front(); e != nil {
		value = l.list.Remove(e).(uint64)
	}
	return
}

func (l *Uint64) pops(max int, front bool) (values []uint64) {
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
		values = make([]uint64, length)
		for i := 0; i < length; i++ {
			if front {
				values[i] = l.list.Remove(l.list.Front()).(uint64)
			} else {
				values[i] = l.list.Remove(l.list.Back()).(uint64)
			}
		}
	}
	return
}

// PopBacks 队尾弹出至多max个元素
func (l *Uint64) PopBacks(max int) (values []uint64) {
	return l.pops(max, false)
}

// PopFronts 队头弹出至多max个元素
func (l *Uint64) PopFronts(max int) (values []uint64) {
	return l.pops(max, true)
}

// PopBackAll 队尾弹出所有元素
func (l *Uint64) PopBackAll() []uint64 {
	return l.PopBacks(-1)
}

// PopFrontAll 队头弹出所有元素
func (l *Uint64) PopFrontAll() []uint64 {
	return l.PopFronts(-1)
}

// FrontAll 队头获取所有元素，拷贝操作
func (l *Uint64) FrontAll() (values []uint64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]uint64, length)
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value.(uint64)
		}
	}
	return
}

// BackAll 队尾获取所有元素，拷贝操作
func (l *Uint64) BackAll() (values []uint64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]uint64, length)
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value.(uint64)
		}
	}
	return
}

// FrontValue 获取队头元素
func (l *Uint64) FrontValue() (value uint64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Front(); e != nil {
		value = e.Value.(uint64)
	}
	return
}

// BackValue 获取队尾元素
func (l *Uint64) BackValue() (value uint64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Back(); e != nil {
		value = e.Value.(uint64)
	}
	return
}

// Front returns the first element of list l or nil if the list is empty.
func (l *Uint64) Front() (e *ElementUint64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Front()
	return
}

// Back returns the last element of list l or nil if the list is empty.
func (l *Uint64) Back() (e *ElementUint64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Back()
	return
}

// Len 获取长度，空返回0
func (l *Uint64) Len() (length int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length = l.list.Len()
	return
}

// Size Len的alias方法
func (l *Uint64) Size() int {
	return l.Len()
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *Uint64) MoveBefore(e, mark *ElementUint64) {
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
func (l *Uint64) MoveAfter(e, p *ElementUint64) {
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
func (l *Uint64) MoveToFront(e *ElementUint64) {
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
func (l *Uint64) MoveToBack(e *ElementUint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveToBack(e)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *Uint64) PushBackList(other *Uint64) {
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
func (l *Uint64) PushFrontList(other *Uint64) {
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
func (l *Uint64) InsertAfter(p *ElementUint64, v uint64) (e *ElementUint64) {
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
func (l *Uint64) InsertBefore(p *ElementUint64, v uint64) (e *ElementUint64) {
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
func (l *Uint64) Remove(e *ElementUint64) (value uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	value = l.list.Remove(e).(uint64)
	return
}

// Removes 删除多个元素，底层调用Remove
func (l *Uint64) Removes(es []*ElementUint64) {
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
func (l *Uint64) RemoveAll() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// Clear See RemoveAll().
func (l *Uint64) Clear() {
	l.RemoveAll()
}

// RLockFunc 读操作调用f方法
func (l *Uint64) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list != nil {
		f(l.list)
	}
}

// LockFunc 写操作调用f方法
func (l *Uint64) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	f(l.list)
}

// Iterator is alias of IteratorAsc.
func (l *Uint64) Iterator(f func(e *ElementUint64) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc 正序遍历，如果f返回false则停止遍历
func (l *Uint64) IteratorAsc(f func(e *ElementUint64) bool) {
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
func (l *Uint64) IteratorDesc(f func(e *ElementUint64) bool) {
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

type localRWMutexVTypeUint64 struct {
	*sync.RWMutex
}

func newLocalRWMutexVTypeUint64(safe bool) *localRWMutexVTypeUint64 {
	mu := localRWMutexVTypeUint64{}
	if safe {
		mu.RWMutex = new(sync.RWMutex)
	}
	return &mu
}

func (mu *localRWMutexVTypeUint64) IsSafe() bool {
	return mu.RWMutex != nil
}

func (mu *localRWMutexVTypeUint64) Lock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Lock()
	}
}

func (mu *localRWMutexVTypeUint64) Unlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Unlock()
	}
}

func (mu *localRWMutexVTypeUint64) RLock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RLock()
	}
}

func (mu *localRWMutexVTypeUint64) RUnlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RUnlock()
	}
}

//template format
var __formatToUint64 = func(i interface{}) uint64 {
	switch ii := i.(type) {
	case int:
		return uint64(ii)
	case int8:
		return uint64(ii)
	case int16:
		return uint64(ii)
	case int32:
		return uint64(ii)
	case int64:
		return uint64(ii)
	case uint:
		return uint64(ii)
	case uint8:
		return uint64(ii)
	case uint16:
		return uint64(ii)
	case uint32:
		return uint64(ii)
	case uint64:
		return uint64(ii)
	case float32:
		return uint64(ii)
	case float64:
		return uint64(ii)
	case string:
		iv, err := strconv.ParseUint(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return uint64(iv)
	default:
		panic("unknown type")
	}
}
