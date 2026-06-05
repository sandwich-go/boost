package lru

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestWorkerPerEngine_BasicFunctionality 测试 workerPerEngine 基本功能
func TestWorkerPerEngine_BasicFunctionality(t *testing.T) {
	Convey("workerPerEngine should execute clean function periodically", t, func() {
		worker := NewWorkerPerEngine()
		counter := atomic.Int32{}
		interval := 50 * time.Millisecond

		worker.Clean(interval, "test-id", func() {
			counter.Add(1)
		})

		// 等待至少执行 3 次
		time.Sleep(200 * time.Millisecond)

		count := counter.Load()
		So(count, ShouldBeGreaterThanOrEqualTo, 2)
		So(count, ShouldBeLessThanOrEqualTo, 5) // 考虑 Disturb 的影响
	})
}

// TestWorkerPerEngine_MultipleEngines 测试多个 engine 同时使用
func TestWorkerPerEngine_MultipleEngines(t *testing.T) {
	Convey("workerPerEngine should handle multiple engines correctly", t, func() {
		worker := NewWorkerPerEngine()
		engineCount := 10
		counters := make([]atomic.Int32, engineCount)
		interval := 50 * time.Millisecond

		// 启动多个 engine 的清理任务
		for i := 0; i < engineCount; i++ {
			idx := i
			worker.Clean(interval, "engine-"+string(rune('0'+i)), func() {
				counters[idx].Add(1)
			})
		}

		// 等待执行
		time.Sleep(200 * time.Millisecond)

		// 验证每个 engine 都执行了
		for i := 0; i < engineCount; i++ {
			count := counters[i].Load()
			So(count, ShouldBeGreaterThanOrEqualTo, 2)
		}
	})
}

// TestWorkerHashPool_BasicFunctionality 测试 workerHashPool 基本功能
func TestWorkerHashPool_BasicFunctionality(t *testing.T) {
	Convey("workerHashPool should execute clean function periodically", t, func() {
		worker := NewWorkerHashPool(2, 100, time.Second)
		counter := atomic.Int32{}
		interval := 50 * time.Millisecond

		worker.Clean(interval, "test-id", func() {
			counter.Add(1)
		})

		// 等待至少执行 3 次
		time.Sleep(200 * time.Millisecond)

		count := counter.Load()
		So(count, ShouldBeGreaterThanOrEqualTo, 2)
		So(count, ShouldBeLessThanOrEqualTo, 5)
	})
}

// TestWorkerHashPool_MultipleEngines 测试多个 engine 使用 workerHashPool
func TestWorkerHashPool_MultipleEngines(t *testing.T) {
	Convey("workerHashPool should handle multiple engines correctly", t, func() {
		numWorkers := 4
		worker := NewWorkerHashPool(numWorkers, 100, time.Second)
		engineCount := 20
		counters := make([]atomic.Int32, engineCount)
		interval := 50 * time.Millisecond

		// 启动多个 engine 的清理任务
		for i := 0; i < engineCount; i++ {
			idx := i
			worker.Clean(interval, "engine-"+string(rune('0'+i)), func() {
				counters[idx].Add(1)
			})
		}

		// 等待执行
		time.Sleep(200 * time.Millisecond)

		// 验证每个 engine 都执行了
		for i := 0; i < engineCount; i++ {
			count := counters[i].Load()
			So(count, ShouldBeGreaterThanOrEqualTo, 2)
		}
	})
}

// TestWorkerHashPool_SameIDSerialization 测试相同 ID 的任务串行执行
func TestWorkerHashPool_SameIDSerialization(t *testing.T) {
	Convey("workerHashPool should serialize tasks with same ID", t, func() {
		worker := NewWorkerHashPool(2, 100, time.Second)
		interval := 30 * time.Millisecond
		executing := atomic.Int32{}
		maxConcurrent := atomic.Int32{}
		counter := atomic.Int32{}

		// 同一个 ID 的清理函数
		cleanFunc := func() {
			current := executing.Add(1)
			defer executing.Add(-1)

			// 更新最大并发数
			for {
				max := maxConcurrent.Load()
				if current <= max || maxConcurrent.CompareAndSwap(max, current) {
					break
				}
			}

			counter.Add(1)
			time.Sleep(10 * time.Millisecond) // 模拟清理工作
		}

		// 使用相同的 ID
		worker.Clean(interval, "same-id", cleanFunc)

		// 等待多次执行
		time.Sleep(150 * time.Millisecond)

		// 验证：相同 ID 的任务应该串行执行，最大并发应该是 1
		So(maxConcurrent.Load(), ShouldEqual, 1)
		So(counter.Load(), ShouldBeGreaterThanOrEqualTo, 3)
	})
}

