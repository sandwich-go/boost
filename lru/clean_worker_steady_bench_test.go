package lru

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

// 稳态 GC 基线测试：先注册 engineCount 个清理任务，让循环跑稳态后
// 测量「单位时间产生的 alloc 字节 / GC 次数」，反映 reschedule 路径
// 上的 timer / channel / closure 持续分配压力。
//
// 不用 b.N（b.N 控制注册数量，但稳态压力来自时间累积），用固定 wallTime
// 配合 b.ReportMetric 输出每秒 alloc/GC，用 -benchtime 控制采样时长。
//
// 跑法：go test -run='^$' -bench=Steady -benchmem ./lru/

const (
	steadyEngineCount = 200
	steadyInterval    = 5 * time.Millisecond // 让一秒内多轮 reschedule 累积
	steadyWarmup      = 200 * time.Millisecond
	steadySample      = 2 * time.Second
)

func benchSteady(b *testing.B, name string, mk func() CleanWorker) {
	b.Helper()
	worker := mk()
	var counter atomic.Int64

	for i := 0; i < steadyEngineCount; i++ {
		id := fmt.Sprintf("steady-%s-%d", name, i)
		worker.Clean(steadyInterval, id, func() {
			counter.Add(1)
		})
	}

	// 等稳态
	time.Sleep(steadyWarmup)

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	startCount := counter.Load()
	start := time.Now()

	time.Sleep(steadySample)

	runtime.ReadMemStats(&after)
	elapsed := time.Since(start)
	executions := counter.Load() - startCount

	allocBytes := after.TotalAlloc - before.TotalAlloc
	allocCount := after.Mallocs - before.Mallocs
	gcCount := after.NumGC - before.NumGC

	if executions == 0 {
		b.Fatalf("%s: no executions in steady window", name)
	}

	// 归一到「每次清理回调」的成本——用户关心的就是单次循环的 GC 摊销
	b.ReportMetric(float64(allocBytes)/float64(executions), "B/exec")
	b.ReportMetric(float64(allocCount)/float64(executions), "allocs/exec")
	b.ReportMetric(float64(executions)/elapsed.Seconds(), "exec/s")
	b.ReportMetric(float64(gcCount), "gc-cycles")
	b.ReportMetric(float64(allocBytes)/elapsed.Seconds()/1024, "KB/s")
}

func BenchmarkSteady_WorkerPerEngine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSteady(b, "per", NewWorkerPerEngine)
	}
}

func BenchmarkSteady_WorkerHashPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSteady(b, "hash", func() CleanWorker {
			return NewWorkerHashPool(8, 4096, time.Second)
		})
	}
}
