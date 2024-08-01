package misc

import (
	"github.com/sandwich-go/boost/misc/xbit"
	"github.com/sandwich-go/boost/xpool"
)

// Archive orm record 回滚使用的协助类
type Archive[T any] struct {
	Fields xbit.FieldSet
	Data   T
}

type ArchivePool[T any] struct {
	pool *xpool.Pool[*Archive[T]]
}

// NewArchivePool 创建 Archive 池
func NewArchivePool[T any]() *ArchivePool[T] {
	return &ArchivePool[T]{
		pool: xpool.NewPool(func() *Archive[T] {
			return &Archive[T]{}
		}),
	}
}

func (p *ArchivePool[T]) Get(data *T) *Archive[T] {
	a := p.pool.Get()
	a.Data = *data
	return a
}

func (p *ArchivePool[T]) Put(a *Archive[T]) {
	if a == nil {
		return
	}
	var zero T
	a.Data = zero
	a.Fields = 0
	p.pool.Put(a)
}
