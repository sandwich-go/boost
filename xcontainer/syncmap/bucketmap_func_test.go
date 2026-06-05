package syncmap

import (
	"errors"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// bucketmap.go 的功能 + 并发测试。bucketmap 是生产级分桶并发 map（双缓冲
// keySnapshots + sync.Pool keySlicePool + 后台 refresh goroutine），是
// AGENTS.md §8 列出的高风险包。
//
// 测试策略：
//   1. 单线程行为契约：所有 32 个公开方法（Load/Store/Delete/Range/...）
//      在单线程下行为正确
//   2. 双缓冲 keySnapshots：cacheKey 模式下 Keys() 返回 snapshot；
//      forceFresh 立即刷新；stopFunc 关闭后台 goroutine
//   3. 并发安全：100 goroutine × 50 次混合 CRUD，race detector 守关；
//      最终状态一致性
//   4. 边界：bucketNum<=0 退化为 1；indexByKey 接收负 hash 用 0 兜底

// stringHash 简单 hash：取字符串长度。仅用于测试。
func stringHash(s string) int64 {
	return int64(len(s))
}

// TestNewBucketMap_DefaultPath 不带 cacheKey 的基础工厂。
func TestNewBucketMap_DefaultPath(t *testing.T) {
	Convey("NewBucketMap 创建无后台 refresh 的 map", t, func() {
		m := NewBucketMap[string, int](4, stringHash)
		So(m, ShouldNotBeNil)
		So(m.BucketNum(), ShouldEqual, 4)

		m.Store("hi", 1)
		v, ok := m.Load("hi")
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, 1)
	})

	Convey("bucketNum<=0 退化为 1，indexByKey 走单 bucket 早返路径", t, func() {
		m := NewBucketMap[string, int](0, stringHash)
		So(m.BucketNum(), ShouldEqual, 1)
		// 在单 bucket 上 Store/Load 走 indexByKey 的 bucketNum==1 早返分支
		m.Store("any-key", 7)
		v, ok := m.Load("any-key")
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, 7)

		m2 := NewBucketMap[string, int](-5, stringHash)
		So(m2.BucketNum(), ShouldEqual, 1)
	})

	Convey("indexByKey hash<=0 兜底返回 bucket 0", t, func() {
		// 让 hashFunc 返回负数
		m := NewBucketMap[string, int](4, func(k string) int64 {
			return -100
		})
		m.Store("foo", 42)
		So(m.Get("foo"), ShouldEqual, 42)
	})
}

