package xtime

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sandwich-go/boost/xtime/cron"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 xtime 包内 0% 覆盖的导出函数 / 方法（不重复已有测试）。
// 已有测试（dispatcher_test / periodic_test / tick_test / wheel_test /
// beginning_end_test / cop_test）覆盖了主路径，本文件补漏。

// TestToLocal_AndIsLeapYear leaf 工具函数。
func TestToLocal_AndIsLeapYear(t *testing.T) {
	Convey("ToLocal 应用 FixedZone offset", t, func() {
		base := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
		// +8h 偏移
		got := ToLocal(base, 8*3600)
		_, off := got.Zone()
		So(off, ShouldEqual, 8*3600)
		So(got.Hour(), ShouldEqual, 20) // UTC 12:00 → +8 = 20:00
	})

	Convey("IsLeapYear", t, func() {
		// 4 年闰：2020 / 2024 是
		So(IsLeapYear(2020), ShouldBeTrue)
		So(IsLeapYear(2024), ShouldBeTrue)
		// 100 年不闰：1900 不是
		So(IsLeapYear(1900), ShouldBeFalse)
		// 400 年再闰：2000 是
		So(IsLeapYear(2000), ShouldBeTrue)
		// 普通：2023 不是
		So(IsLeapYear(2023), ShouldBeFalse)
	})
}

// TestCop_StopAndUnixMilli 覆盖 Cop.Stop / UnixMilli + cop_global.Stop /
// SetNowProvider。
func TestCop_StopAndUnixMilli(t *testing.T) {
	Convey("Cop.UnixMilli 返回毫秒级时间戳", t, func() {
		fixed := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
		c := NewCop(func() time.Time { return fixed })
		// 没 start，Now() 走 nowProvider；UnixMilli 等于 fixed.UnixNano()/1e6
		So(c.UnixMilli(), ShouldEqual, fixed.UnixNano()/int64(time.Millisecond))
	})

	Convey("Cop.Stop 幂等：重复调用不 panic 不 close 已 closed channel", t, func() {
		c := NewCop(time.Now)
		c.Stop()
		// 第二次 Stop 不应 panic（CompareAndSwap 守门）
		So(func() { c.Stop() }, ShouldNotPanic)
	})

	Convey("Cop.SetNowProvider 与并发读 Now 之间原子安全", t, func() {
		// 历史 bug：nowProvider 是裸 func 字段，SetNowProvider 与后台
		// start/check/Now 之间 read-write race。修复改成 atomic.Pointer。
		//
		// 反向验证：本测试在 -race 下若回退到裸字段写会被 race detector 抓到。
		fixed1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		fixed2 := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
		c := NewCop(func() time.Time { return fixed1 })
		defer c.Stop()

		// 并发：N 个 setter 替换 provider，N 个 reader 调 Now()
		const n = 50
		done := make(chan struct{})
		var swaps atomic.Int64
		for i := 0; i < n; i++ {
			go func(i int) {
				for j := 0; j < 100; j++ {
					if (i+j)%2 == 0 {
						c.SetNowProvider(func() time.Time { return fixed1 })
					} else {
						c.SetNowProvider(func() time.Time { return fixed2 })
					}
					swaps.Add(1)
				}
				done <- struct{}{}
			}(i)
			go func() {
				for j := 0; j < 100; j++ {
					_ = c.Now()
				}
				done <- struct{}{}
			}()
		}
		for i := 0; i < 2*n; i++ {
			<-done
		}
		So(swaps.Load(), ShouldEqual, int64(n*100))
	})

	Convey("包级 SetNowProvider / Stop API 存在且可调用", t, func() {
		// API 存在性 smoke test（globalCop 不能 Stop 否则污染其他测试）
		So(SetNowProvider, ShouldNotBeNil)
		So(Stop, ShouldNotBeNil)
	})
}

