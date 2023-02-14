// Code generated by gotemplate. DO NOT EDIT.

package syncmap

import (
	"sort"
	"strconv"
	"sync"
)

//template type SyncMap(KType,VType)

type Int8Int16 struct {
	sm     sync.Map
	locker sync.RWMutex
}

func NewInt8Int16() *Int8Int16 {
	return &Int8Int16{}
}

func (s *Int8Int16) Keys() (ret []int8) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(int8))
		return true
	})
	return ret
}

func (s *Int8Int16) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

func (s *Int8Int16) Contains(key int8) (ok bool) {
	_, ok = s.Load(key)
	return
}

func (s *Int8Int16) Get(key int8) (value int16) {
	value, _ = s.Load(key)
	return
}

func (s *Int8Int16) Load(key int8) (value int16, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(int16), true
	}
	return
}
func (s *Int8Int16) DeleteMultiple(keys ...int8) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}
func (s *Int8Int16) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}
func (s *Int8Int16) Delete(key int8)           { s.sm.Delete(key) }
func (s *Int8Int16) Store(key int8, val int16) { s.sm.Store(key, val) }
func (s *Int8Int16) LoadAndDelete(key int8) (value int16, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(int16), true
	}
	return
}
func (s *Int8Int16) GetOrSetFuncErrorLock(key int8, cf func(key int8) (int16, error)) (value int16, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

func (s *Int8Int16) LoadOrStoreFuncErrorLock(key int8, cf func(key int8) (int16, error)) (value int16, loaded bool, err error) {
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

func (s *Int8Int16) GetOrSetFuncLock(key int8, cf func(key int8) int16) (value int16, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

func (s *Int8Int16) LoadOrStoreFuncLock(key int8, cf func(key int8) int16) (value int16, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key int8) (int16, error) {
		return cf(key), nil
	})
	return value, loaded
}

func (s *Int8Int16) LoadOrStore(key int8, val int16) (int16, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(int16), ok
}

func (s *Int8Int16) Range(f func(key int8, value int16) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(int8), v.(int16))
	})
}

func (s *Int8Int16) RangeDeterministic(f func(key int8, value int16) bool, sortableGetter func([]int8) sort.Interface) {
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
var __formatKTypeToInt8Int16 = func(i interface{}) int8 {
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
var __formatVTypeToInt8Int16 = func(i interface{}) int16 {
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
