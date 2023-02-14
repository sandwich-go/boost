package syncmap

import (
	"sort"
	"sync"
)

//template type SyncMap(KType,VType)

type KType string
type VType interface{}

type SyncMap struct {
	sm     sync.Map
	locker sync.RWMutex
}

func NewSyncMap() *SyncMap {
	return &SyncMap{}
}

func (s *SyncMap) Keys() (ret []KType) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(KType))
		return true
	})
	return ret
}

func (s *SyncMap) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

func (s *SyncMap) Contains(key KType) (ok bool) {
	_, ok = s.Load(key)
	return
}

func (s *SyncMap) Get(key KType) (value VType) {
	value, _ = s.Load(key)
	return
}

func (s *SyncMap) Load(key KType) (value VType, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(VType), true
	}
	return
}
func (s *SyncMap) DeleteMultiple(keys ...KType) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}
func (s *SyncMap) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}
func (s *SyncMap) Delete(key KType)           { s.sm.Delete(key) }
func (s *SyncMap) Store(key KType, val VType) { s.sm.Store(key, val) }
func (s *SyncMap) LoadAndDelete(key KType) (value VType, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(VType), true
	}
	return
}
func (s *SyncMap) GetOrSetFuncErrorLock(key KType, cf func(key KType) (VType, error)) (value VType, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

func (s *SyncMap) LoadOrStoreFuncErrorLock(key KType, cf func(key KType) (VType, error)) (value VType, loaded bool, err error) {
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

func (s *SyncMap) GetOrSetFuncLock(key KType, cf func(key KType) VType) (value VType, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

func (s *SyncMap) LoadOrStoreFuncLock(key KType, cf func(key KType) VType) (value VType, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key KType) (VType, error) {
		return cf(key), nil
	})
	return value, loaded
}

func (s *SyncMap) LoadOrStore(key KType, val VType) (VType, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(VType), ok
}

func (s *SyncMap) Range(f func(key KType, value VType) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(KType), v.(VType))
	})
}

func (s *SyncMap) RangeDeterministic(f func(key KType, value VType) bool, sortableGetter func([]KType) sort.Interface) {
	var keys []KType
	s.sm.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(KType))
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
