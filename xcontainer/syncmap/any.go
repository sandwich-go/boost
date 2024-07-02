package syncmap

import (
	"sort"
	"sync"
)

// Map 定义并发安全的映射，使用 sync.Map 来实现
type Map[K any, V any] struct {
	native sync.Map
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *Map[K, V]) Load(key K) (V, bool) {
	v, ok := m.native.Load(key)
	if !ok {
		var v2 V
		return v2, false
	}
	return v.(V), ok
}

// Store sets the value for a key.
func (m *Map[K, V]) Store(key K, value V) { m.native.Store(key, value) }

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, l := m.native.LoadOrStore(key, value)
	return a.(V), l
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, l := m.native.LoadAndDelete(key)
	if !l {
		var v2 V
		return v2, false
	}
	return v.(V), l
}

// Delete deletes the value for a key.
func (m *Map[K, V]) Delete(key K) { m.native.Delete(key) }

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.native.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (m *Map[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.native.CompareAndSwap(key, old, new)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) Swap(key K, new V) (V, bool) {
	prev, loaded := m.native.Swap(key, new)
	if !loaded {
		var v2 V
		return v2, false
	}
	return prev.(V), loaded
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
func (m *Map[K, V]) CompareAndDelete(key K, new V) bool {
	return m.native.CompareAndDelete(key, new)
}

// Keys 获取映射中的所有键，返回一个 Key 类型的切片
func (m *Map[K, V]) Keys() (ret []K) {
	m.native.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(K))
		return true
	})
	return ret
}

// Len 获取映射中键值对的数量
func (m *Map[K, V]) Len() (c int) {
	m.native.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

// Contains 检查映射中是否包含指定键
func (m *Map[K, V]) Contains(key K) (ok bool) {
	_, ok = m.Load(key)
	return
}

// Get 获取映射中的值
func (m *Map[K, V]) Get(key K) (value V) {
	value, _ = m.Load(key)
	return
}

// DeleteMultiple 删除映射中的多个键
func (m *Map[K, V]) DeleteMultiple(keys ...K) {
	for _, k := range keys {
		m.native.Delete(k)
	}
}

// Clear 清空映射
func (m *Map[K, V]) Clear() {
	m.native.Range(func(key, value interface{}) bool {
		m.native.Delete(key)
		return true
	})
}

// RangeDeterministic 按照 key 的顺序遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数 f, f返回false则中断退出
// 参数 sortableGetter 接收一个 mapKey 切片并返回一个可排序接口，用于对key进行排序
func (m *Map[K, V]) RangeDeterministic(f func(key K, value V) bool, sortableGetter func([]K) sort.Interface) {
	var keys []K
	m.native.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(K))
		return true
	})
	sort.Sort(sortableGetter(keys))
	for _, k := range keys {
		if v, ok := m.Load(k); ok {
			if !f(k, v) {
				break
			}
		}
	}
}

// SyncMap 定义并发安全的映射，使用 sync.Map 来实现
type SyncMap[K any, V any] struct {
	Map[K, V]
	locker sync.RWMutex
}

// New 构造函数，返回一个新的 SyncMap
func New[K any, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

// GetOrSetFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *SyncMap[K, V]) GetOrSetFuncErrorLock(key K, cf func(key K) (V, error)) (value V, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

// LoadOrStoreFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *SyncMap[K, V]) LoadOrStoreFuncErrorLock(key K, cf func(key K) (V, error)) (value V, loaded bool, err error) {
	if v, ok := s.Load(key); ok {
		return v, true, nil
	}
	// 如果不存在，则加写锁，再次查找，如果获取到则直接返回
	s.locker.Lock()
	defer s.locker.Unlock()
	// 再次重试，如果获取到则直接返回
	if v, ok := s.Load(key); ok {
		return v, true, nil
	}
	// 如果还是不存在，则执行cf函数计算出value，并存储到 SyncMap 中
	value, err = cf(key)
	if err != nil {
		return value, false, err
	}
	s.Store(key, value)
	return value, false, nil
}

// GetOrSetFuncLock 根据key获取对应的value，若不存在则通过cf回调创建value并存储
func (s *SyncMap[K, V]) GetOrSetFuncLock(key K, cf func(key K) V) (value V, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

// LoadOrStoreFuncLock 根据key获取对应的value，若不存在则通过cf回调创建value并存储
func (s *SyncMap[K, V]) LoadOrStoreFuncLock(key K, cf func(key K) V) (value V, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key K) (V, error) {
		return cf(key), nil
	})
	return value, loaded
}
