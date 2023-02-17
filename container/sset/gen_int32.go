// Code generated by gotemplate. DO NOT EDIT.

// sset 包提供了多种类型的集合
// 可以产生一个带读写锁的线程安全的SyncSet，也可以产生一个非线程安全的SyncSet
// New 产生非协程安全的版本
// NewSync 产生协程安全的版本
package sarray

import (
	"strconv"
	"sync"
)

//template type SyncSet(VType)

type Int32 struct {
	mu   *localRWMutexVTypeInt32
	data map[int32]struct{}
}

// New 创建非协程安全版本
func NewInt32() *Int32 { return newWithSafeInt32(false) }

// NewSync 创建协程安全版本
func NewSyncInt32() *Int32 { return newWithSafeInt32(true) }

func newWithSafeInt32(safe bool) *Int32 {
	return &Int32{data: make(map[int32]struct{}), mu: newLocalRWMutexVTypeInt32(safe)}
}

// Iterator 遍历
func (set *Int32) Iterator(f func(v int32) bool) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k := range set.data {
		if !f(k) {
			break
		}
	}
}

// Add 添加元素
func (set *Int32) Add(items ...int32) {
	set.mu.Lock()
	if set.data == nil {
		set.data = make(map[int32]struct{})
	}
	for _, v := range items {
		set.data[v] = struct{}{}
	}
	set.mu.Unlock()
}

// AddIfNotExist 如果元素不存在则添加，如添加成功则返回true
func (set *Int32) AddIfNotExist(item int32) (addOK bool) {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[int32]struct{})
		}
		if _, ok := set.data[item]; !ok {
			set.data[item] = struct{}{}
			return true
		}
	}
	return false
}

// AddIfNotExistFunc 如果元素不存在且f返回true则添加，如添加成功则返回true
// f函数运行在lock之外
func (set *Int32) AddIfNotExistFunc(item int32, f func() bool) bool {
	if !set.Contains(item) {
		if f() {
			set.mu.Lock()
			defer set.mu.Unlock()
			if set.data == nil {
				set.data = make(map[int32]struct{})
			}
			if _, ok := set.data[item]; !ok {
				set.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// AddIfNotExistFuncLock 如果元素不存在且f返回true则添加，如添加成功则返回true
// f函数运行在lock之内
func (set *Int32) AddIfNotExistFuncLock(item int32, f func() bool) bool {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[int32]struct{})
		}
		if f() {
			if _, ok := set.data[item]; !ok {
				set.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// Contains 是否存在元素
func (set *Int32) Contains(item int32) bool {
	var ok bool
	set.mu.RLock()
	if set.data != nil {
		_, ok = set.data[item]
	}
	set.mu.RUnlock()
	return ok
}

// Remove 移除指定元素
func (set *Int32) Remove(item int32) {
	set.mu.Lock()
	if set.data != nil {
		delete(set.data, item)
	}
	set.mu.Unlock()
}

// Size 返回长度
func (set *Int32) Size() int {
	set.mu.RLock()
	l := len(set.data)
	set.mu.RUnlock()
	return l
}

// Clear 清理元素
func (set *Int32) Clear() {
	set.mu.Lock()
	set.data = make(map[int32]struct{})
	set.mu.Unlock()
}

// Slice 返回元素slice
func (set *Int32) Slice() []int32 {
	set.mu.RLock()
	var i = 0
	var ret = make([]int32, len(set.data))
	for item := range set.data {
		ret[i] = item
		i++
	}
	set.mu.RUnlock()
	return ret
}

// LockFunc 锁住当前set调用方法f
func (set *Int32) LockFunc(f func(m map[int32]struct{})) {
	set.mu.Lock()
	defer set.mu.Unlock()
	f(set.data)
}

// RLockFunc 读锁住当前set调用方法f
func (set *Int32) RLockFunc(f func(m map[int32]struct{})) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	f(set.data)
}

// Equal 是否相等
func (set *Int32) Equal(other *Int32) bool {
	if set == other {
		return true
	}
	set.mu.RLock()
	defer set.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	if len(set.data) != len(other.data) {
		return false
	}
	for key := range set.data {
		if _, ok := other.data[key]; !ok {
			return false
		}
	}
	return true
}

// Merge 合并set，返回当前set
func (set *Int32) Merge(others ...*Int32) *Int32 {
	set.mu.Lock()
	defer set.mu.Unlock()
	for _, other := range others {
		if set != other {
			other.mu.RLock()
		}
		for k, v := range other.data {
			set.data[k] = v
		}
		if set != other {
			other.mu.RUnlock()
		}
	}
	return set
}

// Walk 对每个元素作用f方法
func (set *Int32) Walk(f func(item int32) int32) *Int32 {
	set.mu.Lock()
	defer set.mu.Unlock()
	m := make(map[int32]struct{}, len(set.data))
	for k, v := range set.data {
		m[f(k)] = v
	}
	set.data = m
	return set
}

type localRWMutexVTypeInt32 struct {
	*sync.RWMutex
}

func newLocalRWMutexVTypeInt32(safe bool) *localRWMutexVTypeInt32 {
	mu := localRWMutexVTypeInt32{}
	if safe {
		mu.RWMutex = new(sync.RWMutex)
	}
	return &mu
}

func (mu *localRWMutexVTypeInt32) IsSafe() bool {
	return mu.RWMutex != nil
}

func (mu *localRWMutexVTypeInt32) Lock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Lock()
	}
}

func (mu *localRWMutexVTypeInt32) Unlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Unlock()
	}
}

func (mu *localRWMutexVTypeInt32) RLock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RLock()
	}
}

func (mu *localRWMutexVTypeInt32) RUnlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RUnlock()
	}
}

//template format
var __formatToInt32 = func(i interface{}) int32 {
	switch ii := i.(type) {
	case int:
		return int32(ii)
	case int8:
		return int32(ii)
	case int16:
		return int32(ii)
	case int32:
		return int32(ii)
	case int64:
		return int32(ii)
	case uint:
		return int32(ii)
	case uint8:
		return int32(ii)
	case uint16:
		return int32(ii)
	case uint32:
		return int32(ii)
	case uint64:
		return int32(ii)
	case float32:
		return int32(ii)
	case float64:
		return int32(ii)
	case string:
		iv, err := strconv.ParseInt(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return int32(iv)
	default:
		panic("unknown type")
	}
}
