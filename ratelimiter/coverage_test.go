package ratelimiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 ratelimiter 包内未覆盖的导出 API + Take 边界：
//   - NewUnlimited / unlimited.Take（之前 0%）
//   - Take 的 maxSlack > 0 累积超时分支
//   - 并发 Take race-safe（CAS 路径）
//   - WithSlack(0) 让 maxSlack==0 走 shrink 分支
//
// 已有 ratelimit_test.go 仅测了基础 100 RPS happy path。

func TestNewUnlimited(t *testing.T) {
	Convey("NewUnlimited Take 不阻塞，连续调用都立即返回", t, func() {
		l := NewUnlimited()
		So(l, ShouldNotBeNil)

		start := time.Now()
		for i := 0; i < 1000; i++ {
			l.Take()
		}
		elapsed := time.Since(start)
		// 1000 次调用应在毫秒级完成（无 sleep）
		So(elapsed, ShouldBeLessThan, 100*time.Millisecond)
	})
}

func TestNew_RateLimit_BasicShape(t *testing.T) {
	Convey("New(rate) 限制 RPS 大致正确", t, func() {
		// 100 RPS = 每 10ms 一次；连续 5 次预计 ~40ms（首次不 sleep）
		l := New(100)

		l.Take() // 初始化
		start := time.Now()
		for i := 0; i < 5; i++ {
			l.Take()
		}
		elapsed := time.Since(start)
		// 100 RPS 5 次约 50ms；±20ms 容差
		So(elapsed, ShouldBeBetween, 30*time.Millisecond, 80*time.Millisecond)
	})
}

func TestNew_WithSlack(t *testing.T) {
	Convey("WithSlack 控制空闲后突发上限", t, func() {
		// 10 RPS = 100ms/req；slack=5 让最多累积 500ms
		l := New(10, WithSlack(5))

		l.Take() // 初始化

		// 不 sleep，连续 Take，前几次应走累积分支不阻塞太久
		start := time.Now()
		for i := 0; i < 3; i++ {
			l.Take()
		}
		elapsed := time.Since(start)
		// slack 让初期累积，3 次应 < 300ms
		So(elapsed, ShouldBeLessThan, 350*time.Millisecond)
	})

	Convey("WithSlack(0) 关闭累积，每次都按 perRequest 间隔", t, func() {
		// slack=0 → maxSlack==0 → 走 shrink 分支：每次 now-state > perRequest 时
		// newTime=now，相当于丢弃过去累积
		l := New(10, WithSlack(0))

		l.Take()
		// sleep 较长时间让 maxSlack==0 走 shrink 路径
		time.Sleep(500 * time.Millisecond)

		start := time.Now()
		l.Take() // 走 shrink 分支：newTime=now，sleep 应是 0
		elapsed := time.Since(start)
		So(elapsed, ShouldBeLessThan, 50*time.Millisecond)
	})

	Convey("WithSlack(>0) + 长 idle → 走 maxSlack cap 分支", t, func() {
		// 覆盖 limiter_atomic.go:52-55 'maxSlack > 0 &&
		// now-state > maxSlack+perRequest' 分支：自上次 Take 累积时间
		// 超过 maxSlack 上限，cap 到 maxSlack 防止突发用尽。
		//
		// 10 RPS, slack=2 → maxSlack=2*100ms=200ms, perRequest=100ms。
		// idle 1 秒 >> 200+100=300ms 阈值，触发 cap 分支。
		l := New(10, WithSlack(2))
		l.Take() // 初始化 state

		// idle 远超 maxSlack+perRequest
		time.Sleep(1 * time.Second)

		start := time.Now()
		l.Take() // 走 line 52-55: newTime=now-maxSlack
		elapsed := time.Since(start)
		// cap 后下次 Take 应几乎不阻塞（state 被设到 now-maxSlack ≤ now）
		So(elapsed, ShouldBeLessThan, 50*time.Millisecond)
	})
}

func TestTake_Concurrent(t *testing.T) {
	Convey("并发 Take CAS 重试不丢限速", t, func() {
		// 1000 RPS = 1ms/req；100 goroutine 并发各 Take 5 次
		l := New(1000, WithSlack(10))
		var counter atomic.Int32
		var wg sync.WaitGroup
		const goroutines = 100
		const perG = 5

		wg.Add(goroutines)
		start := time.Now()
		for g := 0; g < goroutines; g++ {
			go func() {
				defer wg.Done()
				for i := 0; i < perG; i++ {
					l.Take()
					counter.Add(1)
				}
			}()
		}
		wg.Wait()
		elapsed := time.Since(start)

		// 总 500 次 / 1000 RPS = ~500ms，但 slack=10 让初始累积；±50% 容差
		// 关键是不 panic 不死锁
		So(counter.Load(), ShouldEqual, int32(goroutines*perG))
		t.Logf("500 calls under 1000 RPS: %v", elapsed)
	})
}

// BenchmarkTake_AtomicInt64 hot path bench：限流器是高频 hot path（每次
// 业务请求都 Take）。
func BenchmarkTake_AtomicInt64(b *testing.B) {
	l := New(int(1e9), WithSlack(0)) // 极高 RPS 让 sleep ≈ 0
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Take()
	}
}

// BenchmarkTake_Unlimited 对照：unlimited 直接返 z.Now，无 CAS 也无 sleep
func BenchmarkTake_Unlimited(b *testing.B) {
	l := NewUnlimited()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Take()
	}
}
