package middleware

import (
	"context"
)

type (
	// Handler 处理器
	Handler func(context.Context) error
	// Middleware 中间件
	Middleware func(ctx context.Context, next Handler) error
	// Middlewares 中间件集合
	Middlewares []Middleware
)

// Chain 中间件链
func Chain(middlewares ...Middleware) Middleware {
	n := len(middlewares)
	return func(ctx context.Context, next Handler) error {
		chain := func(currMiddleware Middleware, currHandler Handler) Handler {
			return func(currCtx context.Context) error {
				return currMiddleware(currCtx, currHandler)
			}
		}
		chainHandler := next
		for i := n - 1; i >= 0; i-- {
			if middlewares[i] == nil {
				continue
			}
			chainHandler = chain(middlewares[i], chainHandler)
		}
		return chainHandler(ctx)
	}
}

// Use 中间件链
func Use(first Middleware, middlewares ...Middleware) Middleware {
	if first == nil {
		first = Chain(middlewares...)
	} else {
		first = Chain(first, Chain(middlewares...))
	}
	return first
}
