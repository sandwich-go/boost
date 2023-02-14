// Code generated by gotemplate. DO NOT EDIT.

package sarray

import (
	"container/list"
	"strconv"

	"sync"
)

//template type SyncList(VType)

type ElementUint = list.Element

type Uint struct {
	mu   *localRWMutexVTypeUint
	list *list.List
}

func newWithSafeUint(safe bool) *Uint {
	return &Uint{
		mu:   newLocalRWMutexVTypeUint(safe),
		list: list.New(),
	}
}
func NewUint() *Uint     { return newWithSafeUint(false) }
func NewSyncUint() *Uint { return newWithSafeUint(true) }

// PushFront 队头添加
func (l *Uint) PushFront(v uint) (e *ElementUint) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushFront(v)
	l.mu.Unlock()
	return
}

// PushBack 队尾添加
func (l *Uint) PushBack(v uint) (e *ElementUint) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushBack(v)
	l.mu.Unlock()
	return
}

// PushFronts 队头添加多个元素
func (l *Uint) PushFronts(values []uint) {
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
func (l *Uint) PushBacks(values []uint) {
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
func (l *Uint) PopBack() (value uint) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Back(); e != nil {
		value = (l.list.Remove(e)).(uint)
	}
	return
}

// PopFront 队头弹出元素
func (l *Uint) PopFront() (value uint) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Front(); e != nil {
		value = l.list.Remove(e).(uint)
	}
	return
}

func (l *Uint) pops(max int, front bool) (values []uint) {
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
		values = make([]uint, length)
		for i := 0; i < length; i++ {
			if front {
				values[i] = l.list.Remove(l.list.Front()).(uint)
			} else {
				values[i] = l.list.Remove(l.list.Back()).(uint)
			}
		}
	}
	return
}

// PopBacks 队尾弹出至多max个元素
func (l *Uint) PopBacks(max int) (values []uint) {
	return l.pops(max, false)
}

// PopFronts 队头弹出至多max个元素
func (l *Uint) PopFronts(max int) (values []uint) {
	return l.pops(max, true)
}

// PopBackAll 队尾弹出所有元素
func (l *Uint) PopBackAll() []uint {
	return l.PopBacks(-1)
}

// PopFrontAll 队头弹出所有元素
func (l *Uint) PopFrontAll() []uint {
	return l.PopFronts(-1)
}

// FrontAll 队头获取所有元素，拷贝操作
func (l *Uint) FrontAll() (values []uint) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]uint, length)
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value.(uint)
		}
	}
	return
}

// BackAll 队尾获取所有元素，拷贝操作
func (l *Uint) BackAll() (values []uint) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]uint, length)
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value.(uint)
		}
	}
	return
}

// FrontValue 获取队头元素
func (l *Uint) FrontValue() (value uint) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Front(); e != nil {
		value = e.Value.(uint)
	}
	return
}

// BackValue 获取队尾元素
func (l *Uint) BackValue() (value uint) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Back(); e != nil {
		value = e.Value.(uint)
	}
	return
}

// Front returns the first element of list l or nil if the list is empty.
func (l *Uint) Front() (e *ElementUint) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Front()
	return
}

// Back returns the last element of list l or nil if the list is empty.
func (l *Uint) Back() (e *ElementUint) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Back()
	return
}

// Len 获取长度，空返回0
func (l *Uint) Len() (length int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length = l.list.Len()
	return
}

// Size Len的alias方法
func (l *Uint) Size() int {
	return l.Len()
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *Uint) MoveBefore(e, mark *ElementUint) {
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
func (l *Uint) MoveAfter(e, p *ElementUint) {
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
func (l *Uint) MoveToFront(e *ElementUint) {
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
func (l *Uint) MoveToBack(e *ElementUint) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveToBack(e)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *Uint) PushBackList(other *Uint) {
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
func (l *Uint) PushFrontList(other *Uint) {
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
func (l *Uint) InsertAfter(p *ElementUint, v uint) (e *ElementUint) {
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
func (l *Uint) InsertBefore(p *ElementUint, v uint) (e *ElementUint) {
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
func (l *Uint) Remove(e *ElementUint) (value uint) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	value = l.list.Remove(e).(uint)
	return
}

// Removes 删除多个元素，底层调用Remove
func (l *Uint) Removes(es []*ElementUint) {
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
func (l *Uint) RemoveAll() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// Clear See RemoveAll().
func (l *Uint) Clear() {
	l.RemoveAll()
}

// RLockFunc 读操作调用f方法
func (l *Uint) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list != nil {
		f(l.list)
	}
}

// LockFunc 写操作调用f方法
func (l *Uint) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	f(l.list)
}

// Iterator is alias of IteratorAsc.
func (l *Uint) Iterator(f func(e *ElementUint) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc 正序遍历，如果f返回false则停止遍历
func (l *Uint) IteratorAsc(f func(e *ElementUint) bool) {
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
func (l *Uint) IteratorDesc(f func(e *ElementUint) bool) {
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

type localRWMutexVTypeUint struct {
	*sync.RWMutex
}

func newLocalRWMutexVTypeUint(safe bool) *localRWMutexVTypeUint {
	mu := localRWMutexVTypeUint{}
	if safe {
		mu.RWMutex = new(sync.RWMutex)
	}
	return &mu
}

func (mu *localRWMutexVTypeUint) IsSafe() bool {
	return mu.RWMutex != nil
}

func (mu *localRWMutexVTypeUint) Lock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Lock()
	}
}

func (mu *localRWMutexVTypeUint) Unlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Unlock()
	}
}

func (mu *localRWMutexVTypeUint) RLock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RLock()
	}
}

func (mu *localRWMutexVTypeUint) RUnlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RUnlock()
	}
}

//template format
var __formatToUint = func(i interface{}) uint {
	switch ii := i.(type) {
	case int:
		return uint(ii)
	case int8:
		return uint(ii)
	case int16:
		return uint(ii)
	case int32:
		return uint(ii)
	case int64:
		return uint(ii)
	case uint:
		return uint(ii)
	case uint8:
		return uint(ii)
	case uint16:
		return uint(ii)
	case uint32:
		return uint(ii)
	case uint64:
		return uint(ii)
	case float32:
		return uint(ii)
	case float64:
		return uint(ii)
	case string:
		iv, err := strconv.ParseUint(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return uint(iv)
	default:
		panic("unknown type")
	}
}
