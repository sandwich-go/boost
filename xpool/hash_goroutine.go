package xpool

import (
	"github.com/sandwich-go/boost/z"
	"time"
)

type HashGoroutinePool struct {
	pools []*GoroutinePool
}

var DefaultHashFunc = z.MemHashString

func NewHashGoroutinePool(poolCount, jobQueueLen int, timeout time.Duration) *HashGoroutinePool {
	pools := make([]*GoroutinePool, poolCount)
	for i := range pools {
		pools[i] = NewGoroutinePool(1, jobQueueLen, timeout)
	}
	return &HashGoroutinePool{pools: pools}
}

func (h *HashGoroutinePool) getPool(key string) *GoroutinePool {
	index := DefaultHashFunc(key) % uint64(len(h.pools))
	return h.pools[index]
}

func (h *HashGoroutinePool) PushJob(key string, job Job) error {
	return h.getPool(key).Push(job)
}

func (h *HashGoroutinePool) Close() {
	for _, pool := range h.pools {
		pool.Close()
	}
}
