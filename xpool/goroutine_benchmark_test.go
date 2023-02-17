package xpool

import (
	"sync"
	"testing"
	"time"
)

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

	for i := 0; i < b.N; i++ {
		wg.Add(runTimes)
		for j := 0; j < runTimes; j++ {
			pool.jobQueue <- func() {
				defer wg.Done()
				demoTask()
			}
		}
	}
}
