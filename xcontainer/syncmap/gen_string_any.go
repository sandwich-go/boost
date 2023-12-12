// Code generated by gotemplate. DO NOT EDIT.

// syncmap 提供了一个同步的映射实现，允许安全并发的访问
package syncmap

import (
	"fmt"
	"sort"
	"sync"
)

//template type SyncMap(KType,VType)

// SyncMap 定义并发安全的映射，使用 sync.Map 来实现
type StringAny struct {
	sm     sync.Map
	locker sync.RWMutex
}

// NewSyncMap 构造函数，返回一个新的 SyncMap
func NewStringAny() *StringAny {
	return &StringAny{}
}

// Keys 获取映射中的所有键，返回一个 Key 类型的切片
func (s *StringAny) Keys() (ret []string) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(string))
		return true
	})
	return ret
}

// Len 获取映射中键值对的数量
func (s *StringAny) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

// Contains 检查映射中是否包含指定键
func (s *StringAny) Contains(key string) (ok bool) {
	_, ok = s.Load(key)
	return
}

// Get 获取映射中的值
func (s *StringAny) Get(key string) (value interface{}) {
	value, _ = s.Load(key)
	return
}

// Load 获取映射中的值和是否成功加载的标志
func (s *StringAny) Load(key string) (value interface{}, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(interface{}), true
	}
	return
}

// DeleteMultiple 删除映射中的多个键
func (s *StringAny) DeleteMultiple(keys ...string) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}

// Clear 清空映射
func (s *StringAny) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}

// Delete 删除映射中的值
func (s *StringAny) Delete(key string) { s.sm.Delete(key) }

// Store 往映射中存储一个键值对
func (s *StringAny) Store(key string, val interface{}) { s.sm.Store(key, val) }

// LoadAndDelete 获取映射中的值，并将其从映射中删除
func (s *StringAny) LoadAndDelete(key string) (value interface{}, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(interface{}), true
	}
	return
}

// GetOrSetFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *StringAny) GetOrSetFuncErrorLock(key string, cf func(key string) (interface{}, error)) (value interface{}, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

// LoadOrStoreFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *StringAny) LoadOrStoreFuncErrorLock(key string, cf func(key string) (interface{}, error)) (value interface{}, loaded bool, err error) {
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
func (s *StringAny) GetOrSetFuncLock(key string, cf func(key string) interface{}) (value interface{}, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

// LoadOrStoreFuncLock 根据key获取对应的value，若不存在则通过cf回调创建value并存储
func (s *StringAny) LoadOrStoreFuncLock(key string, cf func(key string) interface{}) (value interface{}, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key string) (interface{}, error) {
		return cf(key), nil
	})
	return value, loaded
}

// LoadOrStore 存储一个 key-value 对，若key已存在则返回已存在的value
func (s *StringAny) LoadOrStore(key string, val interface{}) (interface{}, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(interface{}), ok
}

// Range 遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数f
func (s *StringAny) Range(f func(key string, value interface{}) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(string), v.(interface{}))
	})
}

// RangeDeterministic 按照 key 的顺序遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数 f, f返回false则中断退出
// 参数 sortableGetter 接收一个 KType 切片并返回一个可排序接口，用于对key进行排序
func (s *StringAny) RangeDeterministic(f func(key string, value interface{}) bool, sortableGetter func([]string) sort.Interface) {
	var keys []string
	s.sm.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
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
var __formatKTypeToStringAny = func(i interface{}) string {
	switch ii := i.(type) {
	case string:
		return ii
	default:
		return fmt.Sprintf("%d", i)
	}
}

//template format
var __formatVTypeToStringAny = func(i interface{}) interface{} {
	return i
}

// add 6,6 success

// add 7, 7 failed

// add 7, 7 success