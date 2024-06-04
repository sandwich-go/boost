package syncmap

import (
	"sort"
	"sync"
)

// SyncMap 定义并发安全的映射，使用 sync.Map 来实现
type SyncMap[K any, V any] struct {
	sm     sync.Map
	locker sync.RWMutex
}

// New 构造函数，返回一个新的 SyncMap
func New[K any, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

// Keys 获取映射中的所有键，返回一个 Key 类型的切片
func (s *SyncMap[K, V]) Keys() (ret []K) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(K))
		return true
	})
	return ret
}

// Len 获取映射中键值对的数量
func (s *SyncMap[K, V]) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

// Contains 检查映射中是否包含指定键
func (s *SyncMap[K, V]) Contains(key K) (ok bool) {
	_, ok = s.Load(key)
	return
}

// Get 获取映射中的值
func (s *SyncMap[K, V]) Get(key K) (value V) {
	value, _ = s.Load(key)
	return
}

// Load 获取映射中的值和是否成功加载的标志
func (s *SyncMap[K, V]) Load(key K) (value V, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(V), true
	}
	return
}

// DeleteMultiple 删除映射中的多个键
func (s *SyncMap[K, V]) DeleteMultiple(keys ...K) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}

// Clear 清空映射
func (s *SyncMap[K, V]) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}

// Delete 删除映射中的值
func (s *SyncMap[K, V]) Delete(key K) { s.sm.Delete(key) }

// Store 往映射中存储一个键值对
func (s *SyncMap[K, V]) Store(key K, val V) { s.sm.Store(key, val) }

// LoadAndDelete 获取映射中的值，并将其从映射中删除
func (s *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(V), true
	}
	return
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

// LoadOrStore 存储一个 key-value 对，若key已存在则返回已存在的value
func (s *SyncMap[K, V]) LoadOrStore(key K, val V) (V, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(V), ok
}

// Range 遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数f
func (s *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(K), v.(V))
	})
}

// RangeDeterministic 按照 key 的顺序遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数 f, f返回false则中断退出
// 参数 sortableGetter 接收一个 mapKey 切片并返回一个可排序接口，用于对key进行排序
func (s *SyncMap[K, V]) RangeDeterministic(f func(key K, value V) bool, sortableGetter func([]K) sort.Interface) {
	var keys []K
	s.sm.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(K))
		return true
	})
	sort.Sort(sortableGetter(keys))
	for _, k := range keys {
		if v, ok := s.Load(k); ok {
			if !f(k, v) {
				break
			}
		}
	}
}
