package slist

import (
	"container/list"
	"sync"
)

//template type SyncList(VType)

type VType interface{}
type Element = list.Element

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
func New() *SyncList     { return newWithSafe(false) }
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
