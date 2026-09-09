package xpool

import (
	"context"
	"fmt"
	rtdebug "runtime/debug"
	"sync"
	"time"

	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/xerror"
	"github.com/sandwich-go/boost/xpanic"
	"github.com/sandwich-go/boost/xsync"
)

// SemaphoreGoroutinePool 按需起协程、以计数信号量限流的协程池。
//
// 与 GoroutinePool 的取舍差异：
//
//   - GoroutinePool 预先起满 numWorkers 个 worker 常驻竞争 jobQueue，空闲时这些协程和
//     队列缓冲一直占着内存；本类型不预起协程，每个 Job 由 go 直接派发，空闲成本为零，
//     协程数随负载自然升降。
//   - 两者对并发的约束同样是硬上限：GoroutinePool 上限是 worker 数，本类型上限是信号量
//     额度 maxConcurrency，超出后 Push 排队等额度，绝不会突破上限。
//   - GoroutinePool 有 jobQueueLen 长度的队列可以吸收突发，Push 入队即返回；本类型没有
//     队列，额度耗尽时 Push 直接阻塞。需要"提交方快速脱手、由队列缓冲"的场景仍应用
//     GoroutinePool。
//
// 之所以不再包一层常驻 worker：Go 运行时本身就按 P 缓存空闲 g 的栈（gfree 链表），
// go 派发短任务已是复用路径，再经 channel 转交给常驻 worker 只多一次调度往返。
//
// 语义约定：
//
//   - Push 返回非 nil error 时，job 一定没有被执行。调用方必须处理该错误，不能只记日志
//     就当作已提交——那会变成静默丢任务。
//   - Job panic 一定被 recover，不会终止进程：一个 Job 的 bug 不该带走整个进程和它承载的
//     全部连接。默认记录错误日志与调用栈，业务传 WithSemaphorePoolOnPanic 时改用业务的
//     处理函数。这与 GoroutinePool 不同——后者 job() 裸调，panic 直接终止进程。
//   - Job panic 不会削弱池的容量：额度由 defer 归还；且每个 Job 用独立协程，panic 只影响
//     那一个 Job。GoroutinePool 则是常驻 worker 抢队列，worker 一旦因 panic 退出，
//     p.workers 仍把它计在数内，容量静默缩水且无法恢复。
//   - recover 的代价是 bug 不再以进程崩溃的形式暴露。Job 里请把"完成信号"写成 defer
//     （如 defer wg.Done()），否则 panic 被 recover 后等待方会永久阻塞——不 recover 时
//     进程已经崩了，看不出这个问题。建议调用方在 onPanic 里同时打点计数，别只记日志。
//   - Close 会等待所有在飞 Job 结束后才返回。
type SemaphoreGoroutinePool struct {
	// tokens 是计数信号量：发送占用一个额度，接收归还。
	// 元素类型 struct{} 大小为 0，容量再大也不占缓冲字节。
	tokens    chan struct{}
	timeout   time.Duration
	closeFlag xsync.AtomicInt32
	onPanic   func(reason interface{})
}

// SemaphorePoolOption 构造选项。
type SemaphorePoolOption func(*SemaphoreGoroutinePool)

// WithSemaphorePoolOnPanic 用业务自己的处理函数替换默认的 panic 处理。
//
// 无论是否设定，Job panic 都会被 recover；差别只在拿到 panic 之后做什么：
// 默认记 Error 日志并附调用栈，设定后改调 fn。
//
// fn 在 recover 的 defer 内被调用，因此需要栈时可在 fn 里直接取 rtdebug.Stack()。
// fn 自身 panic 不再被兜住。
func WithSemaphorePoolOnPanic(fn func(reason interface{})) SemaphorePoolOption {
	return func(p *SemaphoreGoroutinePool) { p.onPanic = fn }
}

