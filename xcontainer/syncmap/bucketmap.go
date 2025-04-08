package syncmap

import (
	"sort"
	"sync"
	"sync/atomic"
)

// BucketMap 定义并发安全的分桶映射，使用 sync.Map 来实现
type BucketMap[K comparable, V any] struct {
	hashFunc  func(K) int64
	bucketNum int64
	natives   []sync.Map
	locks     []sync.RWMutex

	keyCaches []atomic.Pointer[[]K]
	cacheKey  bool
}

func NewBucketMapWithCacheKey[K comparable, V any](bucketNum int, hashFunc func(K) int64) *BucketMap[K, V] {
	return newBucketMap[K, V](bucketNum, true, hashFunc)
}
func NewBucketMap[K comparable, V any](bucketNum int, hashFunc func(K) int64) *BucketMap[K, V] {
	return newBucketMap[K, V](bucketNum, false, hashFunc)
}
func newBucketMap[K comparable, V any](bucketNum int, cacheKey bool, hashFunc func(K) int64) *BucketMap[K, V] {
	if bucketNum <= 0 {
		bucketNum = 1
	}
	b := &BucketMap[K, V]{hashFunc: hashFunc, bucketNum: int64(bucketNum)}
	b.natives = make([]sync.Map, bucketNum)
	b.locks = make([]sync.RWMutex, bucketNum)
	b.cacheKey = cacheKey
	if cacheKey {
		b.keyCaches = make([]atomic.Pointer[[]K], bucketNum)
		for i := range b.keyCaches {
			b.keyCaches[i].Store(new([]K))
		}
	}
	return b
}
func (m *BucketMap[K, V]) indexByKey(key K) int {
	if m.bucketNum == 1 {
		return 0
	}
	i := m.hashFunc(key)
	if i <= 0 {
		return 0
	}
	return int(i % m.bucketNum)
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *BucketMap[K, V]) Load(key K) (V, bool) {
	v, ok := m.natives[m.indexByKey(key)].Load(key)
	if !ok {
		var v2 V
		return v2, false
	}
	return v.(V), ok
}

// Store sets the value for a key.
func (m *BucketMap[K, V]) Store(key K, value V) {
	idx := m.indexByKey(key)
	m.natives[idx].Store(key, value)
	m.addKeyToCache(idx, key)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *BucketMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	idx := m.indexByKey(key)
	a, l := m.natives[idx].LoadOrStore(key, value)
	if !loaded {
		m.addKeyToCache(idx, key)
	}
	return a.(V), l
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *BucketMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	idx := m.indexByKey(key)
	v, isLoaded := m.natives[idx].LoadAndDelete(key)
	if !isLoaded {
		var v2 V
		return v2, false
	}
	// is delete
	m.removeKeyFromCache(idx, key)
	return v.(V), isLoaded
}

// Delete deletes the value for a key.
func (m *BucketMap[K, V]) Delete(key K) {
	idx := m.indexByKey(key)
	m.natives[idx].Delete(key)
	m.removeKeyFromCache(idx, key)
}

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
func (m *BucketMap[K, V]) Range(f func(key K, value V) bool) {
	for i := range m.natives {
		m.natives[i].Range(func(key, value any) bool {
			return f(key.(K), value.(V))
		})
	}
}

func (m *BucketMap[K, V]) BucketNum() int {
	return len(m.natives)
}