// TestBucketMap_CRUD 覆盖 Load/Store/LoadOrStore/LoadAndDelete/Delete/
// Swap/CompareAndSwap/CompareAndDelete/Contains/Get 等 CRUD 类方法。
func TestBucketMap_CRUD(t *testing.T) {
	m := NewBucketMap[string, int](8, stringHash)

	Convey("Store/Load/Contains/Get", t, func() {
		m.Store("a", 100)
		v, ok := m.Load("a")
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, 100)
		So(m.Contains("a"), ShouldBeTrue)
		So(m.Get("a"), ShouldEqual, 100)

		// 不存在的 key
		v2, ok2 := m.Load("nope")
		So(ok2, ShouldBeFalse)
		So(v2, ShouldEqual, 0) // V 零值
		So(m.Contains("nope"), ShouldBeFalse)
		So(m.Get("nope"), ShouldEqual, 0)
	})

	Convey("LoadOrStore", t, func() {
		// 第一次：store，loaded=false
		actual, loaded := m.LoadOrStore("new-key", 42)
		So(loaded, ShouldBeFalse)
		So(actual, ShouldEqual, 42)
		// 第二次：load，loaded=true，返回原值（忽略入参 99）
		actual2, loaded2 := m.LoadOrStore("new-key", 99)
		So(loaded2, ShouldBeTrue)
		So(actual2, ShouldEqual, 42)
	})

	Convey("LoadAndDelete", t, func() {
		m.Store("temp", 7)
		v, loaded := m.LoadAndDelete("temp")
		So(loaded, ShouldBeTrue)
		So(v, ShouldEqual, 7)
		So(m.Contains("temp"), ShouldBeFalse)

		// 不存在 key 返回零值
		v2, loaded2 := m.LoadAndDelete("never")
		So(loaded2, ShouldBeFalse)
		So(v2, ShouldEqual, 0)
	})

	Convey("Delete / DeleteMultiple", t, func() {
		m.Store("x", 1)
		m.Store("y", 2)
		m.Store("zzz", 3)
		m.Delete("x")
		So(m.Contains("x"), ShouldBeFalse)
		So(m.Contains("y"), ShouldBeTrue)
		m.DeleteMultiple("y", "zzz", "nope")
		So(m.Contains("y"), ShouldBeFalse)
		So(m.Contains("zzz"), ShouldBeFalse)
	})

	Convey("Swap", t, func() {
		m.Store("swap-k", 10)
		prev, loaded := m.Swap("swap-k", 20)
		So(loaded, ShouldBeTrue)
		So(prev, ShouldEqual, 10)
		So(m.Get("swap-k"), ShouldEqual, 20)

		// 新 key 上 swap 返回零值 + loaded=false
		prev2, loaded2 := m.Swap("new-swap", 5)
		So(loaded2, ShouldBeFalse)
		So(prev2, ShouldEqual, 0)
		So(m.Get("new-swap"), ShouldEqual, 5)
	})

	Convey("CompareAndSwap / CompareAndDelete", t, func() {
		m.Store("cas-k", 1)
		ok := m.CompareAndSwap("cas-k", 1, 2)
		So(ok, ShouldBeTrue)
		So(m.Get("cas-k"), ShouldEqual, 2)
		// old 不匹配 swap 失败
		ok2 := m.CompareAndSwap("cas-k", 1, 99)
		So(ok2, ShouldBeFalse)
		So(m.Get("cas-k"), ShouldEqual, 2)

		// CompareAndDelete
		ok3 := m.CompareAndDelete("cas-k", 99)
		So(ok3, ShouldBeFalse)
		So(m.Contains("cas-k"), ShouldBeTrue)
		ok4 := m.CompareAndDelete("cas-k", 2)
		So(ok4, ShouldBeTrue)
		So(m.Contains("cas-k"), ShouldBeFalse)
	})

	Convey("Clear 清空所有 bucket", t, func() {
		m.Store("a", 1)
		m.Store("bb", 2)
		m.Store("ccc", 3)
		m.Clear()
		So(m.Len(), ShouldEqual, 0)
	})
}

// TestBucketMap_Range / RangeBucket / Len / Keys
func TestBucketMap_RangeAndKeys(t *testing.T) {
	Convey("Range 遍历所有 entry", t, func() {
		m := NewBucketMap[string, int](4, stringHash)
		m.Store("a", 1)
		m.Store("bb", 2)
		m.Store("ccc", 3)

		seen := make(map[string]int)
		m.Range(func(k string, v int) bool {
			seen[k] = v
			return true
		})
		So(seen, ShouldResemble, map[string]int{"a": 1, "bb": 2, "ccc": 3})

		// f 返回 false 立即停止某个 bucket 的 Range（注意 BucketMap.Range
		// 是逐 bucket 调 sync.Map.Range，每个 bucket 内能 break）
		var count int
		m.Range(func(k string, v int) bool {
			count++
			return false // 立刻停止当前 bucket
		})
		So(count, ShouldBeBetweenOrEqual, 1, 3)
	})

	Convey("Len 返回 entry 总数", t, func() {
		m := NewBucketMap[int, int](4, func(k int) int64 { return int64(k) })
		So(m.Len(), ShouldEqual, 0)
		for i := 0; i < 50; i++ {
			m.Store(i, i*10)
		}
		So(m.Len(), ShouldEqual, 50)
	})

	Convey("Keys (no cache) fallback 直接 Range 收集", t, func() {
		m := NewBucketMap[string, int](4, stringHash)
		m.Store("k1", 1)
		m.Store("k22", 2)
		keys := m.Keys()
		sort.Strings(keys)
		So(keys, ShouldResemble, []string{"k1", "k22"})
	})

	Convey("RangeBucket 只遍历指定 bucket", t, func() {
		m := NewBucketMap[string, int](4, stringHash)
		m.Store("k1", 1)
		m.Store("k22", 2)
		m.Store("k333", 3)
		// stringHash 用长度做 hash：k1(2%4=2) k22(3%4=3) k333(4%4=0)
		var bucket0 []string
		m.RangeBucket(0, func(k string, v int) bool {
			bucket0 = append(bucket0, k)
			return true
		})
		So(bucket0, ShouldContain, "k333")
	})
}

