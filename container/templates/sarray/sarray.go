package sarray

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"sort"
	"sync"
	"testing"
)

//template type SyncArray(VType)

type VType interface{}

type SyncArray struct {
	mu    *localRWMutexVType
	array []VType
}

// New 创建非协程安全版本
func New() *SyncArray { return newWithSafe(false) }

// NewSync 创建协程安全版本
func NewSync() *SyncArray { return newWithSafe(true) }

func newWithSafe(safe bool) *SyncArray {
	return &SyncArray{
		mu:    newLocalRWMutexVType(safe),
		array: make([]VType, 0),
	}
}

// At 返回指定位置元素，如果越界则返回默认空值
func (a *SyncArray) At(index int) (value VType) {
	value, _ = a.Get(index)
	return
}

// Get 返回指定位置元素，found标识元素是否存在
func (a *SyncArray) Get(index int) (value VType, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.array) {
		return
	}
	return a.array[index], true
}

func (a *SyncArray) errorIndexOutRangeUnLock(index int) error {
	return fmt.Errorf("index %d out of array range %d", index, len(a.array))
}

// Set 设定指定位置数据
func (a *SyncArray) Set(index int, value VType) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return a.errorIndexOutRangeUnLock(index)
	}
	a.array[index] = value
	return nil
}

// SetArray 替换底层存储
func (a *SyncArray) SetArray(array []VType) *SyncArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = array
	return a
}

// Replace 替换指定位置元素
func (a *SyncArray) Replace(given []VType) *SyncArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	max := len(given)
	if max > len(a.array) {
		max = len(a.array)
	}
	for i := 0; i < max; i++ {
		a.array[i] = given[i]
	}
	return a
}

// SortFunc  根据指定的方法进行排序
func (a *SyncArray) SortFunc(less func(v1, v2 VType) bool) *SyncArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.array, func(i, j int) bool {
		return less(a.array[i], a.array[j])
	})
	return a
}

// InsertBefore 在index位置前插入数据
func (a *SyncArray) InsertBefore(index int, value VType) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return a.errorIndexOutRangeUnLock(index)
	}
	rear := append([]VType{}, a.array[index:]...)
	a.array = append(a.array[0:index], value)
	a.array = append(a.array, rear...)
	return nil
}

// InsertAfter 在index位置后插入数据
func (a *SyncArray) InsertAfter(index int, value VType) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return a.errorIndexOutRangeUnLock(index)
	}
	rear := append([]VType{}, a.array[index+1:]...)
	a.array = append(a.array[0:index+1], value)
	a.array = append(a.array, rear...)
	return nil
}

// Contains  是否存在value
func (a *SyncArray) Contains(value VType) bool {
	return a.Search(value) != -1
}

// Search 查找元素，不存在返回-1
func (a *SyncArray) Search(value VType) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return -1
	}
	result := -1
	for index, v := range a.array {
		if v == value {
			result = index
			break
		}
	}
	return result
}

// DeleteValue 查找并删除找到的第一个元素，不存在返回false
func (a *SyncArray) DeleteValue(value VType) (found bool) {
	if i := a.Search(value); i != -1 {
		_, found = a.LoadAndDelete(i)
		return found
	}
	return false
}

// LoadAndDelete 删除元素，如果删除成功返回被删除的元素
func (a *SyncArray) LoadAndDelete(index int) (value VType, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doDeleteWithoutLock(index)
}

// doRemoveWithoutLock 不加锁移除元素
func (a *SyncArray) doDeleteWithoutLock(index int) (value VType, found bool) {
	if index < 0 || index >= len(a.array) {
		return
	}
	if index == 0 {
		value = a.array[0]
		a.array = a.array[1:]
		return value, true
	} else if index == len(a.array)-1 {
		value = a.array[index]
		a.array = a.array[:index]
		return value, true
	}
	value = a.array[index]
	a.array = append(a.array[:index], a.array[index+1:]...)
	return value, true
}

// PushLeft 头插入
func (a *SyncArray) PushLeft(value ...VType) *SyncArray {
	a.mu.Lock()
	a.array = append(value, a.array...)
	a.mu.Unlock()
	return a
}

// PushRight 尾插入
func (a *SyncArray) PushRight(value ...VType) *SyncArray {
	a.mu.Lock()
	a.array = append(a.array, value...)
	a.mu.Unlock()
	return a
}

