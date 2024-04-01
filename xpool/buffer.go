package xpool

import (
	"fmt"
	"sync"
)

// BytesPool bytes pool
type BytesPool interface {
	// Alloc try alloc a []byte from internal slab class if no free chunk in slab class Alloc will make one.
	Alloc(size int) []byte
	// Free release a []byte that alloc from BytesPool.Alloc.
	Free(mem []byte)
}

var debug bool

// SyncBytesPool is a sync.Pool base slab allocation memory pool
type SyncBytesPool struct {
	chunks  []sync.Pool
	sizes   []int
	minSize int
	maxSize int

	// for testing
	allocTimesFromPool int
	freeTimesToPool    int
}

// NewSyncBytesPool create a sync.Pool base slab allocation memory pool.
// minSize is the smallest chunk size.
// maxSize is the largest chunk size.
// factor is used to control growth of chunk size.
func NewSyncBytesPool(minSize, maxSize, factor int) BytesPool {
	n := 0
	if minSize <= 0 || factor <= 0 {
		panic(fmt.Sprintf("invalid paramter, minSize/factor should greater than 0"))
	}
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		n++
	}
	pool := &SyncBytesPool{
		chunks:  make([]sync.Pool, n),
		sizes:   make([]int, n),
		minSize: minSize, maxSize: maxSize,
	}
	n = 0
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		pool.sizes[n] = chunkSize
		pool.chunks[n].New = func(size int) func() interface{} {
			return func() interface{} {
				buf := make([]byte, size)
				return &buf
			}
		}(chunkSize)
		n++
	}
	return pool
}

// Alloc try alloc a []byte from internal slab class if no free chunk in slab class Alloc will make one.
func (p *SyncBytesPool) Alloc(size int) []byte {
	if size <= p.maxSize {
		for i := 0; i < len(p.sizes); i++ {
			if p.sizes[i] >= size {
				mem := p.chunks[i].Get().(*[]byte)
				if debug {
					p.allocTimesFromPool++
				}
				return (*mem)[:size]
			}
		}
	}
	return make([]byte, size)
}

// Free release a []byte that alloc from SyncBytesPool.Alloc.
func (p *SyncBytesPool) Free(mem []byte) {
	if size := cap(mem); size <= p.maxSize {
		for i := 0; i < len(p.sizes); i++ {
			if p.sizes[i] == size {
				p.chunks[i].Put(&mem)
				if debug {
					p.freeTimesToPool++
				}
				return
			} else if p.sizes[i] > size && i > 0 && p.sizes[i-1] <= size {
				p.chunks[i-1].Put(&mem)
				if debug {
					p.freeTimesToPool++
				}
				return
			}
		}
	}
}
