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
		panic(fmt.Sprintf("invalid parameter, minSize/factor should greater than 0"))
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

// Alloc 从 slab 分级中借一块 []byte（无空闲则新建）。
//
// 契约（调用方必须遵守，否则 Free 会静默误投/丢弃）：
//   - 返回的 slice: len==size，但 cap==所命中档位的 chunk size（可能大于 size）。
//   - Free 时必须传回 cap 未被改变的原 slice：不要在 Free 前 append 越过 cap，
//     也不要用 s[:n:n] 这类三索引 reslice 缩小 cap——那会让 Free 依据错误的 cap
//     误投或丢弃，poison 对应桶。
//   - 不要在 Free 之后继续持有/读写该 slice（内存已归还池，可能被他人复用）。
//   - Free 一个非 Alloc 来源、或超出 [minSize,maxSize] 的 slice 是故意的 silent
//     no-op（debug 模式下会 panic 以便测试期捕获误用）。
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

// Free 归还一块由 Alloc 借出的 []byte（按 cap 定位桶）。见 Alloc 的契约说明。
// Free 一个非 Alloc 来源、cap 越界或被改变的 slice 是 silent no-op（按 cap 无法
// 可靠区分"合法的 oversize 往返"与"越界误用"，故不做运行时断言）。
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