func (m *BucketMap[K, V]) RangeBucket(i int, f func(key K, value V) bool) {
	m.natives[i].Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (m *BucketMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.natives[m.indexByKey(key)].CompareAndSwap(key, old, new)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *BucketMap[K, V]) Swap(key K, new V) (V, bool) {
	prev, loaded := m.natives[m.indexByKey(key)].Swap(key, new)
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
func (m *BucketMap[K, V]) CompareAndDelete(key K, new V) bool {
	idx := m.indexByKey(key)
	deleted := m.natives[idx].CompareAndDelete(key, new)
	if deleted {
		m.removeKeyFromCache(idx, key)
	}
	return deleted
}

func (m *BucketMap[K, V]) Keys() (ret []K) {
	if m.cacheKey {
		for i := range m.keyCaches {
			ptr := m.keyCaches[i].Load()
			ret = append(ret, *ptr...)
		}
		return ret
	}
	for i := range m.natives {
		m.natives[i].Range(func(key, value interface{}) bool {
			ret = append(ret, key.(K))
			return true
		})
	}
	return ret
}

// Len 获取映射中键值对的数量
func (m *BucketMap[K, V]) Len() (c int) {
	if m.cacheKey {
		for i := range m.keyCaches {
			ptr := m.keyCaches[i].Load()
			if ptr != nil {
				c += len(*ptr)
			}
		}
		return c
	}
	for i := range m.natives {
		m.natives[i].Range(func(key, value interface{}) bool {
			c++
			return true
		})
	}
	return c
}

// Contains 检查映射中是否包含指定键
func (m *BucketMap[K, V]) Contains(key K) (ok bool) {
	_, ok = m.Load(key)
	return
}

// Get 获取映射中的值
func (m *BucketMap[K, V]) Get(key K) (value V) {
	value, _ = m.Load(key)
	return
}

// DeleteMultiple 删除映射中的多个键
func (m *BucketMap[K, V]) DeleteMultiple(keys ...K) {
	for _, k := range keys {
		idx := m.indexByKey(k)
		_, isLoaded := m.natives[idx].LoadAndDelete(k)
		if isLoaded {
			m.removeKeyFromCache(idx, k)
		}
	}
}

func (m *BucketMap[K, V]) clear(index int) {
	m.natives[index].Range(func(key, value interface{}) bool {
		m.natives[index].Delete(key)
		if m.cacheKey {
			m.keyCaches[index].Store(new([]K))
		}
		return true
	})
}

// Clear 清空映射
func (m *BucketMap[K, V]) Clear() {
	for i := range m.natives {
		m.clear(i)
	}
}

// RangeDeterministic 按照 key 的顺序遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数 f, f返回false则中断退出
// 参数 sortableGetter 接收一个 mapKey 切片并返回一个可排序接口，用于对key进行排序
func (m *BucketMap[K, V]) RangeDeterministic(f func(key K, value V) bool, sortableGetter func([]K) sort.Interface) {
	var keys []K
	for i := range m.natives {
		m.natives[i].Range(func(key, value interface{}) bool {
			keys = append(keys, key.(K))
			return true
		})
	}
	sort.Sort(sortableGetter(keys))
	for _, k := range keys {
		if v, ok := m.Load(k); ok {
			if !f(k, v) {
				break
			}
		}
	}
}

// LoadOrStoreFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (m *BucketMap[K, V]) LoadOrStoreFuncErrorLock(key K, newValFunc func(key K) (V, error)) (value V, loaded bool, err error) {
	idx := m.indexByKey(key)
	bucket := &m.natives[idx]
	if val, ok := bucket.Load(key); ok {
		return val.(V), true, nil
	}
	m.locks[idx].Lock()
	defer m.locks[idx].Unlock()

	if val, ok := bucket.Load(key); ok {
		return val.(V), true, nil
	}
	v, err := newValFunc(key)
	if err != nil {
		return v, false, err
	}
	bucket.Store(key, v)
	m.addKeyToCache(idx, key)
	return v, false, nil
}

// LoadOrStoreFuncLock 根据key获取对应的value，若不存在则通过cf回调创建value并存储
func (m *BucketMap[K, V]) LoadOrStoreFuncLock(key K, cf func(key K) V) (value V, loaded bool) {
	value, loaded, _ = m.LoadOrStoreFuncErrorLock(key, func(key K) (V, error) {
		return cf(key), nil
	})
	return value, loaded
}
func (m *BucketMap[K, V]) addKeyToCache(index int, key K) {
	if !m.cacheKey {
		return
	}
	for {
		oldPtr := m.keyCaches[index].Load()
		var old []K
		if oldPtr != nil {
			old = *oldPtr
			for _, k := range old {
				if k == key {
					return
				}
			}
		}

		newSlice := make([]K, len(old), len(old)+1)
		copy(newSlice, old)
		newSlice = append(newSlice, key)

		newPtr := &newSlice
		if m.keyCaches[index].CompareAndSwap(oldPtr, newPtr) {
			return // 更新成功，退出
		}
		// 否则重新再来一次
	}
}
func (m *BucketMap[K, V]) removeKeyFromCache(index int, key K) {
	if !m.cacheKey {
		return
	}
	for {
		oldPtr := m.keyCaches[index].Load()
		if oldPtr == nil {
			return // 无内容，直接返回
		}
		old := *oldPtr

		found := false
		newSlice := make([]K, 0, len(old))
		for _, k := range old {
			if k == key {
				found = true
				continue
			}
			newSlice = append(newSlice, k)
		}
		if !found {
			return // key 不存在，直接返回
		}

		newPtr := &newSlice
		if m.keyCaches[index].CompareAndSwap(oldPtr, newPtr) {
			return // 更新成功
		}
		// 否则重新获取最新值，再试一次
	}
}
