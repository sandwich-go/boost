// Code generated by gotemplate. DO NOT EDIT.

package syncmap

import (
	"sort"
	"strconv"
	"sync"
)

//template type SyncMap(KType,VType)

type Int8Uint struct {
	sm     sync.Map
	locker sync.RWMutex
}

func NewInt8Uint() *Int8Uint {
	return &Int8Uint{}
}

func (s *Int8Uint) Keys() (ret []int8) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(int8))
		return true
	})
	return ret
}

func (s *Int8Uint) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

func (s *Int8Uint) Contains(key int8) (ok bool) {
	_, ok = s.Load(key)
	return
}

func (s *Int8Uint) Get(key int8) (value uint) {
	value, _ = s.Load(key)
	return
}

func (s *Int8Uint) Load(key int8) (value uint, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(uint), true
	}
	return
}
func (s *Int8Uint) DeleteMultiple(keys ...int8) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}
func (s *Int8Uint) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}
func (s *Int8Uint) Delete(key int8)          { s.sm.Delete(key) }
func (s *Int8Uint) Store(key int8, val uint) { s.sm.Store(key, val) }
func (s *Int8Uint) LoadAndDelete(key int8) (value uint, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(uint), true
	}
	return
}
func (s *Int8Uint) GetOrSetFuncErrorLock(key int8, cf func(key int8) (uint, error)) (value uint, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

func (s *Int8Uint) LoadOrStoreFuncErrorLock(key int8, cf func(key int8) (uint, error)) (value uint, loaded bool, err error) {
	if v, ok := s.Load(key); ok {
		return v, true, nil
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	// 再次重试，如果获取到则直接返回
	if v, ok := s.Load(key); ok {
		return v, true, nil
	}
	value, err = cf(key)
	if err != nil {
		return value, false, err
	}
	s.Store(key, value)
	return value, false, nil
}

func (s *Int8Uint) GetOrSetFuncLock(key int8, cf func(key int8) uint) (value uint, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

func (s *Int8Uint) LoadOrStoreFuncLock(key int8, cf func(key int8) uint) (value uint, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key int8) (uint, error) {
		return cf(key), nil
	})
	return value, loaded
}

func (s *Int8Uint) LoadOrStore(key int8, val uint) (uint, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(uint), ok
}

func (s *Int8Uint) Range(f func(key int8, value uint) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(int8), v.(uint))
	})
}

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
