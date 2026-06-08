package xpool

import (
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