// TestBucketMap_RangeDeterministic
func TestBucketMap_RangeDeterministic(t *testing.T) {
	Convey("RangeDeterministic 按自定义排序遍历", t, func() {
		m := NewBucketMap[string, int](4, stringHash)
		m.Store("c", 3)
		m.Store("a", 1)
		m.Store("b", 2)

		var got []string
		m.RangeDeterministic(
			func(k string, v int) bool {
				got = append(got, k)
				return true
			},
			func(keys []string) sort.Interface {
				return sort.StringSlice(keys)
			},
		)
		So(got, ShouldResemble, []string{"a", "b", "c"})
	})

	Convey("RangeDeterministic f 返回 false 中止", t, func() {
		m := NewBucketMap[int, int](4, func(k int) int64 { return int64(k) })
		for i := 1; i <= 5; i++ {
			m.Store(i, i)
		}
		var visited []int
		m.RangeDeterministic(
			func(k, v int) bool {
				visited = append(visited, k)
				return k < 3 // 走到 k=3 break
			},
			func(keys []int) sort.Interface {
				return sort.IntSlice(keys)
			},
		)
		So(visited, ShouldResemble, []int{1, 2, 3})
	})
}

// TestBucketMap_LoadOrStoreFunc[Error]Lock 双 check + bucket 锁路径。
func TestBucketMap_LoadOrStoreFuncLock(t *testing.T) {
	Convey("LoadOrStoreFuncLock 第一次创建 / 第二次 load", t, func() {
		m := NewBucketMap[string, int](4, stringHash)
		var calls atomic.Int32
		v1, loaded1 := m.LoadOrStoreFuncLock("new", func(k string) int {
			calls.Add(1)
			return 42
		})
		So(loaded1, ShouldBeFalse)
		So(v1, ShouldEqual, 42)
		So(calls.Load(), ShouldEqual, 1)

		// 再调一次：loaded=true，cf 不再被调
		v2, loaded2 := m.LoadOrStoreFuncLock("new", func(k string) int {
			calls.Add(1)
			return 999
		})
		So(loaded2, ShouldBeTrue)
		So(v2, ShouldEqual, 42)
		So(calls.Load(), ShouldEqual, 1)
	})

	Convey("LoadOrStoreFuncErrorLock cf 出错时不存值", t, func() {
		m := NewBucketMap[string, int](4, stringHash)
		errBoom := errors.New("boom")
		v, loaded, err := m.LoadOrStoreFuncErrorLock("err-k", func(k string) (int, error) {
			return 0, errBoom
		})
		So(err, ShouldEqual, errBoom)
		So(loaded, ShouldBeFalse)
		So(v, ShouldEqual, 0)
		// 后续 cf 成功则正常存
		v2, loaded2, err2 := m.LoadOrStoreFuncErrorLock("err-k", func(k string) (int, error) {
			return 7, nil
		})
		So(err2, ShouldBeNil)
		So(loaded2, ShouldBeFalse)
		So(v2, ShouldEqual, 7)
	})
}

// TestBucketMap_CacheKey 双缓冲 keySnapshots + 后台 refresh goroutine。
func TestBucketMap_CacheKey(t *testing.T) {
	Convey("NewBucketMapWithCacheKey snapshot + forceFresh + stop", t, func() {
		m, forceFresh, stop := NewBucketMapWithCacheKey[string, int](
			4, 50*time.Millisecond, stringHash)
		defer stop()

		// 初始：还没 Store，snapshot 应是空切片
		keys := m.Keys()
		So(keys, ShouldBeEmpty)

		m.Store("a", 1)
		m.Store("bb", 2)

		// 还没 forceFresh / 还没 tick 到，snapshot 还是空的
		keys = m.Keys()
		So(keys, ShouldBeEmpty)

		forceFresh()
		keys = m.Keys()
		sort.Strings(keys)
		So(keys, ShouldResemble, []string{"a", "bb"})

		// 等 tick：refreshDur=50ms，等 150ms 让后台 goroutine 跑过至少一轮
		m.Store("ccc", 3)
		time.Sleep(150 * time.Millisecond)
		keys = m.Keys()
		sort.Strings(keys)
		So(keys, ShouldResemble, []string{"a", "bb", "ccc"})
	})
}

