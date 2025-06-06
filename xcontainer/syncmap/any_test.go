package syncmap

import (
	"errors"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"

	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestSyncMap(t *testing.T) {
	Convey("test sync map", t, func() {
		for _, tr := range []*SyncMap[int16, uint64]{New[int16, uint64]()} {
			So(tr.Len(), ShouldBeZeroValue)
			var k, v = int16(3), uint64(4)
			So(tr.Len(), ShouldEqual, 0)
			tr.Store(k, v)
			v1, ok := tr.Load(k)
			So(ok, ShouldBeTrue)
			So(v1, ShouldEqual, v)

			So(tr.Keys(), ShouldResemble, []int16{int16(3)})
			So(tr.Get(int16(3)), ShouldEqual, uint64(4))
			So(tr.Contains(int16(3)), ShouldBeTrue)

			tr.Store(int16(4), uint64(5))
			tr.Store(int16(5), uint64(6))
			ol := tr.Len()
			tr.DeleteMultiple(int16(4), int16(5))
			So(tr.Len(), ShouldEqual, ol-2)

			ol = tr.Len()
			tr.Store(int16(4), uint64(5))
			tr.Store(int16(5), uint64(6))
			vl, ok := tr.LoadAndDelete(int16(4))
			So(vl, ShouldEqual, uint64(5))
			So(ok, ShouldBeTrue)
			So(tr.Len(), ShouldEqual, ol+1)

			tr.Store(int16(4), uint64(5))
			fge := []func(key int16, cf func(key int16) (uint64, error)) (value uint64, loaded bool, err error){tr.GetOrSetFuncErrorLock}
			defv, defv2 := uint64(6), uint64(7)
			for _, f := range fge {
				v, l, e := f(int16(6), func(key int16) (uint64, error) {
					return defv, nil
				})
				So(v, ShouldEqual, defv)
				So(l, ShouldBeFalse)
				So(e, ShouldBeNil)

				v, l, e = f(int16(7), func(key int16) (uint64, error) {
					return defv2, errors.New("")
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
				So(e, ShouldNotBeNil)
			}
			fg := []func(key int16, cf func(key int16) uint64) (value uint64, loaded bool){tr.GetOrSetFuncLock}
			for _, f := range fg {
				v, l := f(int16(7), func(key int16) uint64 {
					return defv2
				})
				So(v, ShouldEqual, defv2)
				So(l, ShouldBeFalse)
			}

			v, ok = tr.LoadOrStore(int16(8), uint64(9))
			So(v, ShouldEqual, uint64(9))
			So(ok, ShouldBeFalse)

			So(func() {
				tr.Range(func(key int16, value uint64) bool {
					return true
				})
			}, ShouldNotPanic)

		}
	})
}

func TestMap(t *testing.T) {
	Convey("SyncMap should work ok", t, func() {
		var sm SyncMap[int, string]

		v, ok := sm.Load(1)
		So(ok, ShouldBeFalse)
		So(v, ShouldBeEmpty)

		v, ok = sm.LoadOrStore(1, "1")
		So(ok, ShouldBeFalse)
		So(v, ShouldEqual, "1")

		v, ok = sm.LoadOrStore(1, "2")
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "1")

		v, ok = sm.Load(1)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "1")

		sm.Delete(1)
		v, ok = sm.LoadAndDelete(1)
		So(ok, ShouldBeFalse)
		So(v, ShouldBeEmpty)

		v, ok = sm.Swap(1, "1")
		So(ok, ShouldBeFalse)
		So(v, ShouldBeEmpty)

		sm.Store(1, "1")

		v, ok = sm.Swap(1, "2")
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "1")

		ok = sm.CompareAndDelete(1, "1")
		So(ok, ShouldBeFalse)

		ok = sm.CompareAndDelete(1, "2")
		So(ok, ShouldBeTrue)
	})
}

