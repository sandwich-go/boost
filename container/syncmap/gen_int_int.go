// Code generated by gotemplate. DO NOT EDIT.

package syncmap

import (
	"sort"
	"strconv"
	"sync"
)

//template type SyncMap(KType,VType)

type IntInt struct {
	sm     sync.Map
	locker sync.RWMutex
}

func NewIntInt() *IntInt {
	return &IntInt{}
}

func (s *IntInt) Keys() (ret []int) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(int))
		return true
	})
	return ret
}

func (s *IntInt) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

func (s *IntInt) Contains(key int) (ok bool) {
	_, ok = s.Load(key)
	return
}

func (s *IntInt) Get(key int) (value int) {
	value, _ = s.Load(key)
	return
}

func (s *IntInt) Load(key int) (value int, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(int), true
	}
	return
}
func (s *IntInt) DeleteMultiple(keys ...int) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}
func (s *IntInt) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}
func (s *IntInt) Delete(key int)         { s.sm.Delete(key) }
func (s *IntInt) Store(key int, val int) { s.sm.Store(key, val) }
func (s *IntInt) LoadAndDelete(key int) (value int, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(int), true
	}
	return
}
func (s *IntInt) GetOrSetFuncErrorLock(key int, cf func(key int) (int, error)) (value int, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

func (s *IntInt) LoadOrStoreFuncErrorLock(key int, cf func(key int) (int, error)) (value int, loaded bool, err error) {
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

func (s *IntInt) GetOrSetFuncLock(key int, cf func(key int) int) (value int, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

func (s *IntInt) LoadOrStoreFuncLock(key int, cf func(key int) int) (value int, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key int) (int, error) {
		return cf(key), nil
	})
	return value, loaded
}

func (s *IntInt) LoadOrStore(key int, val int) (int, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(int), ok
}

func (s *IntInt) Range(f func(key int, value int) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(int), v.(int))
	})
}

func (s *IntInt) RangeDeterministic(f func(key int, value int) bool, sortableGetter func([]int) sort.Interface) {
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
var __formatKTypeToIntInt = func(i interface{}) int {
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
var __formatVTypeToIntInt = func(i interface{}) int {
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
