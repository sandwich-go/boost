package xrand

import (
	"sync"
	"testing"
	_ "unsafe"
)

// 并发测试方法
func benchmarkConcurrentRandomInt(b *testing.B, randomFunc func(int, int) int, goroutines int) {
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < b.N/goroutines; j++ {
				_ = randomFunc(1, 100)
			}
		}()
	}

	wg.Wait()
}

// 并发 Benchmark: FastRandInt
func BenchmarkFastRandInt_Concurrent(b *testing.B) {
	benchmarkConcurrentRandomInt(b, FastRandInt, 1000)
}

// 并发 Benchmark: RandomInt
func BenchmarkRandomInt_Concurrent(b *testing.B) {
	benchmarkConcurrentRandomInt(b, RandomInt, 1000)
}
