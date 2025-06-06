package sarray

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

type Array[T comparable] struct {
	mu    *localRWMutex
	array []T
}

// New 创建非协程安全版本
func New[T comparable]() *Array[T] { return newWithSafe[T](false) }

// NewSync 创建协程安全版本
func NewSync[T comparable]() *Array[T] { return newWithSafe[T](true) }

func newWithSafe[T comparable](safe bool) *Array[T] {
	return &Array[T]{
		mu:    newLocalRWMutex(safe),
		array: make([]T, 0),
	}
}

// At 返回指定位置元素，如果越界则返回默认空值
func (a *Array[T]) At(index int) (value T) {
	value, _ = a.Get(index)
	return
}

// Get 返回指定位置元素，found标识元素是否存在
func (a *Array[T]) Get(index int) (value T, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.array) {
		return
	}
	return a.array[index], true
}

func (a *Array[T]) errorIndexOutRangeUnLock(index int) error {
	return fmt.Errorf("index %d out of array range %d", index, len(a.array))
}

// Set 设定指定位置数据
func (a *Array[T]) Set(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return a.errorIndexOutRangeUnLock(index)
	}
	a.array[index] = value
	return nil
}

// SetArray 替换底层存储
func (a *Array[T]) SetArray(array []T) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = array
	return a
}

// Replace 替换指定位置元素
func (a *Array[T]) Replace(given []T) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	l := len(given)
	if l > len(a.array) {
		l = len(a.array)
	}
	for i := 0; i < l; i++ {
		a.array[i] = given[i]
	}
	return a
}

// SortFunc  根据指定的方法进行排序
func (a *Array[T]) SortFunc(less func(v1, v2 T) bool) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.array, func(i, j int) bool {
		return less(a.array[i], a.array[j])
	})
	return a
}

// InsertBefore 在index位置前插入数据
func (a *Array[T]) InsertBefore(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return a.errorIndexOutRangeUnLock(index)
	}
	rear := append([]T{}, a.array[index:]...)
	a.array = append(a.array[0:index], value)
	a.array = append(a.array, rear...)
	return nil
}

// InsertAfter 在index位置后插入数据
func (a *Array[T]) InsertAfter(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return a.errorIndexOutRangeUnLock(index)
	}
	rear := append([]T{}, a.array[index+1:]...)
	a.array = append(a.array[0:index+1], value)
	a.array = append(a.array, rear...)
	return nil
}

// Contains  是否存在value
func (a *Array[T]) Contains(value T) bool {
	return a.Search(value) != -1
}

// Search 查找元素，不存在返回-1
func (a *Array[T]) Search(value T) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return -1
	}
	result := -1
	for index, v := range a.array {
		if v == value {
			result = index
			break
		}
	}
	return result
}

func (a *Array[T]) DeleteValue(value T) (found bool) {
	if i := a.Search(value); i != -1 {
		_, found = a.LoadAndDelete(i)
		return found
	}
	return false
}

// LoadAndDelete 删除元素，如果删除成功返回被删除的元素
func (a *Array[T]) LoadAndDelete(index int) (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doDeleteWithoutLock(index)
}

// doRemoveWithoutLock 不加锁移除元素
func (a *Array[T]) doDeleteWithoutLock(index int) (value T, found bool) {
	if index < 0 || index >= len(a.array) {
		return
	}
	if index == 0 {
		value = a.array[0]
		a.array = a.array[1:]
		return value, true
	} else if index == len(a.array)-1 {
		value = a.array[index]
		a.array = a.array[:index]
		return value, true
	}
	value = a.array[index]
	a.array = append(a.array[:index], a.array[index+1:]...)
	return value, true
}

// PushLeft 头插入
func (a *Array[T]) PushLeft(value ...T) *Array[T] {
	a.mu.Lock()
	a.array = append(value, a.array...)
	a.mu.Unlock()
	return a
}

// PushRight 尾插入
func (a *Array[T]) PushRight(value ...T) *Array[T] {
	a.mu.Lock()
	a.array = append(a.array, value...)
	a.mu.Unlock()
	return a
}

// PopLeft 头弹出
func (a *Array[T]) PopLeft() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		return
	}
	value = a.array[0]
	a.array = a.array[1:]
	return value, true
}

// PopRight 尾弹出
func (a *Array[T]) PopRight() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	index := len(a.array) - 1
	if index < 0 {
		return
	}
	value = a.array[index]
	a.array = a.array[:index]
	return value, true
}

// PopRand 随机弹出
func (a *Array[T]) PopRand() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doDeleteWithoutLock(rand.Intn(len(a.array)))
}

// PopRands 随机n个元素并弹出，如果size大于数组尺寸则全部弹出
func (a *Array[T]) PopRands(size int) []T {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	if size >= len(a.array) {
		size = len(a.array)
	}
	array := make([]T, size)
	for i := 0; i < size; i++ {
		array[i], _ = a.doDeleteWithoutLock(rand.Intn(len(a.array)))
	}
	return array
}

// Append 尾添加元素 alias of PushRight
func (a *Array[T]) Append(value ...T) *Array[T] { return a.PushRight(value...) }

// Len 获取长度
func (a *Array[T]) Len() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// Slice 获取底层数据存储，如果为sync安全模式则返回一份拷贝，否则直接返回底层数据指针
func (a *Array[T]) Slice() []T {
	array := ([]T)(nil)
	if a.mu.IsSafe() {
		a.mu.RLock()
		defer a.mu.RUnlock()
		array = make([]T, len(a.array))
		copy(array, a.array)
	} else {
		array = a.array
	}
	return array
}

// Clear 清空存储
func (a *Array[T]) Clear() *Array[T] {
	a.mu.Lock()
	if len(a.array) > 0 {
		a.array = make([]T, 0)
	}
	a.mu.Unlock()
	return a
}

// LockFunc 写锁操作array
func (a *Array[T]) LockFunc(f func(array []T)) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	f(a.array)
	return a
}

// RLockFunc 读锁操作array
func (a *Array[T]) RLockFunc(f func(array []T)) *Array[T] {
	a.mu.RLock()
	defer a.mu.RUnlock()
	f(a.array)
	return a
}

// Rand 随机一个元素
func (a *Array[T]) Rand() (value T, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return
	}
	return a.array[rand.Intn(len(a.array))], true
}

// WalkAsc 按照index从小到大的顺序进行遍历，并将k,v作为参数执行f。如果f执行返回false则中止
func (a *Array[T]) WalkAsc(f func(k int, v T) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for k, v := range a.array {
		if !f(k, v) {
			break
		}
	}
}

// WalkDesc 按照index从大到小的顺序进行遍历，并将k,v作为参数执行f。如果f执行返回false则中止
func (a *Array[T]) WalkDesc(f func(k int, v T) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for i := len(a.array) - 1; i >= 0; i-- {
		if !f(i, a.array[i]) {
			break
		}
	}
}

// MarshalJSON 序列化到json
func (a *Array[T]) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a.array)
}

// UnmarshalJSON 由json反序列化
func (a *Array[T]) UnmarshalJSON(b []byte) error {
	if a.array == nil {
		a.array = make([]T, 0)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := json.Unmarshal(b, &a.array); err != nil {
		return err
	}
	return nil
}

// Empty 是否为空
func (a *Array[T]) Empty() bool { return a.Len() == 0 }

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
