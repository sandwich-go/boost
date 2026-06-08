package middleware

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	someValue     = 1
	parentContext = context.WithValue(context.TODO(), "parent", someValue)
)

func TestChain(t *testing.T) {
	Convey("TestChain", t, func() {
		first := func(ctx context.Context, next Handler) error {
			requireContextValue(ctx, "parent", "first interceptor must know the parent context value")
			ctx = context.WithValue(ctx, "first", 1)
			return next(ctx)
		}
		second := func(ctx context.Context, next Handler) error {
			requireContextValue(ctx, "parent", "second interceptor must know the parent context value")
			requireContextValue(ctx, "first", "second interceptor must know the first context value")
			ctx = context.WithValue(ctx, "second", 1)
			return next(ctx)
		}

		handler := func(ctx context.Context) error {
			requireContextValue(ctx, "parent", "handler must know the parent context value")
			requireContextValue(ctx, "first", "handler must know the first context value")
			requireContextValue(ctx, "second", "handler must know the second context value")
			return nil
		}
		chain := Chain(first, second)
		err := chain(parentContext, handler)
		So(err, ShouldBeNil)
	})
}

func requireContextValue(ctx context.Context, key string, msg string) {
	Convey(msg, func() {
		val := ctx.Value(key)
		So(val, ShouldNotBeNil)
		So(someValue, ShouldEqual, val)
	})
}

// TestUse 覆盖 Use 函数（first==nil 走 Chain(middlewares...)，
// first!=nil 走 Chain(first, Chain(...))）。
func TestUse(t *testing.T) {
	hit := func(label string) Middleware {
		return func(ctx context.Context, next Handler) error {
			ctx = context.WithValue(ctx, label, true)
			return next(ctx)
		}
	}

	Convey("Use first==nil → Chain(middlewares...)", t, func() {
		mw := Use(nil, hit("a"), hit("b"))
		var seenA, seenB bool
		err := mw(parentContext, func(ctx context.Context) error {
			seenA = ctx.Value("a") == true
			seenB = ctx.Value("b") == true
			return nil
		})
		So(err, ShouldBeNil)
		So(seenA, ShouldBeTrue)
		So(seenB, ShouldBeTrue)
	})

	Convey("Use first!=nil → Chain(first, Chain(middlewares...))", t, func() {
		mw := Use(hit("first"), hit("a"), hit("b"))
		var seenFirst, seenA, seenB bool
		err := mw(parentContext, func(ctx context.Context) error {
			seenFirst = ctx.Value("first") == true
			seenA = ctx.Value("a") == true
			seenB = ctx.Value("b") == true
			return nil
		})
		So(err, ShouldBeNil)
		So(seenFirst, ShouldBeTrue)
		So(seenA, ShouldBeTrue)
		So(seenB, ShouldBeTrue)
	})

	Convey("Use 仅 first 一个（middlewares 为空）", t, func() {
		mw := Use(hit("only"))
		var seen bool
		err := mw(parentContext, func(ctx context.Context) error {
			seen = ctx.Value("only") == true
			return nil
		})
		So(err, ShouldBeNil)
		So(seen, ShouldBeTrue)
	})
}

// TestChain_NilMiddleware 覆盖 Chain 内 'middlewares[i] == nil → continue'
// 跳过分支（line 27-29）。原 TestChain 三个 middleware 都非 nil。
func TestChain_NilMiddleware(t *testing.T) {
	Convey("Chain 内含 nil middleware 跳过不报错", t, func() {
		hit := func(label string) Middleware {
			return func(ctx context.Context, next Handler) error {
				ctx = context.WithValue(ctx, label, true)
				return next(ctx)
			}
		}
		chain := Chain(hit("a"), nil, hit("b"))
		var seenA, seenB bool
		err := chain(parentContext, func(ctx context.Context) error {
			seenA = ctx.Value("a") == true
			seenB = ctx.Value("b") == true
			return nil
		})
		So(err, ShouldBeNil)
		So(seenA, ShouldBeTrue)
		So(seenB, ShouldBeTrue)
	})
}
