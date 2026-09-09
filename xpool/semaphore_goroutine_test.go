package xpool

import (
	"context"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// closeAndWait 关闭并等待收干，带超时上限，确保任何用例都不会挂住测试进程。
func closeAndWait(t *testing.T, p *SemaphoreGoroutinePool) {
	t.Helper()
	p.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := p.Wait(ctx); err != nil {
		t.Errorf("Wait 未在 2s 内收干: %v（可能有 Job 未结束或额度泄漏）", err)
	}
}

// TestSemaphorePoolNeverExceedsCap 锁住最关键的承诺：同时执行的 Job 数永不超过上限。
// 提交量远大于上限，用 atomic 记录并发峰值。
func TestSemaphorePoolNeverExceedsCap(t *testing.T) {
	const limit = 8
	const jobs = 2000
	p := NewSemaphoreGoroutinePool(limit, 0)
	defer closeAndWait(t, p)

	var cur, peak int64
	var wg sync.WaitGroup
	wg.Add(jobs)
	for i := 0; i < jobs; i++ {
		err := p.Push(context.Background(), func() {
			defer wg.Done()
			n := atomic.AddInt64(&cur, 1)
			for {
				old := atomic.LoadInt64(&peak)
				if n <= old || atomic.CompareAndSwapInt64(&peak, old, n) {
					break
				}
			}
			runtime.Gosched()
			atomic.AddInt64(&cur, -1)
		})
		if err != nil {
			t.Fatalf("job %d push failed: %v", i, err)
		}
	}
	wg.Wait()

	if got := atomic.LoadInt64(&peak); got > limit {
		t.Fatalf("并发峰值 %d 超过上限 %d", got, limit)
	}
	if got := atomic.LoadInt64(&peak); got < 2 {
		t.Fatalf("并发峰值 %d，未真正并行，测试没有验证到上限", got)
	}
}

// TestSemaphorePoolIdleCostsNoGoroutine 锁住"空闲零常驻"：构造后不应多出协程。
func TestSemaphorePoolIdleCostsNoGoroutine(t *testing.T) {
	runtime.GC()
	before := runtime.NumGoroutine()
	p := NewSemaphoreGoroutinePool(8000, 0)
	defer closeAndWait(t, p)
	if got := runtime.NumGoroutine(); got != before {
		t.Fatalf("构造 8000 额度的池后协程数 %d -> %d，应无变化", before, got)
	}
	if p.Cap() != 8000 || p.Running() != 0 {
		t.Fatalf("Cap=%d Running=%d，期望 8000/0", p.Cap(), p.Running())
	}
}

// TestSemaphorePoolPushErrorMeansJobNotRun 锁住语义承诺：Push 返回错误时 Job 一定未执行。
// 用一个 Job 占满容量 1 的池，再以带 deadline 的 ctx 提交第二个。
// 这里走 ctx 而非 timeout，是因为 timeout 受秒级时间轮下限约束（见 TimeoutClampedByWheel）。
func TestSemaphorePoolPushErrorMeansJobNotRun(t *testing.T) {
	p := NewSemaphoreGoroutinePool(1, 0)
	block := make(chan struct{})
	entered := make(chan struct{})
	if err := p.Push(context.Background(), func() {
		close(entered)
		<-block
	}); err != nil {
		t.Fatalf("首个 job push 失败: %v", err)
	}
	<-entered

	var ran int32
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := p.Push(ctx, func() { atomic.AddInt32(&ran, 1) }); err == nil {
		t.Fatal("额度耗尽时 Push 应返回错误")
	}
	time.Sleep(20 * time.Millisecond) // 给"万一被执行"留出时间
	if got := atomic.LoadInt32(&ran); got != 0 {
		t.Fatalf("Push 报错但 job 执行了 %d 次，违反语义承诺", got)
	}

	close(block)
	closeAndWait(t, p)
}

// TestSemaphorePoolTimeoutClampedByWheel 锁住一个容易踩的框架行为：
// timeout 走包内秒级时间轮，xtime.Wheel.After 把小于 1s 的值抬到 1s。
// 传 1ms 并不会在 1ms 返回，而是约 1s。GoroutinePool 同样受此约束。
func TestSemaphorePoolTimeoutClampedByWheel(t *testing.T) {
	p := NewSemaphoreGoroutinePool(1, time.Millisecond)
	block := make(chan struct{})
	entered := make(chan struct{})
	if err := p.Push(context.Background(), func() { close(entered); <-block }); err != nil {
		t.Fatalf("首个 job push 失败: %v", err)
	}
	<-entered

	start := time.Now()
	if err := p.Push(context.Background(), func() {}); err == nil {
		t.Fatal("额度耗尽时 Push 应返回错误")
	}
	if cost := time.Since(start); cost < 500*time.Millisecond {
		t.Fatalf("耗时 %v，说明 1ms 未被时间轮抬到秒级；"+
			"若 xtime.Wheel 改了下限，请同步本用例与 NewSemaphoreGoroutinePool 的文档", cost)
	}

	close(block)
	closeAndWait(t, p)
}

// TestSemaphorePoolPushRespectsContext 额度耗尽且未设 timeout 时，ctx 取消应能解除等待。
func TestSemaphorePoolPushRespectsContext(t *testing.T) {
	p := NewSemaphoreGoroutinePool(1, 0)
	block := make(chan struct{})
	entered := make(chan struct{})
	if err := p.Push(context.Background(), func() { close(entered); <-block }); err != nil {
		t.Fatalf("首个 job push 失败: %v", err)
	}
	<-entered

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- p.Push(ctx, func() {}) }()
	select {
	case err := <-done:
		t.Fatalf("额度耗尽时 Push 不应立即返回，得到 %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("ctx 取消后 Push 应返回错误")
		}
	case <-time.After(time.Second):
		t.Fatal("ctx 取消后 Push 未返回")
	}

	close(block)
	closeAndWait(t, p)
}

