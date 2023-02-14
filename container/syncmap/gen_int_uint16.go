// Code generated by gotemplate. DO NOT EDIT.

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

func NewIntUint16() *IntUint16 {
	return &IntUint16{}
}

func (s *IntUint16) Keys() (ret []int) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(int))
		return true
	})
	return ret
}

func (s *IntUint16) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

func (s *IntUint16) Contains(key int) (ok bool) {
	_, ok = s.Load(key)
	return
}

func (s *IntUint16) Get(key int) (value uint16) {
	value, _ = s.Load(key)
	return
}

func (s *IntUint16) Load(key int) (value uint16, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(uint16), true
	}
	return
}
func (s *IntUint16) DeleteMultiple(keys ...int) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}
func (s *IntUint16) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}
func (s *IntUint16) Delete(key int)            { s.sm.Delete(key) }
func (s *IntUint16) Store(key int, val uint16) { s.sm.Store(key, val) }
func (s *IntUint16) LoadAndDelete(key int) (value uint16, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(uint16), true
	}
	return
}
func (s *IntUint16) GetOrSetFuncErrorLock(key int, cf func(key int) (uint16, error)) (value uint16, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

func (s *IntUint16) LoadOrStoreFuncErrorLock(key int, cf func(key int) (uint16, error)) (value uint16, loaded bool, err error) {
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

func (s *IntUint16) GetOrSetFuncLock(key int, cf func(key int) uint16) (value uint16, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

func (s *IntUint16) LoadOrStoreFuncLock(key int, cf func(key int) uint16) (value uint16, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key int) (uint16, error) {
		return cf(key), nil
	})
	return value, loaded
}

func (s *IntUint16) LoadOrStore(key int, val uint16) (uint16, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(uint16), ok
}

func (s *IntUint16) Range(f func(key int, value uint16) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(int), v.(uint16))
	})
}

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
