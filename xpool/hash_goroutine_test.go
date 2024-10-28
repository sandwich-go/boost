package xpool

import (
	"fmt"
	"github.com/sandwich-go/boost/z"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewHashGoroutinePool(t *testing.T) {
	const (
		poolCount  = 10
		jobQueue   = 100
		iterations = 300
	)
	var (
		counter [poolCount]uint64
	)

	DefaultHashFunc = func(str string) uint64 {
		i, _ := strconv.Atoi(str)
		return uint64(i)
	}
	defer func() {
		DefaultHashFunc = z.MemHashString
	}()
	pool := NewHashGoroutinePool(poolCount, jobQueue, time.Second)
	defer pool.Close()

	wg := sync.WaitGroup{}
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		idx := uint64(i % poolCount)
		job := func() {
			defer wg.Done()
			atomic.AddUint64(&counter[idx], 1)
		}
		pool.PushJob(fmt.Sprintf("%d", i), job)
	}
	wg.Wait()
	counterFinal := iterations / poolCount
	for i := 0; i < poolCount; i++ {
		t.Log("result:", i, counter[i])
		if counter[i] != uint64(counterFinal) {
			t.Errorf("%v is not equal counterFinal %v", counter[i], counterFinal)
		}
	}
}
