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

type ElementInt8 = list.Element

type Int8 struct {
	mu   *localRWMutexVTypeInt8
	list *list.List
}

func newWithSafeInt8(safe bool) *Int8 {
	return &Int8{
		mu:   newLocalRWMutexVTypeInt8(safe),
		list: list.New(),
	}
}

// New 创建非协程安全版本
func NewInt8() *Int8 { return newWithSafeInt8(false) }

// NewSync 创建协程安全版本
func NewSyncInt8() *Int8 { return newWithSafeInt8(true) }

// PushFront 队头添加
func (l *Int8) PushFront(v int8) (e *ElementInt8) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushFront(v)
	l.mu.Unlock()
	return
}

// PushBack 队尾添加
func (l *Int8) PushBack(v int8) (e *ElementInt8) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushBack(v)
	l.mu.Unlock()
	return
}

// PushFronts 队头添加多个元素
func (l *Int8) PushFronts(values []int8) {
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
func (l *Int8) PushBacks(values []int8) {
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
func (l *Int8) PopBack() (value int8) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Back(); e != nil {
		value = (l.list.Remove(e)).(int8)
	}
	return
}

// PopFront 队头弹出元素
func (l *Int8) PopFront() (value int8) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Front(); e != nil {
		value = l.list.Remove(e).(int8)
	}
	return
}

func (l *Int8) pops(max int, front bool) (values []int8) {
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
		values = make([]int8, length)
		for i := 0; i < length; i++ {
			if front {
				values[i] = l.list.Remove(l.list.Front()).(int8)
			} else {
				values[i] = l.list.Remove(l.list.Back()).(int8)
			}
		}
	}
	return
}

// PopBacks 队尾弹出至多max个元素
func (l *Int8) PopBacks(max int) (values []int8) {
	return l.pops(max, false)
}

// PopFronts 队头弹出至多max个元素
func (l *Int8) PopFronts(max int) (values []int8) {
	return l.pops(max, true)
}

// PopBackAll 队尾弹出所有元素
func (l *Int8) PopBackAll() []int8 {
	return l.PopBacks(-1)
}

// PopFrontAll 队头弹出所有元素
func (l *Int8) PopFrontAll() []int8 {
	return l.PopFronts(-1)
}

// FrontAll 队头获取所有元素，拷贝操作
func (l *Int8) FrontAll() (values []int8) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]int8, length)
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value.(int8)
		}
	}
	return
}

// BackAll 队尾获取所有元素，拷贝操作
func (l *Int8) BackAll() (values []int8) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]int8, length)
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value.(int8)
		}
	}
	return
}

// FrontValue 获取队头元素
func (l *Int8) FrontValue() (value int8) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Front(); e != nil {
		value = e.Value.(int8)
	}
	return
}

// BackValue 获取队尾元素
func (l *Int8) BackValue() (value int8) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Back(); e != nil {
		value = e.Value.(int8)
	}
	return
}

// Front returns the first element of list l or nil if the list is empty.
func (l *Int8) Front() (e *ElementInt8) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Front()
	return
}

// Back returns the last element of list l or nil if the list is empty.
func (l *Int8) Back() (e *ElementInt8) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Back()
	return
}

// Len 获取长度，空返回0
func (l *Int8) Len() (length int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length = l.list.Len()
	return
}

// Size Len的alias方法
func (l *Int8) Size() int {
	return l.Len()
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *Int8) MoveBefore(e, mark *ElementInt8) {
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
func (l *Int8) MoveAfter(e, p *ElementInt8) {
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
func (l *Int8) MoveToFront(e *ElementInt8) {
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
func (l *Int8) MoveToBack(e *ElementInt8) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveToBack(e)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *Int8) PushBackList(other *Int8) {
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
func (l *Int8) PushFrontList(other *Int8) {
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
func (l *Int8) InsertAfter(p *ElementInt8, v int8) (e *ElementInt8) {
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
func (l *Int8) InsertBefore(p *ElementInt8, v int8) (e *ElementInt8) {
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
func (l *Int8) Remove(e *ElementInt8) (value int8) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	value = l.list.Remove(e).(int8)
	return
}

// Removes 删除多个元素，底层调用Remove
func (l *Int8) Removes(es []*ElementInt8) {
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
func (l *Int8) RemoveAll() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// Clear See RemoveAll().
func (l *Int8) Clear() {
	l.RemoveAll()
}

// RLockFunc 读操作调用f方法
func (l *Int8) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list != nil {
		f(l.list)
	}
}

// LockFunc 写操作调用f方法
func (l *Int8) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	f(l.list)
}

// Iterator is alias of IteratorAsc.
func (l *Int8) Iterator(f func(e *ElementInt8) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc 正序遍历，如果f返回false则停止遍历
func (l *Int8) IteratorAsc(f func(e *ElementInt8) bool) {
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
func (l *Int8) IteratorDesc(f func(e *ElementInt8) bool) {
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

type localRWMutexVTypeInt8 struct {
	*sync.RWMutex
}

func newLocalRWMutexVTypeInt8(safe bool) *localRWMutexVTypeInt8 {
	mu := localRWMutexVTypeInt8{}
	if safe {
		mu.RWMutex = new(sync.RWMutex)
	}
	return &mu
}

func (mu *localRWMutexVTypeInt8) IsSafe() bool {
	return mu.RWMutex != nil
}

func (mu *localRWMutexVTypeInt8) Lock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Lock()
	}
}

func (mu *localRWMutexVTypeInt8) Unlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Unlock()
	}
}

func (mu *localRWMutexVTypeInt8) RLock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RLock()
	}
}

func (mu *localRWMutexVTypeInt8) RUnlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RUnlock()
	}
}

//template format
var __formatToInt8 = func(i interface{}) int8 {
	switch ii := i.(type) {
	case int:
		return int8(ii)
	case int8:
		return int8(ii)
	case int16:
		return int8(ii)
	case int32:
		return int8(ii)
	case int64:
		return int8(ii)
	case uint:
		return int8(ii)
	case uint8:
		return int8(ii)
	case uint16:
		return int8(ii)
	case uint32:
		return int8(ii)
	case uint64:
		return int8(ii)
	case float32:
		return int8(ii)
	case float64:
		return int8(ii)
	case string:
		iv, err := strconv.ParseInt(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return int8(iv)
	default:
		panic("unknown type")
	}
}
