// Code generated by gotemplate. DO NOT EDIT.

// syncmap 提供了一个同步的映射实现，允许安全并发的访问
package syncmap

import (
	"sort"
	"strconv"
	"sync"
)

//template type SyncMap(KType,VType)

type IntUint16 struct {
	sm     sync.Map
	locker sync.RWMutex
}

// NewSyncMap 构造函数，返回一个新的 SyncMap
func NewIntUint16() *IntUint16 {
	return &IntUint16{}
}

// Keys 获取映射中的所有键，返回一个 Key 类型的切片
func (s *IntUint16) Keys() (ret []int) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(int))
		return true
	})
	return ret
}

// Len 获取映射中键值对的数量
func (s *IntUint16) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

// Contains 检查映射中是否包含指定键
func (s *IntUint16) Contains(key int) (ok bool) {
	_, ok = s.Load(key)
	return
}

// Get 获取映射中的值
func (s *IntUint16) Get(key int) (value uint16) {
	value, _ = s.Load(key)
	return
}

// Load 获取映射中的值和是否成功加载的标志
func (s *IntUint16) Load(key int) (value uint16, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(uint16), true
	}
	return
}

// DeleteMultiple 删除映射中的多个键
func (s *IntUint16) DeleteMultiple(keys ...int) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}

// Clear 清空映射
func (s *IntUint16) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}

// Delete 删除映射中的值
func (s *IntUint16) Delete(key int) { s.sm.Delete(key) }

// Store 往映射中存储一个键值对
func (s *IntUint16) Store(key int, val uint16) { s.sm.Store(key, val) }

// LoadAndDelete 获取映射中的值，并将其从映射中删除
func (s *IntUint16) LoadAndDelete(key int) (value uint16, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(uint16), true
	}
	return
}

// GetOrSetFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *IntUint16) GetOrSetFuncErrorLock(key int, cf func(key int) (uint16, error)) (value uint16, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

// LoadOrStoreFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *IntUint16) LoadOrStoreFuncErrorLock(key int, cf func(key int) (uint16, error)) (value uint16, loaded bool, err error) {
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
func (s *IntUint16) GetOrSetFuncLock(key int, cf func(key int) uint16) (value uint16, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

// LoadOrStoreFuncLock 根据key获取对应的value，若不存在则通过cf回调创建value并存储
func (s *IntUint16) LoadOrStoreFuncLock(key int, cf func(key int) uint16) (value uint16, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key int) (uint16, error) {
		return cf(key), nil
	})
	return value, loaded
}

// LoadOrStore 存储一个 key-value 对，若key已存在则返回已存在的value
func (s *IntUint16) LoadOrStore(key int, val uint16) (uint16, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(uint16), ok
}

// Range 遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数f
func (s *IntUint16) Range(f func(key int, value uint16) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(int), v.(uint16))
	})
}

// RangeDeterministic 按照 key 的顺序遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数 f, f返回false则中断退出
// 参数 sortableGetter 接收一个 KType 切片并返回一个可排序接口，用于对key进行排序
func (s *IntUint16) RangeDeterministic(f func(key int, value uint16) bool, sortableGetter func([]int) sort.Interface) {
	var keys []int
	s.sm.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(int))
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
var __formatKTypeToIntUint16 = func(i interface{}) int {
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
var __formatVTypeToIntUint16 = func(i interface{}) uint16 {
	switch ii := i.(type) {
	case int:
		return uint16(ii)
	case int8:
		return uint16(ii)
	case int16:
		return uint16(ii)
	case int32:
		return uint16(ii)
	case int64:
		return uint16(ii)
	case uint:
		return uint16(ii)
	case uint8:
		return uint16(ii)
	case uint16:
		return uint16(ii)
	case uint32:
		return uint16(ii)
	case uint64:
		return uint16(ii)
	case float32:
		return uint16(ii)
	case float64:
		return uint16(ii)
	case string:
		iv, err := strconv.ParseUint(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return uint16(iv)
	default:
		panic("unknown type")
	}
}

// add 6,6 success

// add 7, 7 failed

// add 7, 7 success
