package xtime

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// periodicTestInterval 是测试中周期任务的统一间隔。
// 取 200ms 是 std time.AfterFunc 调度抖动 + goroutine 启动开销的合理上界，
// 让测试断言"f 大约触发 N 次"有足够容差。
const periodicTestInterval = 200 * time.Millisecond

// TestPeriodic_Basic 验证 Periodic 周期性触发 f。
func TestPeriodic_Basic(t *testing.T) {
	var counter atomic.Int32
	h := Periodic(periodicTestInterval, 0, func() {
		counter.Add(1)
	})
	defer h.Stop()

	// 跑 5.5 个间隔，期望 4..7 次（容差吸收调度抖动）
	time.Sleep(time.Duration(float64(periodicTestInterval) * 5.5))
	got := counter.Load()
	if got < 4 || got > 7 {
		t.Fatalf("expected 4..7 ticks, got %d", got)
	}
}

// TestPeriodic_Stop 验证 Stop 能停止后续触发。
func TestPeriodic_Stop(t *testing.T) {
	var counter atomic.Int32
	h := Periodic(periodicTestInterval, 0, func() {
		counter.Add(1)
	})

	// 跑 ~2.5 个间隔后 Stop
	time.Sleep(time.Duration(float64(periodicTestInterval) * 2.5))
	h.Stop()
	beforeStop := counter.Load()
	// 多给 ~2 个间隔等可能在飞的 fire；Stop 之后下一轮一定不再触发
	time.Sleep(2 * periodicTestInterval)
	afterStop := counter.Load()

	if delta := afterStop - beforeStop; delta > 1 {
		t.Fatalf("Stop did not stop scheduling: %d more ticks after Stop", delta)
	}
}

// TestPeriodic_StopIdempotent 验证 Stop 重复调用安全。
func TestPeriodic_StopIdempotent(t *testing.T) {
	h := Periodic(time.Hour, 0, func() {})
	h.Stop()
	h.Stop() // must not panic / hang
	h.Stop()
}

// TestPeriodic_Jitter 验证 jitter 真的让间隔随机化。
func TestPeriodic_Jitter(t *testing.T) {
	var mu sync.Mutex
	var fires []time.Time

	h := Periodic(periodicTestInterval, 30, func() {
		mu.Lock()
		fires = append(fires, time.Now())
		mu.Unlock()
	})
	defer h.Stop()

	// 跑 8 个间隔，争取拿到 6+ 个采样点
	time.Sleep(8 * periodicTestInterval)
	mu.Lock()
	defer mu.Unlock()

	if len(fires) < 5 {
		t.Fatalf("not enough fires to test jitter: %d", len(fires))
	}

	// 抖动应该让连续两次间隔有差异
	intervals := make([]time.Duration, len(fires)-1)
	var minI, maxI time.Duration = time.Hour, 0
	for i := 1; i < len(fires); i++ {
		intervals[i-1] = fires[i].Sub(fires[i-1])
		if intervals[i-1] < minI {
			minI = intervals[i-1]
		}
		if intervals[i-1] > maxI {
			maxI = intervals[i-1]
		}
	}
	// 30% 抖动 200ms = ±60ms，极差应远大于 5ms
	if maxI-minI < 5*time.Millisecond {
		t.Fatalf("jitter not effective: intervals min=%v max=%v all=%v", minI, maxI, intervals)
	}
}

// TestPeriodic_Shutdown 验证 shutdown chan 关闭后停止。
func TestPeriodic_Shutdown(t *testing.T) {
	shutdown := make(chan struct{})
	var counter atomic.Int32
	PeriodicWithShutdown(periodicTestInterval, 0, shutdown, func() {
		counter.Add(1)
	})

	time.Sleep(time.Duration(float64(periodicTestInterval) * 2.5))
	close(shutdown)
	beforeShutdown := counter.Load()
	time.Sleep(2 * periodicTestInterval)
	afterShutdown := counter.Load()

	if delta := afterShutdown - beforeShutdown; delta > 1 {
		t.Fatalf("shutdown chan did not stop scheduling: %d more ticks", delta)
	}
}

