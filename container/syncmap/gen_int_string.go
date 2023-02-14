// Code generated by gotemplate. DO NOT EDIT.

package syncmap

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

//template type SyncMap(KType,VType)

type IntString struct {
	sm     sync.Map
	locker sync.RWMutex
}

func NewIntString() *IntString {
	return &IntString{}
}

func (s *IntString) Keys() (ret []int) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(int))
		return true
	})
	return ret
}

func (s *IntString) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

func (s *IntString) Contains(key int) (ok bool) {
	_, ok = s.Load(key)
	return
}

func (s *IntString) Get(key int) (value string) {
	value, _ = s.Load(key)
	return
}

func (s *IntString) Load(key int) (value string, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(string), true
	}
	return
}
func (s *IntString) DeleteMultiple(keys ...int) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}
func (s *IntString) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}
func (s *IntString) Delete(key int)            { s.sm.Delete(key) }
func (s *IntString) Store(key int, val string) { s.sm.Store(key, val) }
func (s *IntString) LoadAndDelete(key int) (value string, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(string), true
	}
	return
}
func (s *IntString) GetOrSetFuncErrorLock(key int, cf func(key int) (string, error)) (value string, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

func (s *IntString) LoadOrStoreFuncErrorLock(key int, cf func(key int) (string, error)) (value string, loaded bool, err error) {
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

func (s *IntString) GetOrSetFuncLock(key int, cf func(key int) string) (value string, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

func (s *IntString) LoadOrStoreFuncLock(key int, cf func(key int) string) (value string, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key int) (string, error) {
		return cf(key), nil
	})
	return value, loaded
}

func (s *IntString) LoadOrStore(key int, val string) (string, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(string), ok
}

func (s *IntString) Range(f func(key int, value string) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(int), v.(string))
	})
}

func (s *IntString) RangeDeterministic(f func(key int, value string) bool, sortableGetter func([]int) sort.Interface) {
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
var __formatKTypeToIntString = func(i interface{}) int {
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
var __formatVTypeToIntString = func(i interface{}) string {
	switch ii := i.(type) {
	case string:
		return ii
	default:
		return fmt.Sprintf("%d", i)
	}
}
