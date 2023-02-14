// Code generated by gotemplate. DO NOT EDIT.

package sarray

import (
	"container/list"
	"strconv"

	"sync"
)

//template type SyncList(VType)

type ElementInt64 = list.Element

type Int64 struct {
	mu   *localRWMutexVTypeInt64
	list *list.List
}

func newWithSafeInt64(safe bool) *Int64 {
	return &Int64{
		mu:   newLocalRWMutexVTypeInt64(safe),
		list: list.New(),
	}
}
func NewInt64() *Int64     { return newWithSafeInt64(false) }
func NewSyncInt64() *Int64 { return newWithSafeInt64(true) }

// PushFront 队头添加
func (l *Int64) PushFront(v int64) (e *ElementInt64) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushFront(v)
	l.mu.Unlock()
	return
}

// PushBack 队尾添加
func (l *Int64) PushBack(v int64) (e *ElementInt64) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushBack(v)
	l.mu.Unlock()
	return
}

// PushFronts 队头添加多个元素
func (l *Int64) PushFronts(values []int64) {
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
func (l *Int64) PushBacks(values []int64) {
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
func (l *Int64) PopBack() (value int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Back(); e != nil {
		value = (l.list.Remove(e)).(int64)
	}
	return
}

// PopFront 队头弹出元素
func (l *Int64) PopFront() (value int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Front(); e != nil {
		value = l.list.Remove(e).(int64)
	}
	return
}

func (l *Int64) pops(max int, front bool) (values []int64) {
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
		values = make([]int64, length)
		for i := 0; i < length; i++ {
			if front {
				values[i] = l.list.Remove(l.list.Front()).(int64)
			} else {
				values[i] = l.list.Remove(l.list.Back()).(int64)
			}
		}
	}
	return
}

// PopBacks 队尾弹出至多max个元素
func (l *Int64) PopBacks(max int) (values []int64) {
	return l.pops(max, false)
}

// PopFronts 队头弹出至多max个元素
func (l *Int64) PopFronts(max int) (values []int64) {
	return l.pops(max, true)
}

// PopBackAll 队尾弹出所有元素
func (l *Int64) PopBackAll() []int64 {
	return l.PopBacks(-1)
}

// PopFrontAll 队头弹出所有元素
func (l *Int64) PopFrontAll() []int64 {
	return l.PopFronts(-1)
}

// FrontAll 队头获取所有元素，拷贝操作
func (l *Int64) FrontAll() (values []int64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]int64, length)
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value.(int64)
		}
	}
	return
}

// BackAll 队尾获取所有元素，拷贝操作
func (l *Int64) BackAll() (values []int64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]int64, length)
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value.(int64)
		}
	}
	return
}

// FrontValue 获取队头元素
func (l *Int64) FrontValue() (value int64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Front(); e != nil {
		value = e.Value.(int64)
	}
	return
}

// BackValue 获取队尾元素
func (l *Int64) BackValue() (value int64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Back(); e != nil {
		value = e.Value.(int64)
	}
	return
}

// Front returns the first element of list l or nil if the list is empty.
func (l *Int64) Front() (e *ElementInt64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Front()
	return
}

// Back returns the last element of list l or nil if the list is empty.
func (l *Int64) Back() (e *ElementInt64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Back()
	return
}

// Len 获取长度，空返回0
func (l *Int64) Len() (length int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length = l.list.Len()
	return
}

// Size Len的alias方法
func (l *Int64) Size() int {
	return l.Len()
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *Int64) MoveBefore(e, mark *ElementInt64) {
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
func (l *Int64) MoveAfter(e, p *ElementInt64) {
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
func (l *Int64) MoveToFront(e *ElementInt64) {
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
func (l *Int64) MoveToBack(e *ElementInt64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveToBack(e)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *Int64) PushBackList(other *Int64) {
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
func (l *Int64) PushFrontList(other *Int64) {
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
func (l *Int64) InsertAfter(p *ElementInt64, v int64) (e *ElementInt64) {
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
func (l *Int64) InsertBefore(p *ElementInt64, v int64) (e *ElementInt64) {
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
func (l *Int64) Remove(e *ElementInt64) (value int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	value = l.list.Remove(e).(int64)
	return
}

// Removes 删除多个元素，底层调用Remove
func (l *Int64) Removes(es []*ElementInt64) {
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
func (l *Int64) RemoveAll() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// Clear See RemoveAll().
func (l *Int64) Clear() {
	l.RemoveAll()
}

// RLockFunc 读操作调用f方法
func (l *Int64) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list != nil {
		f(l.list)
	}
}

// LockFunc 写操作调用f方法
func (l *Int64) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	f(l.list)
}

// Iterator is alias of IteratorAsc.
func (l *Int64) Iterator(f func(e *ElementInt64) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc 正序遍历，如果f返回false则停止遍历
func (l *Int64) IteratorAsc(f func(e *ElementInt64) bool) {
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
func (l *Int64) IteratorDesc(f func(e *ElementInt64) bool) {
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

type localRWMutexVTypeInt64 struct {
	*sync.RWMutex
}

func newLocalRWMutexVTypeInt64(safe bool) *localRWMutexVTypeInt64 {
	mu := localRWMutexVTypeInt64{}
	if safe {
		mu.RWMutex = new(sync.RWMutex)
	}
	return &mu
}

func (mu *localRWMutexVTypeInt64) IsSafe() bool {
	return mu.RWMutex != nil
}

func (mu *localRWMutexVTypeInt64) Lock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Lock()
	}
}

func (mu *localRWMutexVTypeInt64) Unlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Unlock()
	}
}

func (mu *localRWMutexVTypeInt64) RLock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RLock()
	}
}

func (mu *localRWMutexVTypeInt64) RUnlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RUnlock()
	}
}

//template format
var __formatToInt64 = func(i interface{}) int64 {
	switch ii := i.(type) {
	case int:
		return int64(ii)
	case int8:
		return int64(ii)
	case int16:
		return int64(ii)
	case int32:
		return int64(ii)
	case int64:
		return int64(ii)
	case uint:
		return int64(ii)
	case uint8:
		return int64(ii)
	case uint16:
		return int64(ii)
	case uint32:
		return int64(ii)
	case uint64:
		return int64(ii)
	case float32:
		return int64(ii)
	case float64:
		return int64(ii)
	case string:
		iv, err := strconv.ParseInt(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return int64(iv)
	default:
		panic("unknown type")
	}
}