// TestPeriodic_NoOverlap 验证同一 Periodic 的 f 不会并发执行——
// 实现保证：reschedule 在 f 返回后才发生（periodic_std.go.onFire），所以
// 同一 Periodic 的 f 必然串行。
func TestPeriodic_NoOverlap(t *testing.T) {
	var concurrent atomic.Int32
	var maxConcurrent atomic.Int32
	var counter atomic.Int32

	// f 跑得比 interval 慢，制造潜在重叠
	slowFn := periodicTestInterval + 50*time.Millisecond

	h := Periodic(periodicTestInterval, 0, func() {
		now := concurrent.Add(1)
		defer concurrent.Add(-1)
		time.Sleep(slowFn)
		counter.Add(1)
		for {
			old := maxConcurrent.Load()
			if now <= old || maxConcurrent.CompareAndSwap(old, now) {
				break
			}
		}
	})
	defer h.Stop()

	time.Sleep(4 * periodicTestInterval)

	if got := maxConcurrent.Load(); got > 1 {
		t.Fatalf("Periodic overlapped: max concurrent=%d", got)
	}
	if counter.Load() == 0 {
		t.Fatal("no fires observed")
	}
}

// TestPeriodic_NotInCallerGoroutine 验证 f() 不会在 caller goroutine
// 上执行——这是 Periodic 与 Clean 的契约：Clean/Periodic 调用立即返回，
// 后续每轮 f() 在 timer 派发的独立 goroutine 上跑，不阻塞 caller。
//
// 验证手法：让 f() 阻塞等待主 goroutine 的信号；如果 f 真在 caller
// goroutine 跑，主 goroutine 永远拿不到 fired 信号，10s 超时。
func TestPeriodic_NotInCallerGoroutine(t *testing.T) {
	fired := make(chan struct{}, 1)
	release := make(chan struct{})

	h := Periodic(periodicTestInterval, 0, func() {
		select {
		case fired <- struct{}{}:
		default:
		}
		<-release // 故意阻塞，证明 caller 没在等我
	})

	// 等 f 触发；如果 f 跑在 caller goroutine，这里永远收不到——但 Periodic
	// 内部 goroutine 不是 caller 的，所以 release 没关也能 fire。
	select {
	case <-fired:
	case <-time.After(2 * periodicTestInterval):
		close(release)
		h.Stop()
		t.Fatal("f did not fire within 2 intervals")
	}

	// 把 f 放出来，避免 goroutine 泄漏；Stop 阻止后续 fire
	h.Stop()
	close(release)
	// 多等一会儿排空，避免下一轮回调被 Stop 截断后的清理影响后续测试
	time.Sleep(periodicTestInterval)
}

// TestPeriodic_PanicOnInvalid 验证非法参数 panic。
func TestPeriodic_PanicOnInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on d <= 0")
		}
	}()
	Periodic(0, 0, func() {})
}

func TestPeriodic_PanicOnNilFunc(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on nil f")
		}
	}()
	Periodic(time.Second, 0, nil)
}

// BenchmarkPeriodicSteady 量化稳态 alloc：注册 N 个 Periodic 任务，
// 跑足够长的稳态窗口后统计「每次回调触发 / 每秒」的 alloc。
//
// 用 -benchtime=1x 跑：测试逻辑里固定 sample 时长，b.N 仅作为外层 loop。
//
// 期望：稳态接近 0 B/exec —— 每轮 *time.Timer.Reset 不分配，
// 整个 Periodic 生命周期只 alloc 一个底层 timer（构造时一次性）。
func BenchmarkPeriodicSteady(b *testing.B) {
	const (
		taskCount = 200
		warmup    = 300 * time.Millisecond
		sample    = 2 * time.Second
		interval  = 200 * time.Millisecond
	)

	for i := 0; i < b.N; i++ {
		var counter atomic.Int64
		handles := make([]PeriodicHandle, 0, taskCount)
		for j := 0; j < taskCount; j++ {
			h := Periodic(interval, 0, func() { counter.Add(1) })
			handles = append(handles, h)
		}
		time.Sleep(warmup)

		var before, after struct {
			alloc, mallocs uint64
			gc             uint32
		}
		readMem := func(out *struct {
			alloc, mallocs uint64
			gc             uint32
		}) {
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			out.alloc = ms.TotalAlloc
			out.mallocs = ms.Mallocs
			out.gc = ms.NumGC
		}
		startCount := counter.Load()
		readMem(&before)
		start := time.Now()

		time.Sleep(sample)

		readMem(&after)
		elapsed := time.Since(start)
		executions := counter.Load() - startCount

		for _, h := range handles {
			h.Stop()
		}

		if executions == 0 {
			b.Fatal("no executions in steady window")
		}
		b.ReportMetric(float64(after.alloc-before.alloc)/float64(executions), "B/exec")
		b.ReportMetric(float64(after.mallocs-before.mallocs)/float64(executions), "allocs/exec")
		b.ReportMetric(float64(executions)/elapsed.Seconds(), "exec/s")
		b.ReportMetric(float64(after.gc-before.gc), "gc-cycles")
		b.ReportMetric(float64(after.alloc-before.alloc)/elapsed.Seconds()/1024, "KB/s")
	}
}
