// Code generated by gotemplate. DO NOT EDIT.

package smap

import (
	"fmt"
	"strconv"
	"sync"
)

//template type Concurrent(KType,VType,KeyHash)

// A thread safe map.
// To avoid lock bottlenecks this map is dived to several (DefaultShardCount) map shards.

var DefaultShardCountIntString = uint64(32)

type IntString struct {
	shardedList  []*shardedIntString
	shardedCount uint64
}

type shardedIntString struct {
	items map[int]string
	sync.RWMutex
}

type TupleIntString struct {
	Key int
	Val string
}

// NewWithSharedCount 返回协程安全版本
func NewWithSharedCountIntString(sharedCount uint64) *IntString {
	p := &IntString{
		shardedCount: sharedCount,
		shardedList:  make([]*shardedIntString, sharedCount),
	}
	for i := uint64(0); i < sharedCount; i++ {
		p.shardedList[i] = &shardedIntString{items: make(map[int]string)}
	}
	return p
}

// New 返回协程安全版本
func NewIntString() *IntString {
	return NewWithSharedCountIntString(DefaultShardCountIntString)
}

// GetShard 返回key对应的分片
func (m *IntString) GetShard(key int) *shardedIntString {
	return m.shardedList[KeyToHashInt(key)%m.shardedCount]
}

// IsEmpty 返回容器是否为空
func (m *IntString) IsEmpty() bool {
	return m.Count() == 0
}

// Set 设定元素
func (m *IntString) Set(key int, value string) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// Keys 返回所有的key列表
func (m *IntString) Keys() []int {
	var ret []int
	for _, shard := range m.shardedList {
		shard.RLock()
		for key := range shard.items {
			ret = append(ret, key)
		}
		shard.RUnlock()
	}
	return ret
}

// GetAll 返回所有元素副本，其中value浅拷贝
func (m *IntString) GetAll() map[int]string {
	data := make(map[int]string)
	for _, shard := range m.shardedList {
		shard.RLock()
		for key, val := range shard.items {
			data[key] = val
		}
		shard.RUnlock()
	}
	return data
}

// Clear 清空元素
func (m *IntString) Clear() {
	for _, shard := range m.shardedList {
		shard.Lock()
		shard.items = make(map[int]string)
		shard.Unlock()
	}
}

// ClearWithFunc 清空元素,onClear在对应分片的Lock内执行，执行完毕后对容器做整体clear操作
//
// Note: 不要在onClear对当前容器做读写操作，容易死锁
//
// data.ClearWithFuncLock(func(key string,val string) {
//		data.Get(...) // 死锁
// })
//
func (m *IntString) ClearWithFuncLock(onClear func(key int, val string)) {
	for _, shard := range m.shardedList {
		shard.Lock()
		for key, val := range shard.items {
			onClear(key, val)
		}
		shard.items = make(map[int]string)
		shard.Unlock()
	}
}

// MGet 返回多个元素
func (m *IntString) MGet(keys ...int) map[int]string {
	data := make(map[int]string)
	for _, key := range keys {
		if val, ok := m.Get(key); ok {
			data[key] = val
		}
	}
	return data
}

// MSet 同时设定多个元素
func (m *IntString) MSet(data map[int]string) {
	for key, value := range data {
		m.Set(key, value)
	}
}

// SetNX 如果key不存在，则设定为value, 设定成功则返回true，否则返回false
func (m *IntString) SetNX(key int, value string) (isSet bool) {
	shard := m.GetShard(key)
	shard.Lock()
	if _, ok := shard.items[key]; !ok {
		shard.items[key] = value
		isSet = true
	}
	shard.Unlock()
	return isSet
}

// LockFuncWithKey 对key对应的分片加写锁，并用f操作该分片内数据
//
// Note: 不要对f中对容器的该分片做读写操作，可以直接操作shardData数据源
//
//  data.LockFuncWithKey("test",func(shardData map[string]string) {
//     data.Remove("test")      // 当前分片已被加读锁, 死锁
//  })
//
func (m *IntString) LockFuncWithKey(key int, f func(shardData map[int]string)) {
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	f(shard.items)
}

// RLockFuncWithKey 对key对应的分片加读锁，并用f操作该分片内数据
//
// Note: 不要在f内对容器做写操作，否则会引起死锁，可以直接操作shardData数据源
//
//  data.RLockFuncWithKey("test",func(shardData map[string]string) {
//     data.Remove("test")      // 当前分片已被加读锁, 死锁
//  })
//
func (m *IntString) RLockFuncWithKey(key int, f func(shardData map[int]string)) {
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()
	f(shard.items)
}

// LockFunc 遍历容器分片，f在Lock写锁内执行
//
// Note: 不要在f内对容器做读写操作，否则会引起死锁，可以直接操作shardData数据源
//
//  data.LockFunc(func(shardData map[string]string) {
//     data.Count()             // 当前分片已被加写锁, 死锁
//  })
//
func (m *IntString) LockFunc(f func(shardData map[int]string)) {
	for _, shard := range m.shardedList {
		shard.Lock()
		f(shard.items)
		shard.Unlock()
	}
}

