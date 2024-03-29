// Code generated by gotemplate. DO NOT EDIT.

// syncmap 提供了一个同步的映射实现，允许安全并发的访问
package syncmap

import (
	"sort"
	"sync"
)

//template type SyncMap(KType,VType)

// SyncMap 定义并发安全的映射，使用 sync.Map 来实现
type AnyAny struct {
	sm     sync.Map
	locker sync.RWMutex
}

// NewSyncMap 构造函数，返回一个新的 SyncMap
func NewAnyAny() *AnyAny {
	return &AnyAny{}
}

// Keys 获取映射中的所有键，返回一个 Key 类型的切片
func (s *AnyAny) Keys() (ret []interface{}) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(interface{}))
		return true
	})
	return ret
}

// Len 获取映射中键值对的数量
func (s *AnyAny) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

// Contains 检查映射中是否包含指定键
func (s *AnyAny) Contains(key interface{}) (ok bool) {
	_, ok = s.Load(key)
	return
}

// Get 获取映射中的值
func (s *AnyAny) Get(key interface{}) (value interface{}) {
	value, _ = s.Load(key)
	return
}

// Load 获取映射中的值和是否成功加载的标志
func (s *AnyAny) Load(key interface{}) (value interface{}, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(interface{}), true
	}
	return
}

// DeleteMultiple 删除映射中的多个键
func (s *AnyAny) DeleteMultiple(keys ...interface{}) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}

// Clear 清空映射
func (s *AnyAny) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}

// Delete 删除映射中的值
func (s *AnyAny) Delete(key interface{}) { s.sm.Delete(key) }

// Store 往映射中存储一个键值对
func (s *AnyAny) Store(key interface{}, val interface{}) { s.sm.Store(key, val) }

// LoadAndDelete 获取映射中的值，并将其从映射中删除
func (s *AnyAny) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(interface{}), true
	}
	return
}

// GetOrSetFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *AnyAny) GetOrSetFuncErrorLock(key interface{}, cf func(key interface{}) (interface{}, error)) (value interface{}, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

// LoadOrStoreFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *AnyAny) LoadOrStoreFuncErrorLock(key interface{}, cf func(key interface{}) (interface{}, error)) (value interface{}, loaded bool, err error) {
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
func (s *AnyAny) GetOrSetFuncLock(key interface{}, cf func(key interface{}) interface{}) (value interface{}, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

// LoadOrStoreFuncLock 根据key获取对应的value，若不存在则通过cf回调创建value并存储
func (s *AnyAny) LoadOrStoreFuncLock(key interface{}, cf func(key interface{}) interface{}) (value interface{}, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key interface{}) (interface{}, error) {
		return cf(key), nil
	})
	return value, loaded
}

// LoadOrStore 存储一个 key-value 对，若key已存在则返回已存在的value
func (s *AnyAny) LoadOrStore(key interface{}, val interface{}) (interface{}, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(interface{}), ok
}

// Range 遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数f
func (s *AnyAny) Range(f func(key interface{}, value interface{}) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(interface{}), v.(interface{}))
	})
}

// RangeDeterministic 按照 key 的顺序遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数 f, f返回false则中断退出
// 参数 sortableGetter 接收一个 KType 切片并返回一个可排序接口，用于对key进行排序
func (s *AnyAny) RangeDeterministic(f func(key interface{}, value interface{}) bool, sortableGetter func([]interface{}) sort.Interface) {
	var keys []interface{}
	s.sm.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(interface{}))
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

//template format
var __formatKTypeToAnyAny = func(i interface{}) interface{} {
	return i
}

//template format
var __formatVTypeToAnyAny = func(i interface{}) interface{} {
	return i
}

// add 6,6 success

// add 7, 7 failed

// add 7, 7 success