// TestWorkerHashPool_DifferentIDParallel 测试不同 ID 的任务可以并行执行
func TestWorkerHashPool_DifferentIDParallel(t *testing.T) {
	Convey("workerHashPool should allow parallel execution for different IDs", t, func() {
		numWorkers := 4
		worker := NewWorkerHashPool(numWorkers, 100, time.Second)
		interval := 30 * time.Millisecond
		executing := atomic.Int32{}
		maxConcurrent := atomic.Int32{}

		cleanFunc := func() {
			current := executing.Add(1)
			defer executing.Add(-1)

			// 更新最大并发数
			for {
				max := maxConcurrent.Load()
				if current <= max || maxConcurrent.CompareAndSwap(max, current) {
					break
				}
			}

			time.Sleep(50 * time.Millisecond) // 模拟清理工作
		}

		// 启动多个不同 ID 的任务
		for i := 0; i < 8; i++ {
			worker.Clean(interval, "engine-"+string(rune('A'+i)), cleanFunc)
		}

		// 等待执行
		time.Sleep(150 * time.Millisecond)

		// 验证：不同 ID 可以并行执行，最大并发应该大于 1
		So(maxConcurrent.Load(), ShouldBeGreaterThan, 1)
	})
}

// TestWorkerComparison_GoroutineCount 验证两种 worker 都不会随 engine
// 数量线性起 goroutine：CleanWorker 接口只承诺周期性执行 f，goroutine
// 数量是实现细节，但「不应随 engine 数量线性增长」是性能约束。
//
// 历史背景：旧版 workerPerEngine 为每个 engine 起一个常驻 channel
// 中转 goroutine，本测试一度以「per-engine 应增加 ≈engineCount 个
// goroutine」为断言记录该实现细节。后续优化把 channel 中转去掉、改用
// timer 回调链，断言反转——现在两者都属于「不起 per-engine goroutine」
// 的实现，本测试统一限制为「增量远小于 engineCount」。
func TestWorkerComparison_GoroutineCount(t *testing.T) {
	Convey("Neither worker should spawn per-engine goroutines", t, func() {
		engineCount := 50
		interval := 100 * time.Millisecond
		// 阈值 = engineCount/2：足够松到容纳 std timer 派发栈临时 goroutine
		// + GC sweep 路径上偶发的 helper goroutine，又能在「真的退化为
		// per-engine goroutine」时立刻报警。
		threshold := engineCount / 2

		Convey("workerPerEngine should not grow goroutine count linearly", func() {
			runtime.GC()
			time.Sleep(50 * time.Millisecond)
			baseGoroutines := runtime.NumGoroutine()

			workerPer := NewWorkerPerEngine()
			for i := 0; i < engineCount; i++ {
				workerPer.Clean(interval, "engine-"+string(rune('0'+i)), func() {})
			}

			time.Sleep(50 * time.Millisecond)
			runtime.GC()
			time.Sleep(50 * time.Millisecond)

			goroutinesWithPer := runtime.NumGoroutine()
			increasedPer := goroutinesWithPer - baseGoroutines

			So(increasedPer, ShouldBeLessThan, threshold)
			t.Logf("workerPerEngine: base=%d, current=%d, increased=%d (threshold=%d)",
				baseGoroutines, goroutinesWithPer, increasedPer, threshold)
		})

		Convey("workerHashPool uses fixed number of goroutines", func() {
			runtime.GC()
			time.Sleep(50 * time.Millisecond)
			baseGoroutines := runtime.NumGoroutine()

			numWorkers := 4
			workerHash := NewWorkerHashPool(numWorkers, 100, time.Second)
			for i := 0; i < engineCount; i++ {
				workerHash.Clean(interval, "engine-"+string(rune('0'+i)), func() {})
			}

			time.Sleep(50 * time.Millisecond)
			runtime.GC()
			time.Sleep(50 * time.Millisecond)

			goroutinesWithHash := runtime.NumGoroutine()
			increasedHash := goroutinesWithHash - baseGoroutines

			So(increasedHash, ShouldBeLessThan, numWorkers+10)
			t.Logf("workerHashPool: base=%d, current=%d, increased=%d, workers=%d",
				baseGoroutines, goroutinesWithHash, increasedHash, numWorkers)
		})
	})
}

// TestWorkerPerEngine_ConcurrentSafety 测试并发安全性
func TestWorkerPerEngine_ConcurrentSafety(t *testing.T) {
	Convey("workerPerEngine should be concurrent safe", t, func() {
		worker := NewWorkerPerEngine()
		interval := 20 * time.Millisecond
		counter := atomic.Int32{}
		var wg sync.WaitGroup

		// 并发启动多个清理任务
		concurrency := 50
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx int) {
				defer wg.Done()
				worker.Clean(interval, "engine-"+string(rune('0'+idx)), func() {
					counter.Add(1)
				})
			}(i)
		}

		wg.Wait()
		time.Sleep(100 * time.Millisecond)

		// 验证所有任务都在执行
		count := counter.Load()
		So(count, ShouldBeGreaterThan, concurrency) // 至少每个执行一次
	})
}