// TestBucketMap_Concurrent 关键并发测试：100 goroutine × 50 次混合 CRUD。
// 必须配合 -race 跑（CI race job 守关）。
func TestBucketMap_Concurrent(t *testing.T) {
	const goroutines = 100
	const perGoroutine = 50

	m := NewBucketMap[string, int](16, stringHash)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				k := strconv.Itoa(gid*perGoroutine + i)
				m.Store(k, i)
				_, _ = m.Load(k)
				_ = m.Contains(k)
				if i%5 == 0 {
					m.Delete(k)
				}
				if i%7 == 0 {
					m.LoadOrStoreFuncLock(k+"-lazy", func(string) int { return 1 })
				}
			}
		}(g)
	}
	wg.Wait()

	// 最终态：Range 不死循环、Len 返回有效值
	count := m.Len()
	if count < 0 {
		t.Fatalf("Len returned negative: %d", count)
	}
}

// TestBucketMap_CacheKey_Concurrent 后台 goroutine + 并发 Store + 并发
// forceFresh：双缓冲切换路径的 race 检测。
func TestBucketMap_CacheKey_Concurrent(t *testing.T) {
	m, forceFresh, stop := NewBucketMapWithCacheKey[int, int](
		8, 10*time.Millisecond,
		func(k int) int64 { return int64(k) })
	defer stop()

	const goroutines = 50
	const perGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines + 1) // +1 for forceFresh refresher

	// 写入 + 读 snapshot 的多 goroutine
	for g := 0; g < goroutines; g++ {
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				m.Store(gid*perGoroutine+i, i)
				_ = m.Keys() // 读 snapshot 与后台 refresh 竞争
			}
		}(g)
	}

	// 主动 refresh 的 goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			forceFresh()
			time.Sleep(time.Millisecond)
		}
	}()

	wg.Wait()

	// 最后 forceFresh 一次，确认 snapshot 包含全部 key
	forceFresh()
	keys := m.Keys()
	if len(keys) != goroutines*perGoroutine {
		t.Errorf("Keys() lost data: got=%d want=%d", len(keys), goroutines*perGoroutine)
	}
}

// 顺带补 Map (any.go) 未覆盖的 Clear / RangeDeterministic（其它 Map 方法
// 已被 any_test.go 覆盖到 100%）。
func TestMap_ClearAndRangeDeterministic(t *testing.T) {
	Convey("Map.Clear", t, func() {
		var m Map[string, int]
		m.Store("a", 1)
		m.Store("b", 2)
		m.Clear()
		_, ok := m.Load("a")
		So(ok, ShouldBeFalse)
		_, ok2 := m.Load("b")
		So(ok2, ShouldBeFalse)
	})

	Convey("Map.RangeDeterministic 按自定义排序", t, func() {
		var m Map[string, int]
		m.Store("c", 3)
		m.Store("a", 1)
		m.Store("b", 2)

		var got []string
		m.RangeDeterministic(
			func(k string, v int) bool {
				got = append(got, k)
				return true
			},
			func(keys []string) sort.Interface {
				return sort.StringSlice(keys)
			},
		)
		So(got, ShouldResemble, []string{"a", "b", "c"})
	})

	Convey("Map.RangeDeterministic f 返回 false 中止", t, func() {
		var m Map[int, int]
		for i := 1; i <= 5; i++ {
			m.Store(i, i)
		}
		var visited []int
		m.RangeDeterministic(
			func(k, v int) bool {
				visited = append(visited, k)
				return k < 3
			},
			func(keys []int) sort.Interface {
				return sort.IntSlice(keys)
			},
		)
		So(visited, ShouldResemble, []int{1, 2, 3})
	})
}