// TestSemaphorePoolCloseIsNonBlockingAndWaitDrains 锁住 Close/Wait 的分工：
// Close 立即返回且拒绝新 Push；Wait 才等在飞 Job 收干。
func TestSemaphorePoolCloseIsNonBlockingAndWaitDrains(t *testing.T) {
	const n = 16
	p := NewSemaphoreGoroutinePool(n, 0)
	var done int32
	for i := 0; i < n; i++ {
		if err := p.Push(context.Background(), func() {
			time.Sleep(30 * time.Millisecond)
			atomic.AddInt32(&done, 1)
		}); err != nil {
			t.Fatalf("push 失败: %v", err)
		}
	}

	start := time.Now()
	p.Close()
	if cost := time.Since(start); cost > 10*time.Millisecond {
		t.Fatalf("Close 耗时 %v，应立即返回而不等在飞任务", cost)
	}
	if !p.IsClosed() {
		t.Fatal("Close 后 IsClosed 应为 true")
	}
	if err := p.Push(context.Background(), func() {}); err == nil {
		t.Fatal("Close 后 Push 应返回错误")
	}
	p.Close() // 重复 Close 应为 no-op

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := p.Wait(ctx); err != nil {
		t.Fatalf("Wait 未收干: %v", err)
	}
	if got := atomic.LoadInt32(&done); got != n {
		t.Fatalf("Wait 返回时完成 %d/%d 个 job", got, n)
	}
}

// TestSemaphorePoolWaitRespectsContext 锁住 Wait 的逃生口：Job 卡住时 Wait 按 ctx 超时返回，
// 不会永久阻塞。这是把等待从 Close 里拆出来的原因。
func TestSemaphorePoolWaitRespectsContext(t *testing.T) {
	p := NewSemaphoreGoroutinePool(2, 0)
	block := make(chan struct{})
	entered := make(chan struct{})
	if err := p.Push(context.Background(), func() { close(entered); <-block }); err != nil {
		t.Fatalf("push 失败: %v", err)
	}
	<-entered
	p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	if err := p.Wait(ctx); err == nil {
		t.Fatal("Job 未结束时 Wait 应返回 ctx 错误")
	}
	if cost := time.Since(start); cost > time.Second {
		t.Fatalf("Wait 耗时 %v，未按 ctx 及时返回", cost)
	}

	// Wait 超时后应把已占额度还回，池状态不被破坏
	if n := p.Running(); n != 1 {
		t.Fatalf("Running=%d，期望 1（仅卡住的那个 Job 占用）", n)
	}

	close(block)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	if err := p.Wait(ctx2); err != nil {
		t.Fatalf("Job 结束后 Wait 应收干: %v", err)
	}
}

