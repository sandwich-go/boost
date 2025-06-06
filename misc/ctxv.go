package misc

import "context"

type ContextValue[V any] interface {
	// Get 从 context.Context 中获取 V 值
	Get(context.Context) V
	// With 将 V 值存放至 context.Context 中， 并返回存有 V 值的 context.Context
	With(context.Context, V) context.Context
}

// NewContextValue 创建 ContextValue
func NewContextValue[V any]() ContextValue[V] {
	return &contextValueImpl[V]{}
}

type contextValueImpl[V any] struct{}

func (cv *contextValueImpl[V]) Get(ctx context.Context) V {
	v := ctx.Value(cv)
	if v == nil {
		var zero V
		return zero
	}

	return v.(V)
}

func (cv *contextValueImpl[V]) With(ctx context.Context, v V) context.Context {
	return context.WithValue(ctx, cv, v)
}
