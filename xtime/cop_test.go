package xtime

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkTimeUnixWithSystem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NowFunc().Unix()
	}
}

func BenchmarkTimeUnixWithCop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Unix()
	}
}

func BenchmarkCompareSystemAndCop(b *testing.B) {
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			for i := 0; i < b.N; i++ {
				n := Unix()
				if n > (NowFunc().Unix()+1) || n < (NowFunc().Unix()-1) {
					fmt.Println("Error Cop:", n, "NowFunc().Unix():", NowFunc().Unix())
					b.Fail()
					return
				}
				time.Sleep(1 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
}
