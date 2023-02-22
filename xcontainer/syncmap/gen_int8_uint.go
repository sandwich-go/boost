// Code generated by gotemplate. DO NOT EDIT.

// syncmap 提供了一个同步的映射实现，允许安全并发的访问
package syncmap

import (
	"sort"
	"strconv"
	"sync"
)

//template type SyncMap(KType,VType)

// SyncMap 定义并发安全的映射，使用 sync.Map 来实现
type Int8Uint struct {
	sm     sync.Map
	locker sync.RWMutex
}

// NewSyncMap 构造函数，返回一个新的 SyncMap
func NewInt8Uint() *Int8Uint {
	return &Int8Uint{}
}

// Keys 获取映射中的所有键，返回一个 Key 类型的切片
func (s *Int8Uint) Keys() (ret []int8) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(int8))
		return true
	})
	return ret
}

// Len 获取映射中键值对的数量
func (s *Int8Uint) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

// Contains 检查映射中是否包含指定键
func (s *Int8Uint) Contains(key int8) (ok bool) {
	_, ok = s.Load(key)
	return
}

// Get 获取映射中的值
func (s *Int8Uint) Get(key int8) (value uint) {
	value, _ = s.Load(key)
	return
}

// Load 获取映射中的值和是否成功加载的标志
func (s *Int8Uint) Load(key int8) (value uint, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(uint), true
	}
	return
}

// DeleteMultiple 删除映射中的多个键
func (s *Int8Uint) DeleteMultiple(keys ...int8) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}

// Clear 清空映射
func (s *Int8Uint) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}

// Delete 删除映射中的值
func (s *Int8Uint) Delete(key int8) { s.sm.Delete(key) }

// Store 往映射中存储一个键值对
func (s *Int8Uint) Store(key int8, val uint) { s.sm.Store(key, val) }

// LoadAndDelete 获取映射中的值，并将其从映射中删除
func (s *Int8Uint) LoadAndDelete(key int8) (value uint, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(uint), true
	}
	return
}

// GetOrSetFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *Int8Uint) GetOrSetFuncErrorLock(key int8, cf func(key int8) (uint, error)) (value uint, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

// LoadOrStoreFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *Int8Uint) LoadOrStoreFuncErrorLock(key int8, cf func(key int8) (uint, error)) (value uint, loaded bool, err error) {
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
func (s *Int8Uint) GetOrSetFuncLock(key int8, cf func(key int8) uint) (value uint, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

// LoadOrStoreFuncLock 根据key获取对应的value，若不存在则通过cf回调创建value并存储
func (s *Int8Uint) LoadOrStoreFuncLock(key int8, cf func(key int8) uint) (value uint, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key int8) (uint, error) {
		return cf(key), nil
	})
	return value, loaded
}

// LoadOrStore 存储一个 key-value 对，若key已存在则返回已存在的value
func (s *Int8Uint) LoadOrStore(key int8, val uint) (uint, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(uint), ok
}

// Range 遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数f
func (s *Int8Uint) Range(f func(key int8, value uint) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(int8), v.(uint))
	})
}

// RangeDeterministic 按照 key 的顺序遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数 f, f返回false则中断退出
// 参数 sortableGetter 接收一个 KType 切片并返回一个可排序接口，用于对key进行排序
func (s *Int8Uint) RangeDeterministic(f func(key int8, value uint) bool, sortableGetter func([]int8) sort.Interface) {
	var keys []int8
	s.sm.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(int8))
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
var __formatKTypeToInt8Uint = func(i interface{}) int8 {
	switch ii := i.(type) {
	case int:
		return int8(ii)
	case int8:
		return int8(ii)
	case int16:
		return int8(ii)
	case int32:
		return int8(ii)
	case int64:
		return int8(ii)
	case uint:
		return int8(ii)
	case uint8:
		return int8(ii)
	case uint16:
		return int8(ii)
	case uint32:
		return int8(ii)
	case uint64:
		return int8(ii)
	case float32:
		return int8(ii)
	case float64:
		return int8(ii)
	case string:
		iv, err := strconv.ParseInt(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return int8(iv)
	default:
		panic("unknown type")
	}
}

//template format
var __formatVTypeToInt8Uint = func(i interface{}) uint {
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

// add 6,6 success

// add 7, 7 failed

// add 7, 7 success
