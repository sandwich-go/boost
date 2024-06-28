package xpool

import (
	"sync"
)

type Pool[T any] struct {
	native sync.Pool
}

// NewPool 创建池
func NewPool[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		native: sync.Pool{
			New: func() interface{} { return fn() },
		},
	}
}

// Get 从池子里获取对象
func (p *Pool[T]) Get() T {
	return p.native.Get().(T)
}

// Put 将对象放入池子
func (p *Pool[T]) Put(x T) { p.native.Put(x) }