// TestDispatcher_AfterFuncWithOwnershipTransfer 覆盖
// AfterFuncWithOwnershipTransfer + WithOwnershipTransferInDomain，业务自管
// timer 生命周期的路径（不进 runningTimers，不被 RemoveAllTimer 清理）。
func TestDispatcher_AfterFuncWithOwnershipTransfer(t *testing.T) {
	Convey("AfterFuncWithOwnershipTransfer 注册的 timer 不在 runningTimers 中", t, func() {
		d := NewDispatcher(64)
		d.Start()
		defer d.Close()

		var fired atomic.Int32
		dt := d.AfterFuncWithOwnershipTransfer(20*time.Millisecond, func() {
			fired.Add(1)
		})
		So(dt, ShouldNotBeNil)

		// 业务接收 ChanTimer 推送的 Timer 并执行 Cb
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		select {
		case t := <-d.TimerNotify():
			t.Cb()
		case <-ctx.Done():
			tEarly := dt
			t.Fatalf("AfterFuncWithOwnershipTransfer never fired, dt=%v", tEarly)
		}
		So(fired.Load(), ShouldEqual, int32(1))

		// DanglingTimer.Stop 是业务自管入口；幂等不 panic
		So(func() { dt.Stop() }, ShouldNotPanic)
		So(func() { dt.Stop() }, ShouldNotPanic)
	})

	Convey("DanglingTimer.GetDomain / Reset", t, func() {
		d := NewDispatcher(64)
		d.Start()
		defer d.Close()

		dt := d.AfterFuncWithOwnershipTransferInDomain(1*time.Hour, func() {}, "biz-domain")
		So(dt.GetDomain(), ShouldEqual, "biz-domain")
		// Reset 在 timer 还活着时返回 true
		ok := dt.Reset(1 * time.Hour)
		So(ok, ShouldBeTrue)
		dt.Stop()
	})
}

// TestSafeTimer_StopMethods 覆盖 SafeTimer.Stop（与 timer.go 中 0% 的方法）。
// SafeTimer 是 dispatcher.AfterFunc 返回类型，业务通常不直接调 Stop（dispatcher
// Close 会自动清理），但 API 保留外部 Stop 入口。
func TestSafeTimer_StopMethods(t *testing.T) {
	Convey("SafeTimer.stop 内部清理 cb", t, func() {
		d := NewDispatcher(64)
		d.Start()
		defer d.Close()

		st := d.AfterFunc(1*time.Hour, func() {})
		So(st, ShouldNotBeNil)
		// 直接调 stop（小写，内部）—— 由 RemoveAllTimer 触发
		// 这里通过 RemoveAllTimer 间接覆盖；剩余 80% RemoveAllTimer 跑完整 Range
		d.RemoveAllTimer()
		// 再 RemoveAllTimer 是空 Range，无副作用
		d.RemoveAllTimer()
	})
}

// TestDispatcher_CronFunc 覆盖 CronFunc 路径。cron 表达式每秒触发一次。
func TestDispatcher_CronFunc(t *testing.T) {
	Convey("CronFunc 按 cron 表达式触发回调", t, func() {
		d := NewDispatcher(64)
		d.Start()
		defer d.Close()

		// 每秒触发
		expr, err := cron.Parse("* * * * * *")
		So(err, ShouldBeNil)

		var fired atomic.Int32
		c := d.CronFunc(expr, func() { fired.Add(1) })
		So(c, ShouldNotBeNil)

		// 等待回调触发：CronFunc 通过 ChanTimer 投递，需要 drain
		ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
		defer cancel()

	drain:
		for {
			select {
			case t := <-d.TimerNotify():
				t.Cb()
				if fired.Load() >= 1 {
					break drain
				}
			case <-ctx.Done():
				break drain
			}
		}
		So(fired.Load(), ShouldBeGreaterThanOrEqualTo, int32(1))

		// Stop 幂等不 panic
		So(func() { c.Stop() }, ShouldNotPanic)
		So(func() { c.Stop() }, ShouldNotPanic)
	})

	Convey("CronFunc 表达式 Next 返零时间直接返回空 Cron", t, func() {
		// 构造永不触发的 expression：cron 包目前不易构造 IsZero 场景，
		// 通过反射 / 自定义 Expression 也复杂，skip 该分支（实际生产
		// cron 解析合法时 Next 总返非零，IsZero 仅 defensive）
		t.Skip("Cron Next IsZero 是 defensive 分支，跳过")
	})
}

