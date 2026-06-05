package xpanic

import (
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestAutoRecover_BasicRecover 验证 panic 被 recover、OnRecover 被调、
// 然后 AutoRecover 起新 goroutine 重试。
//
// 测试设计：让 f 有计数器，前 N 次 panic、第 N+1 次 return；断言 OnRecover
// 被调 N 次，最终 f 至少被执行 N+1 次（panic 后重启起新 goroutine）。
func TestAutoRecover_BasicRecover(t *testing.T) {
	Convey("AutoRecover restarts goroutine after panic until f returns", t, func() {
		var (
			runCount      atomic.Int32
			recoverCount  atomic.Int32
			lastReason    atomic.Value // any
			panicTimes    int32        = 3
			expectedTotal int32        = panicTimes + 1
		)

		done := make(chan struct{})
		f := func() {
			n := runCount.Add(1)
			if n <= panicTimes {
				panic("trial panic")
			}
			close(done)
		}
		onRecover := func(tag string, reason any) {
			recoverCount.Add(1)
			lastReason.Store(reason)
		}

		go AutoRecover("test-tag", f,
			WithAutoRecoverOptionDelayTime(time.Millisecond),
			WithAutoRecoverOptionOnRecover(onRecover),
		)

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatalf("AutoRecover never reached non-panic state: ran=%d recover=%d",
				runCount.Load(), recoverCount.Load())
		}

		// 给 OnRecover 异步完成留窗口（go AutoRecover 是新 goroutine）
		// 实际 OnRecover 在 panic goroutine 的 defer 里执行，已同步，但
		// runCount.Add(1) 比 close(done) 早，最后那次的 OnRecover 不会跑（没 panic）
		So(runCount.Load(), ShouldBeGreaterThanOrEqualTo, expectedTotal)
		So(recoverCount.Load(), ShouldEqual, panicTimes)
		So(lastReason.Load(), ShouldEqual, "trial panic")
	})
}

// TestAutoRecover_NoPanic 验证 f 不 panic 时 AutoRecover 直接返回，OnRecover 不被调。
func TestAutoRecover_NoPanic(t *testing.T) {
	Convey("AutoRecover returns without restart when f does not panic", t, func() {
		var ran atomic.Int32
		var recovered atomic.Int32
		AutoRecover("test", func() {
			ran.Add(1)
		},
			WithAutoRecoverOptionOnRecover(func(string, any) { recovered.Add(1) }),
		)
		So(ran.Load(), ShouldEqual, 1)
		So(recovered.Load(), ShouldEqual, 0)
	})
}

// TestAutoRecover_NilOnRecover 验证 OnRecover 为 nil 时 panic 仍被吞掉、不崩溃。
// 默认 OnRecover 写日志但仓内可注入 nil（看 option.go 默认值）；这里显式覆盖
// 为 nil 验证容错。
func TestAutoRecover_NilOnRecover(t *testing.T) {
	Convey("AutoRecover tolerates nil OnRecover", t, func() {
		var ran atomic.Int32
		done := make(chan struct{})
		go AutoRecover("test-nil", func() {
			n := ran.Add(1)
			if n == 1 {
				panic("once")
			}
			close(done)
		},
			WithAutoRecoverOptionOnRecover(nil),
			WithAutoRecoverOptionDelayTime(time.Millisecond),
		)
		select {
		case <-done:
			So(ran.Load(), ShouldBeGreaterThanOrEqualTo, int32(2))
		case <-time.After(2 * time.Second):
			t.Fatal("AutoRecover with nil OnRecover did not restart")
		}
	})
}