// NewSemaphoreGoroutinePool 创建按需派发的协程池。
//
// maxConcurrency 为同时执行的 Job 上限，必须大于 0。
//
// timeout 为额度耗尽时 Push 的最长等待时间；为 0 表示只受 ctx 约束地一直等。
// 该等待与 GoroutinePool 共用包内的秒级时间轮 poolTimeWheel，因此精度受限：
// xtime.Wheel.After 会把小于 1s 的 timeout 抬到 1s、大于 19s 的截到 19s。
// 需要亚秒级控制时传 0，改由调用方自己给 ctx 设 deadline。
func NewSemaphoreGoroutinePool(maxConcurrency int, timeout time.Duration, opts ...SemaphorePoolOption) *SemaphoreGoroutinePool {
	xpanic.WhenTrue(maxConcurrency <= 0, "xpool: maxConcurrency must be positive, got %d", maxConcurrency)
	p := &SemaphoreGoroutinePool{
		tokens:  make(chan struct{}, maxConcurrency),
		timeout: timeout,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Push 提交 Job。返回 nil 表示 Job 已被派发执行；返回 error 表示 Job 未执行。
//
// 额度耗尽时按以下顺序结束等待：ctx 取消、timeout 到时、拿到额度。
func (p *SemaphoreGoroutinePool) Push(ctx context.Context, job Job) error {
	if job == nil {
		return xerror.NewText("goroutine pool push nil job")
	}
	run, err := p.Reserve(ctx)
	if err != nil {
		return err
	}
	run(job)
	return nil
}

// Reserve 先占住一个执行额度，把「能不能执行」与「执行什么」分成两步。
//
// 返回 nil error 时额度已在手，此后派发一定会发生；返回的 run 必须且只需调用一次：
// 传入 Job 表示派发执行，传 nil 表示放弃本次预约、立即归还额度。重复调用是安全的空操作，
// 但一次都不调会让该额度永久泄漏，池容量静默缩水。
//
// 等待额度的结束条件与 Push 一致：ctx 取消、timeout 到时、拿到额度。
//
// # 什么时候需要它
//
// 当调用方要在派发前获取「不可回收」的资源时。用 Push 的话，资源已经取走而 Push 可能失败，
// 那份资源就漏在外面且无从补偿。典型场景是取一个全局单调递增的序号：下游按序号连续性做保序
// 处理，缺号会让后续数据被压住等超时，而序号取出来就还不回去了。
//
// 用 Reserve 可以把顺序写成「先占额度 → 取资源 → 派发」，取资源发生在派发已确定之后，
// 也仍处在调用方的顺序执行流里（而不是 Job 内部的并发上下文），资源的先后顺序因此可控。
func (p *SemaphoreGoroutinePool) Reserve(ctx context.Context) (run func(Job), err error) {
	if p.IsClosed() {
		return nil, xerror.NewText("goroutine pool closed")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// timerC 为 nil 时该 case 永久阻塞，等价于不设超时
	var timerC <-chan struct{}
	if p.timeout > 0 {
		timerC = poolTimeWheel.After(p.timeout)
	}
	select {
	case <-ctx.Done():
		return nil, xerror.NewText("goroutine pool job push context done, running: %d", p.Running())
	case <-timerC:
		return nil, xerror.NewText("goroutine pool job push blocked with %s, running: %d", p.timeout, p.Running())
	case p.tokens <- struct{}{}:
	}

	// 拿到额度后再判一次：Close 期间不再启动新工作，让 Close 能尽快收干
	if p.IsClosed() {
		<-p.tokens
		return nil, xerror.NewText("goroutine pool closed")
	}

	// once 保证额度只被消费一次：调用方重复调 run 时不会重复归还，
	// 否则一次误用就会凭空多出一个额度、把并发上限撑破。
	var once sync.Once
	return func(job Job) {
		once.Do(func() {
			if job == nil {
				<-p.tokens
				return
			}
			p.goWithToken(job)
		})
	}, nil
}

// goWithToken 用已经占住的额度执行 job，额度在 job 结束后归还。
func (p *SemaphoreGoroutinePool) goWithToken(job Job) {
	go func() {
		// 归还额度的 defer 先注册、后执行，因此即使 onPanic 自身再 panic，
		// 额度也会在栈展开时归还。
		defer func() { <-p.tokens }()
		defer func() {
			reason := recover()
			if reason == nil {
				return
			}
			if p.onPanic != nil {
				p.onPanic(reason)
				return
			}
			log.Error(fmt.Sprintf("xpool: job panic, reason: %v\nstack:\n%s", reason, rtdebug.Stack()))
		}()
		job()
	}()
}

// Running 返回当前占用的额度数，即正在执行的 Job 数（含刚拿到额度尚未开始的）。
// 可用于导出并发水位指标，判断 maxConcurrency 是否贴近实际峰值。
func (p *SemaphoreGoroutinePool) Running() int { return len(p.tokens) }

// Cap 返回并发上限。
func (p *SemaphoreGoroutinePool) Cap() int { return cap(p.tokens) }

// IsClosed 协程池是否已关闭。
func (p *SemaphoreGoroutinePool) IsClosed() bool { return p.closeFlag.Get() == 1 }

// Close 关闭协程池：此后 Push 一律返回错误。
//
// Close 不等待在飞 Job，立即返回——与 GoroutinePool.Close 一样不阻塞。
// 需要等待收干时紧接着调用 Wait，由调用方用 ctx 决定最长等多久。
// 之所以不把等待做进 Close：一个卡住不返回的 Job 会让 Close 永久阻塞，
// 而优雅退出路径上这种无逃生口的等待比"没等干净"危害更大。
func (p *SemaphoreGoroutinePool) Close() {
	p.closeFlag.CompareAndSwap(0, 1)
}

// Wait 等待在飞 Job 全部结束，或 ctx 结束。返回 nil 表示已收干。
//
// 判定方式是占满全部额度：只有先前派发的 Job 都归还了额度才可能占满，
// 因此无需 WaitGroup，也不存在 Add 与 Wait 并发的竞态。占满后即刻释放，
// 不改变池的可用状态。
//
// 应先 Close 再 Wait：Close 之后不再有新 Job 抢额度，Wait 才能确定收敛；
// 若在持续有 Push 的情况下直接 Wait，可能一直凑不满额度而等到 ctx 超时。
func (p *SemaphoreGoroutinePool) Wait(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	acquired := 0
	defer func() {
		for i := 0; i < acquired; i++ {
			<-p.tokens
		}
	}()
	for acquired < cap(p.tokens) {
		select {
		case p.tokens <- struct{}{}:
			acquired++
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
