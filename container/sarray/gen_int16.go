// Code generated by gotemplate. DO NOT EDIT.

package sarray

import (
	"encoding/json"
	"fmt"
	"strconv"

	"math/rand"
	"sort"
	"sync"
)

//template type SyncArray(VType)

type Int16 struct {
	mu    *localRWMutexVTypeInt16
	array []int16
}

// New 创建非协程安全版本
func NewInt16() *Int16 { return newWithSafeInt16(false) }

// NewSync 创建协程安全版本
func NewSyncInt16() *Int16 { return newWithSafeInt16(true) }

func newWithSafeInt16(safe bool) *Int16 {
	return &Int16{
		mu:    newLocalRWMutexVTypeInt16(safe),
		array: make([]int16, 0),
	}
}

// At 返回指定位置元素，如果越界则返回默认空值
func (a *Int16) At(index int) (value int16) {
	value, _ = a.Get(index)
	return
}

// Get 返回指定位置元素，found标识元素是否存在
func (a *Int16) Get(index int) (value int16, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.array) {
		return
	}
	return a.array[index], true
}

func (a *Int16) errorIndexOutRangeUnLock(index int) error {
	return fmt.Errorf("index %d out of array range %d", index, len(a.array))
}

// Set 设定指定位置数据
func (a *Int16) Set(index int, value int16) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return a.errorIndexOutRangeUnLock(index)
	}
	a.array[index] = value
	return nil
}

// SetArray 替换底层存储
func (a *Int16) SetArray(array []int16) *Int16 {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = array
	return a
}

// Replace 替换指定位置元素
func (a *Int16) Replace(given []int16) *Int16 {
	a.mu.Lock()
	defer a.mu.Unlock()
	max := len(given)
	if max > len(a.array) {
		max = len(a.array)
	}
	for i := 0; i < max; i++ {
		a.array[i] = given[i]
	}
	return a
}

// SortFunc  根据指定的方法进行排序
func (a *Int16) SortFunc(less func(v1, v2 int16) bool) *Int16 {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.array, func(i, j int) bool {
		return less(a.array[i], a.array[j])
	})
	return a
}

// InsertBefore 在index位置前插入数据
func (a *Int16) InsertBefore(index int, value int16) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return a.errorIndexOutRangeUnLock(index)
	}
	rear := append([]int16{}, a.array[index:]...)
	a.array = append(a.array[0:index], value)
	a.array = append(a.array, rear...)
	return nil
}

// InsertAfter 在index位置后插入数据
func (a *Int16) InsertAfter(index int, value int16) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return a.errorIndexOutRangeUnLock(index)
	}
	rear := append([]int16{}, a.array[index+1:]...)
	a.array = append(a.array[0:index+1], value)
	a.array = append(a.array, rear...)
	return nil
}

// Contains  是否存在value
func (a *Int16) Contains(value int16) bool {
	return a.Search(value) != -1
}