// TestSemaphorePoolRejectsBadArgs 校验入参防御。
func TestSemaphorePoolRejectsBadArgs(t *testing.T) {
	t.Run("nil job", func(t *testing.T) {
		p := NewSemaphoreGoroutinePool(1, 0)
		defer closeAndWait(t, p)
		if err := p.Push(context.Background(), nil); err == nil {
			t.Fatal("nil job 应返回错误")
		}
	})
	t.Run("nil ctx", func(t *testing.T) {
		p := NewSemaphoreGoroutinePool(1, 0)
		defer closeAndWait(t, p)
		ran := make(chan struct{})
		//nolint:staticcheck // 显式验证 nil ctx 被兜住而不是 panic
		if err := p.Push(nil, func() { close(ran) }); err != nil {
			t.Fatalf("nil ctx 应被兜住，得到 %v", err)
		}
		select {
		case <-ran:
		case <-time.After(time.Second):
			t.Fatal("job 未执行")
		}
	})
	t.Run("非正上限应 panic", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("maxConcurrency<=0 应 panic")
			}
		}()
		NewSemaphoreGoroutinePool(0, 0)
	})
}

// TestSemaphorePoolPanicUsesBusinessHandler 锁住 WithSemaphorePoolOnPanic 的两条承诺：
// 1) 业务处理函数拿到 panic 值；2) 额度已归还，池仍可用（panic 不泄漏额度、不削弱容量）。
func TestSemaphorePoolPanicUsesBusinessHandler(t *testing.T) {
	got := make(chan interface{}, 1)
	p := NewSemaphoreGoroutinePool(1, 0, WithSemaphorePoolOnPanic(func(reason interface{}) {
		got <- reason
	}))
	defer closeAndWait(t, p)

	if err := p.Push(context.Background(), func() { panic("boom") }); err != nil {
		t.Fatalf("push 失败: %v", err)
	}
	select {
	case r := <-got:
		if r != "boom" {
			t.Fatalf("onPanic 收到 %v，期望 boom", r)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("业务 panic 处理函数未被调用")
	}

	// 容量为 1：若上一个 job panic 后额度没归还，这里将拿不到额度
	ran := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := p.Push(ctx, func() { close(ran) }); err != nil {
		t.Fatalf("panic 后额度未归还，后续 Push 失败: %v", err)
	}
	select {
	case <-ran:
	case <-time.After(2 * time.Second):
		t.Fatal("panic 后池不可用")
	}
	if n := p.Running(); n != 0 {
		t.Fatalf("Running=%d，期望 0（额度应已全部归还）", n)
	}
}

// TestSemaphorePoolPanicDefaultRecovers 锁住默认语义：不传 WithSemaphorePoolOnPanic 时
// Job panic 仍被 recover，进程不受影响，池可继续使用。
//
// 本用例真实触发 panic：若默认改回"不 recover"，测试进程会被带走，直接暴露为失败。
func TestSemaphorePoolPanicDefaultRecovers(t *testing.T) {
	p := NewSemaphoreGoroutinePool(1, 0)
	defer closeAndWait(t, p)

	if err := p.Push(context.Background(), func() { panic("boom-default") }); err != nil {
		t.Fatalf("push 失败: %v", err)
	}

	// 容量 1：只有前一个 Job 的 panic 被 recover 且额度归还，这个才能拿到额度
	ran := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := p.Push(ctx, func() { close(ran) }); err != nil {
		t.Fatalf("默认 panic 处理后池不可用: %v", err)
	}
	select {
	case <-ran:
	case <-time.After(2 * time.Second):
		t.Fatal("后续 Job 未执行")
	}
}

// TestSemaphorePoolPanicDoesNotStrandWaiter 锁住文档里那条约定的实际效果：
// Job 把完成信号写成 defer 时，panic 被 recover 后等待方不会被永久阻塞。
// 这正是"扇出后 wg.Wait()"用法依赖的性质。
func TestSemaphorePoolPanicDoesNotStrandWaiter(t *testing.T) {
	p := NewSemaphoreGoroutinePool(4, 0)
	defer closeAndWait(t, p)

	var wg sync.WaitGroup
	wg.Add(2)
	if err := p.Push(context.Background(), func() { defer wg.Done(); panic("boom-waiter") }); err != nil {
		t.Fatalf("push 失败: %v", err)
	}
	if err := p.Push(context.Background(), func() { defer wg.Done() }); err != nil {
		t.Fatalf("push 失败: %v", err)
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("panic 的 Job 未释放 wg，等待方被永久阻塞")
	}
}

// BenchmarkSemaphorePoolFanout 对比 GoroutinePool：模拟"扇出 N 个任务后等齐"的形态。
//
// 跑法：go test -run '^$' -bench Fanout -benchmem ./xpool/
func BenchmarkSemaphorePoolFanout(b *testing.B) {
	for _, fan := range []int{1, 4, 32} {
		name := strconv.Itoa(fan)
		b.Run("semaphore/fan"+name, func(b *testing.B) {
			p := NewSemaphoreGoroutinePool(8000, time.Second)
			defer p.Close()
			benchFanout(b, fan, p.Push)
		})
		b.Run("goroutinepool/fan"+name, func(b *testing.B) {
			p := NewGoroutinePool(8000, 1000000, time.Second)
			defer p.Close()
			benchFanout(b, fan, p.Push)
		})
	}
}

func benchFanout(b *testing.B, fan int, push func(context.Context, Job) error) {
	ctx := context.Background()
	var sink int64
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(fan)
		for j := 0; j < fan; j++ {
			if err := push(ctx, func() { defer wg.Done(); atomic.AddInt64(&sink, 1) }); err != nil {
				wg.Done()
			}
		}
		wg.Wait()
	}
}

