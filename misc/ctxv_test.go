package misc

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestContextValue(t *testing.T) {
	Convey("ctxv should work ok", t, func() {
		var ctxv = NewContextValue[int]()
		var ctx = context.Background()

		var v = ctxv.Get(ctx)
		So(v, ShouldBeZeroValue)

		ctx = ctxv.With(ctx, 1)
		v = ctxv.Get(ctx)
		So(v, ShouldEqual, 1)
	})
}