// Search 查找元素，不存在返回-1
func (a *Int16) Search(value int16) int {
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

func (a *Int16) DeleteValue(value int16) (found bool) {
	if i := a.Search(value); i != -1 {
		_, found = a.LoadAndDelete(i)
		return found
	}
	return false
}

// LoadAndDelete 删除元素，如果删除成功返回被删除的元素
func (a *Int16) LoadAndDelete(index int) (value int16, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doDeleteWithoutLock(index)
}

// doRemoveWithoutLock 不加锁移除元素
func (a *Int16) doDeleteWithoutLock(index int) (value int16, found bool) {
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
func (a *Int16) PushLeft(value ...int16) *Int16 {
	a.mu.Lock()
	a.array = append(value, a.array...)
	a.mu.Unlock()
	return a
}

// PushRight 尾插入
func (a *Int16) PushRight(value ...int16) *Int16 {
	a.mu.Lock()
	a.array = append(a.array, value...)
	a.mu.Unlock()
	return a
}

// PopLeft 头弹出
func (a *Int16) PopLeft() (value int16, found bool) {
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
func (a *Int16) PopRight() (value int16, found bool) {
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
func (a *Int16) PopRand() (value int16, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doDeleteWithoutLock(rand.Intn(len(a.array)))
}

// PopRands 随机n个元素并弹出，如果size大于数组尺寸则全部弹出
func (a *Int16) PopRands(size int) []int16 {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	if size >= len(a.array) {
		size = len(a.array)
	}
	array := make([]int16, size)
	for i := 0; i < size; i++ {
		array[i], _ = a.doDeleteWithoutLock(rand.Intn(len(a.array)))
	}
	return array
}

// Append 尾添加元素 alias of PushRight
func (a *Int16) Append(value ...int16) *Int16 { return a.PushRight(value...) }

// Len 获取长度
func (a *Int16) Len() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// Slice 获取底层数据存储，如果为sync安全模式则返回一份拷贝，否则直接返回底层数据指针
func (a *Int16) Slice() []int16 {
	array := ([]int16)(nil)
	if a.mu.IsSafe() {
		a.mu.RLock()
		defer a.mu.RUnlock()
		array = make([]int16, len(a.array))
		copy(array, a.array)
	} else {
		array = a.array
	}
	return array
}

// Clear 清空存储
func (a *Int16) Clear() *Int16 {
	a.mu.Lock()
	if len(a.array) > 0 {
		a.array = make([]int16, 0)
	}
	a.mu.Unlock()
	return a
}

// LockFunc 写锁操作array
func (a *Int16) LockFunc(f func(array []int16)) *Int16 {
	a.mu.Lock()
	defer a.mu.Unlock()
	f(a.array)
	return a
}

// RLockFunc 读锁操作array
func (a *Int16) RLockFunc(f func(array []int16)) *Int16 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	f(a.array)
	return a
}

// Rand 随机一个元素
func (a *Int16) Rand() (value int16, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return
	}
	return a.array[rand.Intn(len(a.array))], true
}

// WalkAsc 按照index从小到大的顺序进行遍历，并将k,v作为参数执行f。如果f执行返回false则中止
func (a *Int16) WalkAsc(f func(k int, v int16) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for k, v := range a.array {
		if !f(k, v) {
			break
		}
	}
}

// WalkDesc 按照index从大到小的顺序进行遍历，并将k,v作为参数执行f。如果f执行返回false则中止
func (a *Int16) WalkDesc(f func(k int, v int16) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for i := len(a.array) - 1; i >= 0; i-- {
		if !f(i, a.array[i]) {
			break
		}
	}
}

// MarshalJSON 序列化到json
func (a Int16) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a.array)
}

// UnmarshalJSON 由json反序列化
func (a *Int16) UnmarshalJSON(b []byte) error {
	if a.array == nil {
		a.array = make([]int16, 0)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := json.Unmarshal(b, &a.array); err != nil {
		return err
	}
	return nil
}

// Empty 是否为空
func (a *Int16) Empty() bool { return a.Len() == 0 }

type localRWMutexVTypeInt16 struct {
	*sync.RWMutex
}

func newLocalRWMutexVTypeInt16(safe bool) *localRWMutexVTypeInt16 {
	mu := localRWMutexVTypeInt16{}
	if safe {
		mu.RWMutex = new(sync.RWMutex)
	}
	return &mu
}

func (mu *localRWMutexVTypeInt16) IsSafe() bool {
	return mu.RWMutex != nil
}

func (mu *localRWMutexVTypeInt16) Lock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Lock()
	}
}

func (mu *localRWMutexVTypeInt16) Unlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Unlock()
	}
}

func (mu *localRWMutexVTypeInt16) RLock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RLock()
	}
}

func (mu *localRWMutexVTypeInt16) RUnlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RUnlock()
	}
}

//template format
var __formatToInt16 = func(i interface{}) int16 {
	switch ii := i.(type) {
	case int:
		return int16(ii)
	case int8:
		return int16(ii)
	case int16:
		return int16(ii)
	case int32:
		return int16(ii)
	case int64:
		return int16(ii)
	case uint:
		return int16(ii)
	case uint8:
		return int16(ii)
	case uint16:
		return int16(ii)
	case uint32:
		return int16(ii)
	case uint64:
		return int16(ii)
	case float32:
		return int16(ii)
	case float64:
		return int16(ii)
	case string:
		iv, err := strconv.ParseInt(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return int16(iv)
	default:
		panic("unknown type")
	}
}
