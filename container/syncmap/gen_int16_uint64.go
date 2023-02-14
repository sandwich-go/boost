// Code generated by gotemplate. DO NOT EDIT.

package syncmap

import (
	"sort"
	"strconv"
	"sync"
)

//template type SyncMap(KType,VType)

type Int16Uint64 struct {
	sm     sync.Map
	locker sync.RWMutex
}

func NewInt16Uint64() *Int16Uint64 {
	return &Int16Uint64{}
}

func (s *Int16Uint64) Keys() (ret []int16) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(int16))
		return true
	})
	return ret
}

func (s *Int16Uint64) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

func (s *Int16Uint64) Contains(key int16) (ok bool) {
	_, ok = s.Load(key)
	return
}

func (s *Int16Uint64) Get(key int16) (value uint64) {
	value, _ = s.Load(key)
	return
}

func (s *Int16Uint64) Load(key int16) (value uint64, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(uint64), true
	}
	return
}
func (s *Int16Uint64) DeleteMultiple(keys ...int16) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}
func (s *Int16Uint64) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}
func (s *Int16Uint64) Delete(key int16)            { s.sm.Delete(key) }
func (s *Int16Uint64) Store(key int16, val uint64) { s.sm.Store(key, val) }
func (s *Int16Uint64) LoadAndDelete(key int16) (value uint64, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(uint64), true
	}
	return
}
func (s *Int16Uint64) GetOrSetFuncErrorLock(key int16, cf func(key int16) (uint64, error)) (value uint64, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

func (s *Int16Uint64) LoadOrStoreFuncErrorLock(key int16, cf func(key int16) (uint64, error)) (value uint64, loaded bool, err error) {
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

func (s *Int16Uint64) GetOrSetFuncLock(key int16, cf func(key int16) uint64) (value uint64, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

func (s *Int16Uint64) LoadOrStoreFuncLock(key int16, cf func(key int16) uint64) (value uint64, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key int16) (uint64, error) {
		return cf(key), nil
	})
	return value, loaded
}

func (s *Int16Uint64) LoadOrStore(key int16, val uint64) (uint64, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(uint64), ok
}

func (s *Int16Uint64) Range(f func(key int16, value uint64) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(int16), v.(uint64))
	})
}

func (s *Int16Uint64) RangeDeterministic(f func(key int16, value uint64) bool, sortableGetter func([]int16) sort.Interface) {
	var keys []int16
	s.sm.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(int16))
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
var __formatKTypeToInt16Uint64 = func(i interface{}) int16 {
	switch ii := i.(type) {
	case int:
		return int16(ii)
	case int8:
		return int16(ii)
	case int16:
		return int16(ii)
	case int32:
		return int16(ii)
	case int64:
		return int16(ii)
	case uint:
		return int16(ii)
	case uint8:
		return int16(ii)
	case uint16:
		return int16(ii)
	case uint32:
		return int16(ii)
	case uint64:
		return int16(ii)
	case float32:
		return int16(ii)
	case float64:
		return int16(ii)
	case string:
		iv, err := strconv.ParseInt(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return int16(iv)
	default:
		panic("unknown type")
	}
}

//template format
var __formatVTypeToInt16Uint64 = func(i interface{}) uint64 {
	switch ii := i.(type) {
	case int:
		return uint64(ii)
	case int8:
		return uint64(ii)
	case int16:
		return uint64(ii)
	case int32:
		return uint64(ii)
	case int64:
		return uint64(ii)
	case uint:
		return uint64(ii)
	case uint8:
		return uint64(ii)
	case uint16:
		return uint64(ii)
	case uint32:
		return uint64(ii)
	case uint64:
		return uint64(ii)
	case float32:
		return uint64(ii)
	case float64:
		return uint64(ii)
	case string:
		iv, err := strconv.ParseUint(ii, 10, 64)
		if err != nil {
			panic(err)
		}
		return uint64(iv)
	default:
		panic("unknown type")
	}
}
