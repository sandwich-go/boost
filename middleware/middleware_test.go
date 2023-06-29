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
