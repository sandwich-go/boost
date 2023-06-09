// syncmap 提供了一个同步的映射实现，允许安全并发的访问
package syncmap

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"sync"
	"testing"
)

//template type SyncMap(KType,VType)

type KType string
type VType interface{}

// SyncMap 定义并发安全的映射，使用 sync.Map 来实现
type SyncMap struct {
	sm     sync.Map
	locker sync.RWMutex
}

// NewSyncMap 构造函数，返回一个新的 SyncMap
func NewSyncMap() *SyncMap {
	return &SyncMap{}
}

// Keys 获取映射中的所有键，返回一个 Key 类型的切片
func (s *SyncMap) Keys() (ret []KType) {
	s.sm.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(KType))
		return true
	})
	return ret
}

// Len 获取映射中键值对的数量
func (s *SyncMap) Len() (c int) {
	s.sm.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

// Contains 检查映射中是否包含指定键
func (s *SyncMap) Contains(key KType) (ok bool) {
	_, ok = s.Load(key)
	return
}

// Get 获取映射中的值
func (s *SyncMap) Get(key KType) (value VType) {
	value, _ = s.Load(key)
	return
}

// Load 获取映射中的值和是否成功加载的标志
func (s *SyncMap) Load(key KType) (value VType, loaded bool) {
	if v, ok := s.sm.Load(key); ok {
		return v.(VType), true
	}
	return
}

// DeleteMultiple 删除映射中的多个键
func (s *SyncMap) DeleteMultiple(keys ...KType) {
	for _, k := range keys {
		s.sm.Delete(k)
	}
}

// Clear 清空映射
func (s *SyncMap) Clear() {
	s.sm.Range(func(key, value interface{}) bool {
		s.sm.Delete(key)
		return true
	})
}

// Delete 删除映射中的值
func (s *SyncMap) Delete(key KType) { s.sm.Delete(key) }

// Store 往映射中存储一个键值对
func (s *SyncMap) Store(key KType, val VType) { s.sm.Store(key, val) }

// LoadAndDelete 获取映射中的值，并将其从映射中删除
func (s *SyncMap) LoadAndDelete(key KType) (value VType, loaded bool) {
	if v, ok := s.sm.LoadAndDelete(key); ok {
		return v.(VType), true
	}
	return
}

// GetOrSetFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *SyncMap) GetOrSetFuncErrorLock(key KType, cf func(key KType) (VType, error)) (value VType, loaded bool, err error) {
	return s.LoadOrStoreFuncErrorLock(key, cf)
}

// LoadOrStoreFuncErrorLock 函数根据key查找值，如果key存在则返回对应的值，否则用cf函数计算得到一个新的值，存储到 SyncMap 中并返回。
// 如果执行cf函数时出错，则返回error。
// 函数内部使用读写锁实现并发安全
func (s *SyncMap) LoadOrStoreFuncErrorLock(key KType, cf func(key KType) (VType, error)) (value VType, loaded bool, err error) {
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
func (s *SyncMap) GetOrSetFuncLock(key KType, cf func(key KType) VType) (value VType, loaded bool) {
	return s.LoadOrStoreFuncLock(key, cf)
}

// LoadOrStoreFuncLock 根据key获取对应的value，若不存在则通过cf回调创建value并存储
func (s *SyncMap) LoadOrStoreFuncLock(key KType, cf func(key KType) VType) (value VType, loaded bool) {
	value, loaded, _ = s.LoadOrStoreFuncErrorLock(key, func(key KType) (VType, error) {
		return cf(key), nil
	})
	return value, loaded
}

// LoadOrStore 存储一个 key-value 对，若key已存在则返回已存在的value
func (s *SyncMap) LoadOrStore(key KType, val VType) (VType, bool) {
	actual, ok := s.sm.LoadOrStore(key, val)
	return actual.(VType), ok
}

// Range 遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数f
func (s *SyncMap) Range(f func(key KType, value VType) bool) {
	s.sm.Range(func(k, v interface{}) bool {
		return f(k.(KType), v.(VType))
	})
}

// RangeDeterministic 按照 key 的顺序遍历映射中的 key-value 对，对每个 key-value 对执行给定的函数 f, f返回false则中断退出
// 参数 sortableGetter 接收一个 KType 切片并返回一个可排序接口，用于对key进行排序
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

//template format
var __formatKTypeTo func(interface{}) KType

//template format
var __formatVTypeTo func(interface{}) VType

func TestSyncMap(t *testing.T) {
	Convey("test sync map", t, func() {
		for _, tr := range []*SyncMap{NewSyncMap()} {
			So(tr.Len(), ShouldBeZeroValue)
			var k, v = __formatKTypeTo(3), __formatVTypeTo(4)
			So(tr.Len(), ShouldEqual, 0)
			tr.Store(k, v)
			v1, ok := tr.Load(k)
			So(ok, ShouldBeTrue)
			So(v1, ShouldEqual, v)

			So(tr.Keys(), ShouldResemble, []KType{__formatKTypeTo(3)})
			So(tr.Get(__formatKTypeTo(3)), ShouldEqual, __formatVTypeTo(4))
			So(tr.Contains(__formatKTypeTo(3)), ShouldBeTrue)

			tr.Store(__formatKTypeTo(4), __formatVTypeTo(5))
			tr.Store(__formatKTypeTo(5), __formatVTypeTo(6))
			ol := tr.Len()
			tr.DeleteMultiple(__formatKTypeTo(4), __formatKTypeTo(5))
			So(tr.Len(), ShouldEqual, ol-2)

			ol = tr.Len()
			tr.Store(__formatKTypeTo(4), __formatVTypeTo(5))
			tr.Store(__formatKTypeTo(5), __formatVTypeTo(6))
			vl, ok := tr.LoadAndDelete(__formatKTypeTo(4))
			So(vl, ShouldEqual, __formatVTypeTo(5))
			So(ok, ShouldBeTrue)
			So(tr.Len(), ShouldEqual, ol+1)

			tr.Store(__formatKTypeTo(4), __formatVTypeTo(5))
			fge := []func(key KType, cf func(key KType) (VType, error)) (value VType, loaded bool, err error){tr.GetOrSetFuncErrorLock}
			defv, defv2 := __formatVTypeTo(6), __formatVTypeTo(7)
			for _, f := range fge {
				v, l, e := f(__formatKTypeTo(6), func(key KType) (VType, error) {
					return defv, nil
				})
				So(v, ShouldEqual, defv)
				So(l, ShouldBeFalse)
				So(e, ShouldBeNil)

				v, l, e = f(__formatKTypeTo(7), func(key KType) (VType, error) {
					return defv2, errors.New("")
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
				So(e, ShouldNotBeNil)
			}
			fg := []func(key KType, cf func(key KType) VType) (value VType, loaded bool){tr.GetOrSetFuncLock}
			for _, f := range fg {
				v, l := f(__formatKTypeTo(7), func(key KType) VType {
					return defv2
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
			}

			v, ok = tr.LoadOrStore(__formatKTypeTo(8), __formatVTypeTo(9))
			So(v, ShouldEqual, __formatVTypeTo(9))
			So(ok, ShouldBeFalse)

			So(func() {
				tr.Range(func(key KType, value VType) bool {
					return true
				})
			}, ShouldNotPanic)

		}
	})
}
