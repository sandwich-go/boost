package sset

import (
	"sync"
)

//template type SandwichSet(VType)

type VType interface{}

type SandwichSet struct {
	mu   *localRWMutexVType
	data map[VType]struct{}
}

// New 创建非协程安全版本
func New() *SandwichSet { return newWithSafe(false) }

// NewSync 创建协程安全版本
func NewSync() *SandwichSet { return newWithSafe(true) }

func newWithSafe(safe bool) *SandwichSet {
	return &SandwichSet{data: make(map[VType]struct{}), mu: newLocalRWMutexVType(safe)}
}

// Iterator 遍历
func (set *SandwichSet) Iterator(f func(v VType) bool) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k := range set.data {
		if !f(k) {
			break
		}
	}
}

// Add 添加元素
func (set *SandwichSet) Add(items ...VType) {
	set.mu.Lock()
	if set.data == nil {
		set.data = make(map[VType]struct{})
	}
	for _, v := range items {
		set.data[v] = struct{}{}
	}
	set.mu.Unlock()
}

// AddIfNotExist 如果元素不存在则添加，如添加成功则返回true
func (set *SandwichSet) AddIfNotExist(item VType) (addOK bool) {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[VType]struct{})
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
func (set *SandwichSet) AddIfNotExistFunc(item VType, f func() bool) bool {
	if !set.Contains(item) {
		if f() {
			set.mu.Lock()
			defer set.mu.Unlock()
			if set.data == nil {
				set.data = make(map[VType]struct{})
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
func (set *SandwichSet) AddIfNotExistFuncLock(item VType, f func() bool) bool {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[VType]struct{})
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
func (set *SandwichSet) Contains(item VType) bool {
	var ok bool
	set.mu.RLock()
	if set.data != nil {
		_, ok = set.data[item]
	}
	set.mu.RUnlock()
	return ok
}

// Remove 移除指定元素
func (set *SandwichSet) Remove(item VType) {
	set.mu.Lock()
	if set.data != nil {
		delete(set.data, item)
	}
	set.mu.Unlock()
}

// Size 返回长度
func (set *SandwichSet) Size() int {
	set.mu.RLock()
	l := len(set.data)
	set.mu.RUnlock()
	return l
}

// Clear 清理元素
func (set *SandwichSet) Clear() {
	set.mu.Lock()
	set.data = make(map[VType]struct{})
	set.mu.Unlock()
}

// Slice 返回元素slice
func (set *SandwichSet) Slice() []VType {
	set.mu.RLock()
	var i = 0
	var ret = make([]VType, len(set.data))
	for item := range set.data {
		ret[i] = item
		i++
	}
	set.mu.RUnlock()
	return ret
}

// LockFunc 锁住当前set调用方法f
func (set *SandwichSet) LockFunc(f func(m map[VType]struct{})) {
	set.mu.Lock()
	defer set.mu.Unlock()
	f(set.data)
}

// RLockFunc 读锁住当前set调用方法f
func (set *SandwichSet) RLockFunc(f func(m map[VType]struct{})) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	f(set.data)
}

// Equal 是否相等
func (set *SandwichSet) Equal(other *SandwichSet) bool {
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

// Merge 合并set，发布会当前set
func (set *SandwichSet) Merge(others ...*SandwichSet) *SandwichSet {
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
func (set *SandwichSet) Walk(f func(item VType) VType) *SandwichSet {
	set.mu.Lock()
	defer set.mu.Unlock()
	m := make(map[VType]struct{}, len(set.data))
	for k, v := range set.data {
		m[f(k)] = v
	}
	set.data = m
	return set
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
