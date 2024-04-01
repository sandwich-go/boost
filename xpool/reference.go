package xpool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Reference 引用
type Reference interface {
	set(interface{})
	incr()

	// Decr 引用计数减一，当引用为0时，会添入对象池中
	Decr()
}

// ReferencePool 引用池
type ReferencePool interface {
	// Get 获取对象
	Get() interface{}
	// Stats 池子状态
	Stats() ReferencePoolStats
}

// ReferencePoolStats 引用池子状态
type ReferencePoolStats struct {
	Released   uint32
	Allocated  uint32
	Referenced uint32
}

func (s ReferencePoolStats) String() string {
	return fmt.Sprintf("Released: %d, Allocated: %d, Referenced: %d", s.Released, s.Allocated, s.Referenced)
}

type reference struct {
	count    *uint32           // 引用计数
	pool     *sync.Pool        // 释放的目标池
	released *uint32           // 释放次数
	instance interface{}       // 目标对象
	reset    func(interface{}) // 重置函数
	id       uint32            // 唯一标记
}

func (c *reference) set(i interface{}) { c.instance = i }
func (c *reference) incr()             { atomic.AddUint32(c.count, 1) }

// Decr 减少引用计数
func (c *reference) Decr() {
	if atomic.LoadUint32(c.count) == 0 {
		return
	}
	if atomic.AddUint32(c.count, ^uint32(0)) == 0 {
		atomic.AddUint32(c.released, 1)
		if c.reset != nil {
			c.reset(c.instance)
		}
		c.pool.Put(c.instance)
		c.instance = nil
	}
}

type referencePool struct {
	*sync.Pool
	released   uint32
	allocated  uint32
	referenced uint32
}

func NewReferencePool(builder func(Reference) Reference, reset func(interface{})) ReferencePool {
	p := new(referencePool)
	p.Pool = new(sync.Pool)
	p.Pool.New = func() interface{} {
		atomic.AddUint32(&p.allocated, 1)
		return builder(&reference{
			count:    new(uint32),
			pool:     p.Pool,
			released: &p.released,
			reset:    reset,
			id:       p.allocated,
		})
	}
	return p
}

func (p *referencePool) Get() interface{} {
	c := p.Pool.Get().(Reference)
	c.set(c)
	atomic.AddUint32(&p.referenced, 1)
	c.incr()
	return c
}

func (p *referencePool) Stats() ReferencePoolStats {
	return ReferencePoolStats{
		Allocated:  atomic.LoadUint32(&p.allocated),
		Referenced: atomic.LoadUint32(&p.referenced),
		Released:   atomic.LoadUint32(&p.released),
	}
}