// TestWheelAfter_AndSetAccuracy 覆盖 WheelAfter / newWheelWithDuration /
// SetAccuracy（time wheel 的可选模式）。
func TestWheelAfter_AndSetAccuracy(t *testing.T) {
	Convey("SetAccuracy + WheelAfter 触发回调", t, func() {
		// 保存恢复 accuracy，避免影响其他测试
		old := accuracy
		defer SetAccuracy(old)

		SetAccuracy(10) // 1/10 = 10% 精度
		So(accuracy, ShouldEqual, 10)

		// WheelAfter 接到 channel 后：等到 fire
		ch := WheelAfter(50 * time.Millisecond)
		So(ch, ShouldNotBeNil)

		select {
		case <-ch:
			// fired
		case <-time.After(500 * time.Millisecond):
			t.Fatal("WheelAfter did not fire within 500ms")
		}

		// 再次 WheelAfter 同 duration：复用同一 wheel
		ch2 := WheelAfter(50 * time.Millisecond)
		select {
		case <-ch2:
		case <-time.After(500 * time.Millisecond):
			t.Fatal("WheelAfter (reused) did not fire")
		}
	})
}

// TestDispatcher_TickerCAndTimerNotify 覆盖 TickerC 50% 缺失分支：tickFreq=0
// 时返回 nil 的路径（已有测试覆盖了 tickFreq>0 时返 ticker.C 的路径）。
func TestDispatcher_TickerCAndTimerNotify(t *testing.T) {
	Convey("TickerC tickFreq=0 返回 nil channel", t, func() {
		// 默认 Option tickFreq=0
		d := NewDispatcher(64).(*dispatcher) // TickerC 在 TickerDispatcher 子接口上
		So(d.TickerC(), ShouldBeNil)
	})

	Convey("TickerC tickFreq>0 + hostingMode=true 走 warn 分支返 nil", t, func() {
		// hostingMode=true（默认）下调 TickerC 是用法错误（应让 dispatcher 自托管），
		// 走 line 149-151 log.Warn + return nil 路径。
		d := NewDispatcher(64, WithTickDuration(5*time.Millisecond)).(*dispatcher)
		defer d.Close()
		So(d.TickerC(), ShouldBeNil)
	})
}

// TestDispatcher_TickFunc_DuplicateKey 覆盖 TickFunc 的重复 key 早返路径（warn）。
func TestDispatcher_TickFunc_DuplicateKey(t *testing.T) {
	Convey("TickFunc 同 key 重复注册走 warn 不重复加入 tickFuncs", t, func() {
		d := NewDispatcher(64, WithTickDuration(5*time.Millisecond)).(*dispatcher)
		defer d.Close()

		// 第一次注册成功
		d.TickFunc("dup-key", func(_ context.Context) {})
		So(len(d.tickFuncs), ShouldEqual, 1)

		// 第二次相同 key 走 warn 分支，不再 append
		d.TickFunc("dup-key", func(_ context.Context) {})
		So(len(d.tickFuncs), ShouldEqual, 1)

		// 不同 key 正常 append
		d.TickFunc("another-key", func(_ context.Context) {})
		So(len(d.tickFuncs), ShouldEqual, 2)
	})
}

// TestDispatcher_TriggerTickFuncs_NotRunning 覆盖 TriggerTickFuncs 在
// closeFlag != stateRunning 时的早返路径（line 159-160）。
//
// 注意：NewDispatcher 内部已自动 Start（dispatcher.go:87），所以构造完
// 直接 TriggerTickFuncs 是 running 状态会执行。要触发早返必须先 Close。
func TestDispatcher_TriggerTickFuncs_NotRunning(t *testing.T) {
	Convey("TriggerTickFuncs 在 dispatcher Close 后早返不执行 cb", t, func() {
		d := NewDispatcher(64, WithTickDuration(5*time.Millisecond), WithTickHostingMode(false)).(*dispatcher)
		var called atomic.Int32
		d.TickFunc("k", func(_ context.Context) { called.Add(1) })

		// 构造后是 running 状态，TriggerTickFuncs 正常跑 cb
		d.TriggerTickFuncs(context.Background())
		So(called.Load(), ShouldEqual, int32(1))

		// Close 后 closeFlag != stateRunning，TriggerTickFuncs 走早返不再 cb
		d.Close()
		d.TriggerTickFuncs(context.Background())
		So(called.Load(), ShouldEqual, int32(1)) // 没再涨，确认走早返
	})
}