func TestConcurrentRange(t *testing.T) {
	const mapSize = 1 << 10

	m := new(SyncMap[int64, int64])
	for n := int64(1); n <= mapSize; n++ {
		m.Store(n, int64(n))
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		close(done)
		wg.Wait()
	}()
	for g := int64(runtime.GOMAXPROCS(0)); g > 0; g-- {
		r := rand.New(rand.NewSource(g))
		wg.Add(1)
		go func(g int64) {
			defer wg.Done()
			for i := int64(0); ; i++ {
				select {
				case <-done:
					return
				default:
				}
				for n := int64(1); n < mapSize; n++ {
					if r.Int63n(mapSize) == 0 {
						m.Store(n, n*i*g)
					} else {
						m.Load(n)
					}
				}
			}
		}(g)
	}

	iters := 1 << 10
	if testing.Short() {
		iters = 16
	}
	for n := iters; n > 0; n-- {
		seen := make(map[int64]bool, mapSize)

		m.Range(func(k, v int64) bool {
			if v%k != 0 {
				t.Fatalf("while Storing multiples of %v, Range saw value %v", k, v)
			}
			if seen[k] {
				t.Fatalf("Range visited key %v twice", k)
			}
			seen[k] = true
			return true
		})

		if len(seen) != mapSize {
			t.Fatalf("Range visited %v elements of %v-element Map", len(seen), mapSize)
		}
	}
}

func TestIssue40999(t *testing.T) {
	var m SyncMap[*int, struct{}]

	// Since the miss-counting in missLocked (via Delete)
	// compares the miss count with len(m.dirty),
	// add an initial entry to bias len(m.dirty) above the miss count.
	m.Store(nil, struct{}{})

	var finalized uint32

	// Set finalizers that count for collected keys. A non-zero count
	// indicates that keys have not been leaked.
	for atomic.LoadUint32(&finalized) == 0 {
		p := new(int)
		runtime.SetFinalizer(p, func(*int) {
			atomic.AddUint32(&finalized, 1)
		})
		m.Store(p, struct{}{})
		m.Delete(p)
		runtime.GC()
	}
}

func TestMapRangeNestedCall(t *testing.T) { // Issue 46399
	var m SyncMap[int, string]
	for i, v := range [3]string{"hello", "world", "Go"} {
		m.Store(i, v)
	}
	m.Range(func(key int, value string) bool {
		m.Range(func(key int, value string) bool {
			// We should be able to load the key offered in the Range callback,
			// because there are no concurrent Delete involved in this tested map.
			if v, ok := m.Load(key); !ok || !reflect.DeepEqual(v, value) {
				t.Fatalf("Nested Range loads unexpected value, got %+v want %+v", v, value)
			}

			// We didn't keep 42 and a value into the map before, if somehow we loaded
			// a value from such a key, meaning there must be an internal bug regarding
			// nested range in the Map.
			if _, loaded := m.LoadOrStore(42, "dummy"); loaded {
				t.Fatalf("Nested Range loads unexpected value, want store a new value")
			}

			// Try to Store then LoadAndDelete the corresponding value with the key
			// 42 to the Map. In this case, the key 42 and associated value should be
			// removed from the Map. Therefore any future range won't observe key 42
			// as we checked in above.
			val := "SyncMap"
			m.Store(42, val)
			if v, loaded := m.LoadAndDelete(42); !loaded || !reflect.DeepEqual(v, val) {
				t.Fatalf("Nested Range loads unexpected value, got %v, want %v", v, val)
			}
			return true
		})

		// Remove key from Map on-the-fly.
		m.Delete(key)
		return true
	})

	// After a Range of Delete, all keys should be removed and any
	// further Range won't invoke the callback. Hence length remains 0.
	length := 0
	m.Range(func(key int, value string) bool {
		length++
		return true
	})

	if length != 0 {
		t.Fatalf("Unexpected SyncMap size, got %v want %v", length, 0)
	}
}

func TestCompareAndSwap_NonExistingKey(t *testing.T) {
	m := new(SyncMap[int64, *int64])
	v := int64(42)
	if m.CompareAndSwap(0, nil, &v) {
		// See https://go.dev/issue/51972#issuecomment-1126408637.
		t.Fatalf("CompareAndSwap on a non-existing key succeeded")
	}
}