// TestSemaphorePoolReserveHoldsSlotBeforeRun 锁住 Reserve 的核心承诺：返回成功即额度在手，
// 此时还没派发任何 Job，但容量已被占住。
//
// 这是 Reserve 存在的全部理由——调用方要在派发前取「不可回收」的资源（例如全局单调序号），
// 必须先确认派发一定会发生。若 Reserve 没真占住额度，后续 run 仍可能因额度不足而阻塞或失败，
// 资源就漏在外面且无从补偿。
func TestSemaphorePoolReserveHoldsSlotBeforeRun(t *testing.T) {
	p := NewSemaphoreGoroutinePool(1, 0)
	defer closeAndWait(t, p)

	run, err := p.Reserve(context.Background())
	if err != nil {
		t.Fatalf("Reserve 失败: %v", err)
	}

	// 唯一的额度已被预约但尚未派发，此刻再 Push 必须拿不到额度。
	// 用短 ctx 断言"拿不到"，而不是断言 Running()——Running 统计的是在飞 Job 数，
	// 预约还没派发时它本就是 0，用它区分不出"额度被占住"和"额度还空着"。
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err = p.Push(ctx, func() {}); err == nil {
		t.Error("Reserve 之后额度仍可被别人拿走：预约没有真正占住容量，" +
			"调用方在 run 之前取的不可回收资源可能白取")
	}

	var ran atomic.Bool
	done := make(chan struct{})
	run(func() { ran.Store(true); close(done) })

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("run 提交的 Job 未在 2s 内执行")
	}
	if !ran.Load() {
		t.Error("run 没有执行 Job")
	}
}

// TestSemaphorePoolReserveNilJobReleasesSlot 锁住放弃预约时额度被归还。
//
// 调用方在 Reserve 成功之后仍可能决定不发（例如发现没有可发的内容），此时必须能把额度还回去，
// 否则每次放弃都让池容量永久缩水一格。
func TestSemaphorePoolReserveNilJobReleasesSlot(t *testing.T) {
	p := NewSemaphoreGoroutinePool(1, 0)
	defer closeAndWait(t, p)

	run, err := p.Reserve(context.Background())
	if err != nil {
		t.Fatalf("Reserve 失败: %v", err)
	}
	run(nil) // 放弃

	// 额度归还后，Push 应立即成功。
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan struct{})
	if err = p.Push(ctx, func() { close(done) }); err != nil {
		t.Fatalf("放弃预约后额度没被归还，Push 失败: %v；每次放弃都会让容量缩水一格", err)
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Job 未在 2s 内执行")
	}
}

