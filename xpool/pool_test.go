package xpool

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPool(t *testing.T) {
	Convey("pool should work ok", t, func() {
		var nexStr = "new string"
		var p = NewPool[string](func() string {
			return nexStr
		})
		// 空池子 Get 走 New() factory 返回 nexStr
		var v = p.Get()
		So(v, ShouldEqual, nexStr)

		// Put 之后再 Get：sync.Pool 不保证 Put-then-Get 拿到刚 Put 的值
		// （per-P localPool + GC 可能清 victim cache + goroutine 调度切 P
		// 都会让 Get 走 New() 而非拿到 Put 的值）。这里只验证 Put / Get
		// 不 panic，不断言具体值。
		// 历史 bug: 原断言 'ShouldNotEqual nexStr' 在 CI Linux 多核环境
		// 偶发触发 sync.Pool 实际行为暴露，让测试 fail 但 Pool 本身行为
		// 正确（薄 wrapper 透传 sync.Pool）。
		p.Put("custom")
		v = p.Get()
		// v 可能是 "custom"（拿回刚 Put 的）或 nexStr（factory 兜底），
		// 两种都是合法的 sync.Pool 行为
		So(v, ShouldBeIn, []string{"custom", nexStr})
	})
}
