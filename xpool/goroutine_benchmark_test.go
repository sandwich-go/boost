package xpool

import (
	"context"
	"sync"
	"testing"
	"time"
)

// runTimes / poolSize / queueSize 选取动机：
//   - runTimes=1e6: 让 bench 跑出稳定的 alloc 数（job closure 1 alloc + wg
//     底层 1 alloc + sleep timer 1-2 alloc，理论 ≈ 4M alloc/iter）
//   - poolSize=50000: 50k worker 并行 50k * 10ms = 500ms 干完 1e6 job 的
//     5%，总 200ms 理论；但 queueSize<<runTimes 会让 push 反复阻塞
//   - queueSize=5000: 故意小，让 push 与 worker 消费速率匹配，避免内存爆
//
// 历史 bug（commit `<本批>` 修）：原 BenchmarkGoroutinePool 缺 wg.Wait()，
// 测的只是"push 到 jobQueue 速度"不是端到端 push+job_done 时长，导致与
// BenchmarkGoroutine 不可比。
const (
	runTimes  = 1000000
	poolSize  = 50000
	queueSize = 5000
)

func demoTask() {
	time.Sleep(time.Millisecond * 10)
}

func BenchmarkGoroutine(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(runTimes)

		for j := 0; j < runTimes; j++ {
			go func() {
				defer wg.Done()
				demoTask()
			}()
		}

		wg.Wait()
	}
}

func BenchmarkGoroutinePool(b *testing.B) {
	pool := NewGoroutinePool(poolSize, queueSize, time.Duration(0))
	defer pool.Close()
	var wg sync.WaitGroup
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		wg.Add(runTimes)
		for j := 0; j < runTimes; j++ {
			if err := pool.Push(ctx, func() {
				defer wg.Done()
				demoTask()
			}); err != nil {
				b.Fatalf("Push: %v", err)
			}
		}
		// 等所有 job 跑完才结束本次迭代。原 bench 缺这步 → 测的不是
		// 端到端 latency。Wait 之后才能与 BenchmarkGoroutine 直接对比。
		wg.Wait()
	}
}