// TestSemaphorePoolReserveRunIsIdempotent 锁住 run 重复调用不会多还额度。
//
// 多还一次就凭空多出一个额度，池的并发上限被撑破——而上限本身就是这个池唯一的承诺。
// 误用（重复调用）应当是安全的空操作，而不是静默破坏容量。
func TestSemaphorePoolReserveRunIsIdempotent(t *testing.T) {
	const limit = 2
	p := NewSemaphoreGoroutinePool(limit, 0)
	defer closeAndWait(t, p)

	run, err := p.Reserve(context.Background())
	if err != nil {
		t.Fatalf("Reserve 失败: %v", err)
	}

	block := make(chan struct{})
	var running atomic.Int32
	var peak atomic.Int32
	job := func() {
		cur := running.Add(1)
		for {
			old := peak.Load()
			if cur <= old || peak.CompareAndSwap(old, cur) {
				break
			}
		}
		<-block
		running.Add(-1)
	}

	run(job)
	run(job) // 重复调用：必须是空操作
	run(nil) // 放弃：同样必须是空操作

	// 再填满剩余额度；若上面的重复调用多还了额度，这里能塞进超过 limit 个并发 Job。
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for i := 0; i < limit; i++ {
		_ = p.Push(ctx, job)
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && running.Load() < int32(limit) {
		runtime.Gosched()
	}
	close(block)

	if got := peak.Load(); got > int32(limit) {
		t.Errorf("并发峰值 %d 超过上限 %d：run 被重复调用后多归还了额度，池的限流承诺被打破",
			got, limit)
	}
}

// TestSemaphorePoolReserveRespectsContextAndClose 锁住 Reserve 的失败路径不占额度。
//
// 失败时若把额度留在手里，池容量会随每次失败递减；而失败恰恰发生在关停与超时这类本就异常的
// 时刻，容量泄漏会让问题雪上加霜。
func TestSemaphorePoolReserveRespectsContextAndClose(t *testing.T) {
	t.Run("ctx 取消时失败且不占额度", func(t *testing.T) {
		p := NewSemaphoreGoroutinePool(1, 0)
		defer closeAndWait(t, p)

		// 先占满唯一额度
		block := make(chan struct{})
		if err := p.Push(context.Background(), func() { <-block }); err != nil {
			t.Fatalf("前置 Push 失败: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		if run, err := p.Reserve(ctx); err == nil {
			run(nil)
			t.Error("额度已满且 ctx 到期，Reserve 仍返回成功")
		}

		close(block)
	})

	t.Run("池关闭后失败", func(t *testing.T) {
		p := NewSemaphoreGoroutinePool(1, 0)
		p.Close()
		if run, err := p.Reserve(context.Background()); err == nil {
			run(nil)
			t.Error("池已关闭，Reserve 仍返回成功：关停期间不该再占额度起新工作")
		}
	})
}

// TestSemaphorePoolPushStillWorksViaReserve 锁住 Push 改为复用 Reserve 之后行为不变。
//
// Push 是既有 API，绝大多数调用方在用；拆出 Reserve 时若把 nil job 校验或错误语义弄丢，
// 影响面远大于新 API 本身。
func TestSemaphorePoolPushStillWorksViaReserve(t *testing.T) {
	p := NewSemaphoreGoroutinePool(2, 0)
	defer closeAndWait(t, p)

	if err := p.Push(context.Background(), nil); err == nil {
		t.Error("Push(nil) 应当报错：拆分后不能把 nil job 校验丢掉")
	}

	done := make(chan struct{})
	if err := p.Push(context.Background(), func() { close(done) }); err != nil {
		t.Fatalf("Push 失败: %v", err)
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Push 的 Job 未在 2s 内执行")
	}
}

// TestSemaphorePoolReserveFailsFastWhenClosed 锁住池已关闭时 Reserve 立即失败，而不是先去等额度。
//
// Reserve 里有两道 IsClosed 检查：入口一道、拿到额度后一道。后者处理「等额度期间池被关闭」的
// 竞态，是正确性所必需；前者只影响快慢——删掉它结果依然正确（第二道会拦住），所以用例必须
// 断言「耗时」而不只是「失败」，否则这道检查删掉也没人知道。
//
// 快速失败的意义：额度被长时间占用的 Job 占满时，入口不判就会一直等到那些 Job 结束或 ctx 取消。
// 关停路径上多等这一会儿会拖慢整个进程收干。
func TestSemaphorePoolReserveFailsFastWhenClosed(t *testing.T) {
	p := NewSemaphoreGoroutinePool(1, 0)

	// 占满唯一额度，且让它在用例结束前不释放。
	block := make(chan struct{})
	if err := p.Push(context.Background(), func() { <-block }); err != nil {
		t.Fatalf("前置 Push 失败: %v", err)
	}
	p.Close() // Close 本身不阻塞

	start := time.Now()
	run, err := p.Reserve(context.Background())
	cost := time.Since(start)
	if err == nil {
		run(nil)
		t.Fatal("池已关闭，Reserve 仍返回成功")
	}
	// 阈值取一个宽松值：入口判 IsClosed 是纯内存操作，正常在微秒级；
	// 而少了它就要等上面那个 Job 结束（本用例里它永不结束，只能等到超时）。
	if cost > 200*time.Millisecond {
		t.Errorf("Reserve 在池关闭后耗时 %v 才失败：入口的 IsClosed 检查没生效，"+
			"关停时会白等额度，拖慢进程收干", cost)
	}

	close(block)
	closeAndWait(t, p)
}