// PopLeft 头弹出
func (a *SyncArray) PopLeft() (value VType, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		return
	}
	value = a.array[0]
	a.array = a.array[1:]
	return value, true
}

// PopRight 尾弹出
func (a *SyncArray) PopRight() (value VType, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	index := len(a.array) - 1
	if index < 0 {
		return
	}
	value = a.array[index]
	a.array = a.array[:index]
	return value, true
}

// PopRand 随机弹出
func (a *SyncArray) PopRand() (value VType, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doDeleteWithoutLock(rand.Intn(len(a.array)))
}

// PopRands 随机n个元素并弹出，如果size大于数组尺寸则全部弹出
func (a *SyncArray) PopRands(size int) []VType {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	if size >= len(a.array) {
		size = len(a.array)
	}
	array := make([]VType, size)
	for i := 0; i < size; i++ {
		array[i], _ = a.doDeleteWithoutLock(rand.Intn(len(a.array)))
	}
	return array
}

// Append 尾添加元素 alias of PushRight
func (a *SyncArray) Append(value ...VType) *SyncArray { return a.PushRight(value...) }

// Len 获取长度
func (a *SyncArray) Len() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// Slice 获取底层数据存储，如果为sync安全模式则返回一份拷贝，否则直接返回底层数据指针
func (a *SyncArray) Slice() []VType {
	array := ([]VType)(nil)
	if a.mu.IsSafe() {
		a.mu.RLock()
		defer a.mu.RUnlock()
		array = make([]VType, len(a.array))
		copy(array, a.array)
	} else {
		array = a.array
	}
	return array
}

// Clear 清空存储
func (a *SyncArray) Clear() *SyncArray {
	a.mu.Lock()
	if len(a.array) > 0 {
		a.array = make([]VType, 0)
	}
	a.mu.Unlock()
	return a
}

// LockFunc 写锁操作array
func (a *SyncArray) LockFunc(f func(array []VType)) *SyncArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	f(a.array)
	return a
}

// RLockFunc 读锁操作array
func (a *SyncArray) RLockFunc(f func(array []VType)) *SyncArray {
	a.mu.RLock()
	defer a.mu.RUnlock()
	f(a.array)
	return a
}

// Rand 随机一个元素
func (a *SyncArray) Rand() (value VType, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return
	}
	return a.array[rand.Intn(len(a.array))], true
}

// WalkAsc 按照index从小到大的顺序进行遍历，并将k,v作为参数执行f。如果f执行返回false则中止
func (a *SyncArray) WalkAsc(f func(k int, v VType) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for k, v := range a.array {
		if !f(k, v) {
			break
		}
	}
}

// WalkDesc 按照index从大到小的顺序进行遍历，并将k,v作为参数执行f。如果f执行返回false则中止
func (a *SyncArray) WalkDesc(f func(k int, v VType) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for i := len(a.array) - 1; i >= 0; i-- {
		if !f(i, a.array[i]) {
			break
		}
	}
}

// MarshalJSON 序列化到json
func (a SyncArray) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a.array)
}

// UnmarshalJSON 由json反序列化
func (a *SyncArray) UnmarshalJSON(b []byte) error {
	if a.array == nil {
		a.array = make([]VType, 0)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := json.Unmarshal(b, &a.array); err != nil {
		return err
	}
	return nil
}

// Empty 是否为空
func (a *SyncArray) Empty() bool { return a.Len() == 0 }

type localRWMutexVType struct {
	*sync.RWMutex
}

func newLocalRWMutexVType(safe bool) *localRWMutexVType {
	mu := localRWMutexVType{}
	if safe {
		mu.RWMutex = new(sync.RWMutex)
	}
	return &mu
}

func (mu *localRWMutexVType) IsSafe() bool {
	return mu.RWMutex != nil
}

func (mu *localRWMutexVType) Lock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Lock()
	}
}

func (mu *localRWMutexVType) Unlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.Unlock()
	}
}

func (mu *localRWMutexVType) RLock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RLock()
	}
}

func (mu *localRWMutexVType) RUnlock() {
	if mu.RWMutex != nil {
		mu.RWMutex.RUnlock()
	}
}

//template format
var __formatTo func(interface{}) VType

