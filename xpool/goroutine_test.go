package xpool

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestNewPool1(t *testing.T) {
	pool := NewGoroutinePool(10, 100, time.Duration(0))
	pool.SetSize(5)
	pool.SetSize(8)
	pool.SetSize(20)
}

func TestNewPool(t *testing.T) {
	pool := NewGoroutinePool(10, 100, time.Duration(0))
	defer pool.Close()

	iterations := 20
	var counter uint64 = 0

	wg := sync.WaitGroup{}
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		arg := uint64(1)
		job := func() {
			defer wg.Done()
			time.Sleep(time.Duration(1) * time.Second)
			atomic.AddUint64(&counter, arg)
		}

		pool.jobQueue <- job
	}
	wg.Wait()

	counterFinal := atomic.LoadUint64(&counter)
	if uint64(iterations) != counterFinal {
		t.Errorf("iterations %v is not equal counterFinal %v", iterations, counterFinal)
	}
}

// TestGoroutinePool_Close_WaitsForWorkers 验证 Close 后 worker goroutine 真
// 退出（不再悬挂）。
//
// 历史 bug: worker.Start 原代码 'defer close(w.closedChan)' 写在 Start 函数
// 顶部，会在 Start 返回时（go func 起完即返）立刻 close closedChan，让
// worker.join() 不再等 worker goroutine 真退出。表现：Close 看似立即返回但
// worker goroutine 仍在跑（goroutine 泄漏）。
//
// 反向验证（§3.1）：临时把 'defer close(w.closedChan)' 移回 Start 顶部，
// Close 后 NumGoroutine 仍含 worker (本测试断言会 fail)。
func TestGoroutinePool_Close_WaitsForWorkers(t *testing.T) {
	runtime.GC()
	time.Sleep(50 * time.Millisecond)
	baseGoroutines := runtime.NumGoroutine()

	const numWorkers = 50
	pool := NewGoroutinePool(numWorkers, 100, time.Duration(0))

	// 等 worker goroutine 全起来（每个 worker 1 goroutine）
	time.Sleep(50 * time.Millisecond)
	if got := runtime.NumGoroutine(); got < baseGoroutines+numWorkers {
		t.Fatalf("workers not all started: base=%d, current=%d, want >= %d",
			baseGoroutines, got, baseGoroutines+numWorkers)
	}

	// Close 后所有 worker goroutine 应真退出（join 等 closedChan）。
	// 修复后 Close 调用是同步等 worker 真退；修复前 join 立即返 →
	// Close 立即返但 worker 仍可能在跑（jobQueue close 后才退）。
	// 这里 Close 后**不 sleep 不 GC** 立即检查 NumGoroutine：
	pool.Close()
	got := runtime.NumGoroutine()
	// 给 2 个 goroutine 余量（GC sweep 等系统 goroutine）
	if got > baseGoroutines+2 {
		t.Errorf("worker goroutines leaked after Close: base=%d, current=%d (delta=%d)",
			baseGoroutines, got, got-baseGoroutines)
	}
}

// TestGoroutinePool_Close_DrainsResidualJobs 验证 Close 时 jobQueue 里残余
// 的 job 被 worker drain 跑完才退出，**不**丢失。
//
// 设计动机：dataserver 的 Publish 路径 in-flight job 走 asyncQueue.Produce
// 把落库 task 入 stream——若 Close 时 jobQueue 残余 job 被丢失，等同丢未
// 落库任务。详见 worker.Start 注释"历史 bug 2"段。
//
// 反向验证锚点（§3.1）：把 worker.Start 里 closeChan case 的 drain 内层
// for-select 改回 'return' → 残余 job 不会跑 → counter < N → 测试 fail。
//
// 测试构造：
//
//   - 单 worker（避免并行 drain 的随机性，让"是否 drain"判定确定）
//   - 队列预填 N 个慢 job（每个 sleep）—— Push 完后 worker 还没消费完
//   - Close 触发 SetSize(0) → worker 收 closeChan → 切 drain → 跑完所有
//   - 断言 counter == N
func TestGoroutinePool_Close_DrainsResidualJobs(t *testing.T) {
	const numJobs = 20
	// 单 worker 让 drain 顺序确定；timeout=0 让 Push 阻塞到队列有空
	// （但 capacity 设大于 numJobs 不会阻塞）。
	pool := NewGoroutinePool(1, numJobs+1, time.Duration(0))

	var counter atomic.Uint64
	var wg sync.WaitGroup
	wg.Add(numJobs)
	ctx := context.Background()
	for i := 0; i < numJobs; i++ {
		err := pool.Push(ctx, func() {
			defer wg.Done()
			// 让 job 慢一点确保 Close 时还有残余在 jobQueue 里
			time.Sleep(10 * time.Millisecond)
			counter.Add(1)
		})
		if err != nil {
			t.Fatalf("Push %d: %v", i, err)
		}
	}

	// Push 完立即 Close —— 此时 worker 大概率只跑了 1-2 个 job，jobQueue
	// 里至少 18+ 个 job 残余。Close 应同步等所有残余 drain 完。
	pool.Close()

	// 强不变量 1：drain 完成后所有 job 都跑过
	if got := counter.Load(); got != numJobs {
		t.Errorf("residual jobs not drained: counter=%d, want=%d", got, numJobs)
	}
	// 强不变量 2：wg.Wait 不应阻塞（所有 job.Done 都跑过）
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		// OK
	case <-time.After(2 * time.Second):
		t.Fatal("wg.Wait timeout: 部分 job 未执行 Done() —— drain 不全")
	}
}

// TestGoroutinePool_Close_PushAfterCloseShortCircuits 验证 Close 后 Push
// 立即短路返 "pool closed" error，不阻塞、不 panic。
func TestGoroutinePool_Close_PushAfterCloseShortCircuits(t *testing.T) {
	pool := NewGoroutinePool(2, 10, time.Duration(0))
	pool.Close()

	if !pool.IsClosed() {
		t.Fatal("Close 后 IsClosed 应为 true")
	}

	err := pool.Push(context.Background(), func() {})
	if err == nil {
		t.Fatal("Close 后 Push 应返 error，got nil")
	}
}

// TestGoroutinePool_Close_Idempotent 验证 Close 幂等——重复调不 panic、
// CAS 让第二次起 short-circuit。
func TestGoroutinePool_Close_Idempotent(t *testing.T) {
	pool := NewGoroutinePool(2, 10, time.Duration(0))
	pool.Close()
	pool.Close() // 不 panic
	pool.Close() // 不 panic
}