// RLockFunc 遍历容器分片，f在RLock读锁内执行
//
// Note: 不要在f内对容器做修改操作，否则会引起死锁，可以直接操作shardData数据源
//
//  data.RLockFunc(func(shardData map[string]string) {
//     data.Remove("test")      // 当前分片已被加读锁, 死锁
//  })
//
func (m *IntString) RLockFunc(f func(shardData map[int]string)) {
	for _, shard := range m.shardedList {
		shard.RLock()
		f(shard.items)
		shard.RUnlock()
	}
}

func (m *IntString) doSetWithLockCheck(key int, val string) (result string, isSet bool) {
	shard := m.GetShard(key)
	shard.Lock()

	if got, ok := shard.items[key]; ok {
		shard.Unlock()
		return got, false
	}

	shard.items[key] = val
	isSet = true
	shard.Unlock()
	return
}

func (m *IntString) doSetWithLockCheckWithFunc(key int, f func(key int) string) (result string, isSet bool) {
	shard := m.GetShard(key)
	shard.Lock()

	if got, ok := shard.items[key]; ok {
		shard.Unlock()
		return got, false
	}

	shard.items[key] = f(key)
	isSet = true
	shard.Unlock()
	return
}

// GetOrSetFunc 获取或者设定数值，方法f在Lock写锁外执行, 如元素早已存在则返回false,设定成功返回true
func (m *IntString) GetOrSetFunc(key int, f func(key int) string) (result string, isSet bool) {
	if v, ok := m.Get(key); ok {
		return v, false
	}
	return m.doSetWithLockCheck(key, f(key))
}

// GetOrSetFuncLock 获取或者设定数值，方法f在Lock写锁内执行, 如元素早已存在则返回false,设定成功返回true
//
// Note: 不要在f内对容器做操作，否则会死锁
//
//  data.GetOrSetFuncLock(“test”,func(key string)string {
//     data.Count() // 死锁
//  })
//
func (m *IntString) GetOrSetFuncLock(key int, f func(key int) string) (result string, isSet bool) {
	if v, ok := m.Get(key); ok {
		return v, false
	}
	return m.doSetWithLockCheckWithFunc(key, f)
}

// GetOrSet 获取或设定元素, 如元素早已存在则返回false,设定成功返回true
func (m *IntString) GetOrSet(key int, val string) (string, bool) {
	if v, ok := m.Get(key); ok {
		return v, false
	}
	return m.doSetWithLockCheck(key, val)
}

// Get 返回key对应的元素，不存在返回false
func (m *IntString) Get(key int) (string, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Len Count方法别名
func (m *IntString) Len() int { return m.Count() }

// Size Count方法别名
func (m *IntString) Size() int { return m.Count() }

// Count 返回容器内元素数量
func (m *IntString) Count() int {
	count := 0
	for i := uint64(0); i < m.shardedCount; i++ {
		shard := m.shardedList[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Has 是否存在key对应的元素
func (m *IntString) Has(key int) bool {
	shard := m.GetShard(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Remove 删除key对应的元素
func (m *IntString) Remove(key int) {
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// GetAndRemove 返回key对应的元素并将其由容器中删除，如果元素不存在则返回false
func (m *IntString) GetAndRemove(key int) (string, bool) {
	shard := m.GetShard(key)
	shard.Lock()
	val, ok := shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return val, ok
}

// Iter 迭代当前容器内所有元素，使用无缓冲chan
//
// Note: 不要在迭代过程中对当前容器作修改操作(申请写锁)，容易会产生死锁
//
//  for v:= data.Iter() {
//		data.Remove(v.Key) // 尝试删除元素申请分片Lock,但是Iter内部的迭代协程对分片做了RLock，导致死锁
//  }
//
func (m *IntString) Iter() <-chan TupleIntString {
	ch := make(chan TupleIntString)
	go func() {
		// Foreach shard.
		for _, shard := range m.shardedList {
			shard.RLock()
			// Foreach key, value pair.
			for key, val := range shard.items {
				ch <- TupleIntString{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// IterBuffered 迭代当前容器内所有元素，使用有缓冲chan，缓冲区大小等于容器大小,迭代过程中操作容器是安全的
func (m *IntString) IterBuffered() <-chan TupleIntString {
	ch := make(chan TupleIntString, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m.shardedList {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- TupleIntString{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

//template format
var __formatKTypeToIntString = func(i interface{}) int {
	switch ii := i.(type) {
	case int:
		return int(ii)
	case int8:
		return int(ii)
	case int16:
		return int(ii)
	case int32:
		return int(ii)
	case int64:
		return int(ii)
	case uint:
		return int(ii)
	case uint8:
		return int(ii)
	case uint16:
		return int(ii)
	case uint32:
		return int(ii)
	case uint64:
		return int(ii)
	case float32:
		return int(ii)
	case float64:
		return int(ii)
	case string:
		iv, err := strconv.ParseInt(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return int(iv)
	default:
		panic("unknown type")
	}
}

//template format
var __formatVTypeToIntString = func(i interface{}) string {
	switch ii := i.(type) {
	case string:
		return ii
	default:
		return fmt.Sprintf("%d", i)
	}
}

//So(r, ShouldEqual, __formatVTypeTo(10))