func TestSyncArray(t *testing.T) {
	Convey("test sync array", t, func() {
		for _, tr := range []*SyncArray{New(), NewSync()} {
			So(tr.Len(), ShouldBeZeroValue)
			_, exists := tr.Get(0)
			So(exists, ShouldBeFalse)
			So(tr.Empty(), ShouldBeTrue)
			var e0 = __formatTo(3)
			tr.PushLeft(e0)
			So(tr.Len(), ShouldEqual, 1)

			e := tr.At(0)
			So(e, ShouldEqual, e0)
			v, f := tr.Get(0)
			So(f, ShouldBeTrue)
			So(v, ShouldEqual, e)
			_, f = tr.Get(1)
			So(f, ShouldBeFalse)

			tr.SetArray([]VType{__formatTo(1), __formatTo(1), __formatTo(1), __formatTo(1), __formatTo(1)})
			So(tr.Len(), ShouldEqual, 5)
			tr.SetArray([]VType{__formatTo(1), __formatTo(2), __formatTo(3), __formatTo(4), __formatTo(5)})
			So(tr.Len(), ShouldEqual, 5)

			rpv := __formatTo(20)
			err := tr.Set(2, rpv)
			So(err, ShouldBeEmpty)
			So(tr.At(2), ShouldEqual, rpv)

			tr.Replace([]VType{__formatTo(1), __formatTo(1)})
			So(tr.Len(), ShouldEqual, 5)
			So(tr.At(0), ShouldEqual, tr.At(1))
			So(tr.At(2), ShouldNotEqual, tr.At(0))

			iv1 := __formatTo(11)
			err = tr.InsertBefore(0, iv1)
			So(err, ShouldBeNil)
			So(tr.At(0), ShouldEqual, iv1)

			iv2 := __formatTo(12)
			err = tr.InsertAfter(0, iv2)
			So(err, ShouldBeNil)
			So(tr.At(1), ShouldEqual, iv2)

			So(tr.Contains(iv1), ShouldBeTrue)
			So(tr.Search(iv1), ShouldNotEqual, -1)

			So(tr.DeleteValue(iv2), ShouldBeTrue)
			v, f = tr.LoadAndDelete(0)
			So(f, ShouldBeTrue)
			So(v, ShouldEqual, iv1)

			pl := __formatTo(11)
			tr.PushLeft(pl)
			So(tr.At(0), ShouldEqual, pl)
			pr := __formatTo(21)
			tr.PushRight(pr)
			So(tr.At(tr.Len()-1), ShouldEqual, pr)

			v, f = tr.PopLeft()
			So(v, ShouldEqual, pl)
			v, f = tr.PopRight()
			So(v, ShouldEqual, pr)
			l := tr.Len()
			_, f = tr.PopRand()
			So(f, ShouldBeTrue)
			So(tr.Len()+1, ShouldEqual, l)
			l = tr.Len()
			poplen := 2
			pv := tr.PopRands(poplen)
			So(len(pv), ShouldEqual, poplen)
			So(tr.Len(), ShouldBeGreaterThanOrEqualTo, l-poplen)

			aps := []VType{__formatTo(35), __formatTo(40), __formatTo(45), __formatTo(50)}
			tr.Append(aps...)
			for i := len(aps); i > 0; i-- {
				So(aps[i-1], ShouldEqual, func() VType { v, f = tr.PopRight(); So(f, ShouldBeTrue); return v }())
			}

			tr.Clear()
			So(tr.Len(), ShouldEqual, 0)

			tr.Append(aps...)
			s := tr.Slice()
			So(len(s), ShouldEqual, tr.Len())

			k := 0
			tr.WalkAsc(func(key int, val VType) bool {
				So(key, ShouldEqual, k)
				So(val, ShouldEqual, s[k])
				k++
				return true
			})

			k = len(s) - 1
			tr.WalkDesc(func(key int, val VType) bool {
				So(key, ShouldEqual, k)
				So(val, ShouldEqual, s[k])
				k--
				return true
			})

			So(func() {
				tr.LockFunc(func(array []VType) {
					return
				})
			}, ShouldNotPanic)

			So(func() {
				tr.RLockFunc(func(array []VType) {
					return
				})
			}, ShouldNotPanic)

			j, err := tr.MarshalJSON()
			So(err, ShouldBeNil)
			err = tr.UnmarshalJSON(j)
			So(err, ShouldBeNil)
		}
	})
}
