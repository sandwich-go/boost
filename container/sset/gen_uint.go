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

type Uint struct {
	mu   *localRWMutexVTypeUint
	data map[uint]struct{}
}

// New 创建非协程安全版本
func NewUint() *Uint { return newWithSafeUint(false) }

// NewSync 创建协程安全版本
func NewSyncUint() *Uint { return newWithSafeUint(true) }

func newWithSafeUint(safe bool) *Uint {
	return &Uint{data: make(map[uint]struct{}), mu: newLocalRWMutexVTypeUint(safe)}
}

// Iterator 遍历
func (set *Uint) Iterator(f func(v uint) bool) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k := range set.data {
		if !f(k) {
			break
		}
	}
}

// Add 添加元素
func (set *Uint) Add(items ...uint) {
	set.mu.Lock()
	if set.data == nil {
		set.data = make(map[uint]struct{})
	}
	for _, v := range items {
		set.data[v] = struct{}{}
	}
	set.mu.Unlock()
}

// AddIfNotExist 如果元素不存在则添加，如添加成功则返回true
func (set *Uint) AddIfNotExist(item uint) (addOK bool) {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[uint]struct{})
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
func (set *Uint) AddIfNotExistFunc(item uint, f func() bool) bool {
	if !set.Contains(item) {
		if f() {
			set.mu.Lock()
			defer set.mu.Unlock()
			if set.data == nil {
				set.data = make(map[uint]struct{})
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
func (set *Uint) AddIfNotExistFuncLock(item uint, f func() bool) bool {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[uint]struct{})
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
func (set *Uint) Contains(item uint) bool {
	var ok bool
	set.mu.RLock()
	if set.data != nil {
		_, ok = set.data[item]
	}
	set.mu.RUnlock()
	return ok
}

// Remove 移除指定元素
func (set *Uint) Remove(item uint) {
	set.mu.Lock()
	if set.data != nil {
		delete(set.data, item)
	}
	set.mu.Unlock()
}

// Size 返回长度
func (set *Uint) Size() int {
	set.mu.RLock()
	l := len(set.data)
	set.mu.RUnlock()
	return l
}

// Clear 清理元素
func (set *Uint) Clear() {
	set.mu.Lock()
	set.data = make(map[uint]struct{})
	set.mu.Unlock()
}

// Slice 返回元素slice
func (set *Uint) Slice() []uint {
	set.mu.RLock()
	var i = 0
	var ret = make([]uint, len(set.data))
	for item := range set.data {
		ret[i] = item
		i++
	}
	set.mu.RUnlock()
	return ret
}

// LockFunc 锁住当前set调用方法f
func (set *Uint) LockFunc(f func(m map[uint]struct{})) {
	set.mu.Lock()
	defer set.mu.Unlock()
	f(set.data)
}

// RLockFunc 读锁住当前set调用方法f
func (set *Uint) RLockFunc(f func(m map[uint]struct{})) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	f(set.data)
}

// Equal 是否相等
func (set *Uint) Equal(other *Uint) bool {
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
func (set *Uint) Merge(others ...*Uint) *Uint {
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
func (set *Uint) Walk(f func(item uint) uint) *Uint {
	set.mu.Lock()
	defer set.mu.Unlock()
	m := make(map[uint]struct{}, len(set.data))
	for k, v := range set.data {
		m[f(k)] = v
	}
	set.data = m
	return set
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