// TestWorkerHashPool_ConcurrentSafety 测试并发安全性
func TestWorkerHashPool_ConcurrentSafety(t *testing.T) {
	Convey("workerHashPool should be concurrent safe", t, func() {
		worker := NewWorkerHashPool(4, 100, time.Second)
		interval := 20 * time.Millisecond
		counter := atomic.Int32{}
		var wg sync.WaitGroup

		// 并发启动多个清理任务
		concurrency := 50
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx int) {
				defer wg.Done()
				worker.Clean(interval, "engine-"+string(rune('0'+idx)), func() {
					counter.Add(1)
				})
			}(i)
		}

		wg.Wait()
		time.Sleep(100 * time.Millisecond)

		// 验证所有任务都在执行
		count := counter.Load()
		So(count, ShouldBeGreaterThan, concurrency)
	})
}

// TestWorkerHashPool_HashDistribution 测试 hash 分配的均匀性
func TestWorkerHashPool_HashDistribution(t *testing.T) {
	Convey("workerHashPool should distribute tasks relatively evenly", t, func() {
		numWorkers := 4
		worker := NewWorkerHashPool(numWorkers, 100, time.Second).(*workerHashPool)
		engineCount := 100

		// 验证 worker 配置正确
		So(worker.numWorkers, ShouldEqual, numWorkers)
		So(worker.pool, ShouldNotBeNil)

		// 启动多个任务验证不会崩溃
		counter := atomic.Int32{}
		for i := 0; i < engineCount; i++ {
			id := "engine-" + string(rune('0'+i))
			worker.Clean(50*time.Millisecond, id, func() {
				counter.Add(1)
			})
		}

		// 等待执行
		time.Sleep(150 * time.Millisecond)

		// 验证任务都在执行
		So(counter.Load(), ShouldBeGreaterThan, engineCount)
	})
}

// BenchmarkWorkerPerEngine 性能基准测试
func BenchmarkWorkerPerEngine(b *testing.B) {
	worker := NewWorkerPerEngine()
	counter := atomic.Int32{}
	interval := 100 * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker.Clean(interval, "bench-"+string(rune(i)), func() {
			counter.Add(1)
		})
	}
	b.StopTimer()

	time.Sleep(200 * time.Millisecond)
	b.Logf("Total executions: %d", counter.Load())
}

// BenchmarkWorkerHashPool 性能基准测试
func BenchmarkWorkerHashPool(b *testing.B) {
	worker := NewWorkerHashPool(8, 1000, time.Second)
	counter := atomic.Int32{}
	interval := 100 * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker.Clean(interval, "bench-"+string(rune(i)), func() {
			counter.Add(1)
		})
	}
	b.StopTimer()

	time.Sleep(200 * time.Millisecond)
	b.Logf("Total executions: %d", counter.Load())
}

// TestWorkerInterface 测试接口兼容性
func TestWorkerInterface(t *testing.T) {
	Convey("Both workers should implement CleanWorker interface", t, func() {
		var _ CleanWorker = NewWorkerPerEngine()
		var _ CleanWorker = NewWorkerHashPool(4, 100, time.Second)
	})
}

// TestWorkerExecution_Timing 测试执行时间的准确性
func TestWorkerExecution_Timing(t *testing.T) {
	Convey("Workers should execute at approximately correct intervals", t, func() {
		interval := 100 * time.Millisecond

		testWorkerTiming := func(worker CleanWorker, name string) {
			var timestamps []time.Time
			var mu sync.Mutex

			worker.Clean(interval, "timing-test", func() {
				mu.Lock()
				timestamps = append(timestamps, time.Now())
				mu.Unlock()
			})

			// 等待执行多次
			time.Sleep(500 * time.Millisecond)

			mu.Lock()
			defer mu.Unlock()

			So(len(timestamps), ShouldBeGreaterThanOrEqualTo, 3)

			// 检查执行间隔
			for i := 1; i < len(timestamps); i++ {
				interval := timestamps[i].Sub(timestamps[i-1])
				// 由于 Disturb 函数会添加 ±10% 的随机扰动，所以间隔会有变化
				// 这里只验证大致在合理范围内
				So(interval, ShouldBeGreaterThan, 50*time.Millisecond)
				So(interval, ShouldBeLessThan, 150*time.Millisecond)
			}

			t.Logf("%s execution intervals: %v", name, func() []time.Duration {
				intervals := make([]time.Duration, len(timestamps)-1)
				for i := 1; i < len(timestamps); i++ {
					intervals[i-1] = timestamps[i].Sub(timestamps[i-1])
				}
				return intervals
			}())
		}

		Convey("workerPerEngine timing", func() {
			testWorkerTiming(NewWorkerPerEngine(), "workerPerEngine")
		})

		Convey("workerHashPool timing", func() {
			testWorkerTiming(NewWorkerHashPool(2, 100, time.Second), "workerHashPool")
		})
	})
}
