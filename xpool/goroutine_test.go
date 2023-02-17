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
