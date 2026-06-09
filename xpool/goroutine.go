package xpool

import (
	"context"
	"errors"
	"time"

	"github.com/sandwich-go/boost/xerror"
	"github.com/sandwich-go/boost/xsync"
	"github.com/sandwich-go/boost/xtime"
)

type worker struct {
	jobChan chan Job

	// closeChan can be closed in order to cleanly shut down this worker.
	closeChan chan struct{}
	// closedChan is closed by the run() goroutine when it exits.
	closedChan chan struct{}
}

func (w *worker) Start(jobQueue chan Job) {
	// 历史 bug 1：defer close(w.closedChan) 原本写在 Start 函数顶部，会
	// 在 Start 返回时（即 go func 起完即返）立刻 close closedChan，
	// 让 worker.join() 不再等待 worker goroutine 真退出。表现：
	// GoroutinePool.Close 立即返但 worker goroutine 仍存活（leak）。
	// 修复：把 defer close 移到 worker goroutine 内部，让 closedChan
	// 在 worker 真退出时才 close。
	//
	// 历史 bug 2（drain 不全）：原实现 closeChan 收到后立即 return，不
	// drain jobQueue 里残余的 job —— SetSize(0) / Close 时丢失 N 个未
	// 取的 job。下游 dataserver 的 Publish 路径 in-flight job 走
	// asyncQueue.Produce 把落库 task 入 stream，丢失等同丢未落库任务。
	// 修复：closeChan 收到后切 drain 模式（非阻塞 select 取 jobQueue 直
	// 到空），N 个 worker 同时 drain 互相竞争 jobQueue receive 互斥，
	// 最终 jobQueue 空 → 全 worker 退出。SetSize / Close caller 要保证
	// drain 期间不再 Push 新 job（GoroutinePool.Close 通过 closeFlag 让
	// Push 短路返 error；SetSize 缩容场景由 caller 自己保证）。
	//
	// nil job 检测保留作 defensive：理论上 jobQueue 不会被 send nil，但
	// 防误用（caller 显式 Push(nil) 或 channel close 后取到 zero value）
	// 让 worker 不 panic 而是优雅退出。
	go func() {
		defer close(w.closedChan)
		var job Job
		for {
			select {
			case job = <-jobQueue:
				if job == nil {
					return
				}
				job()
			case <-w.closeChan:
				// drain 残余 job：非阻塞 select，default 退出
				for {
					select {
					case job = <-jobQueue:
						if job == nil {
							return
						}
						job()
					default:
						return
					}
				}
			}
		}
	}()
}

func (w *worker) stop() { close(w.closeChan) }
func (w *worker) join() { <-w.closedChan }

func newWorker() *worker {
	return &worker{
		jobChan:    make(chan Job),
		closeChan:  make(chan struct{}),
		closedChan: make(chan struct{}),
	}
}

// Job 被 worker 竞争的工作
type Job func()

// GoroutinePool 线程池，numWorkers 数量的 worker 竞争 Job
type GoroutinePool struct {
	jobQueue  chan Job
	workers   []*worker
	closeFlag xsync.AtomicInt32
	timeout   time.Duration
}

// NewGoroutinePool 创建新的协程竞争池
// numWorkers 数量的 worker 竞争 Job
// jobQueueLen 设置 job 队列长度
// timeout 若 job 队列满，Push job 的超时时间
func NewGoroutinePool(numWorkers int, jobQueueLen int, timeout time.Duration) *GoroutinePool {
	pool := &GoroutinePool{jobQueue: make(chan Job, jobQueueLen), timeout: timeout}
	pool.SetSize(numWorkers)
	return pool
}

var poolTimeWheel = xtime.NewWheel(time.Second, 20)

// Push 放入 job 至job 队列
// 若设置了 timeout，当 job 队列满，Push 阻塞 timeout 会报错
func (p *GoroutinePool) Push(ctx context.Context, job Job) error {
	if p.IsClosed() {
		return errors.New("pool closed")
	}
	if p.timeout == 0 {
		select {
		case <-ctx.Done():
			return xerror.NewText("goroutine pool job push context done")
		case p.jobQueue <- job:
		}
	} else {
		select {
		case <-poolTimeWheel.After(p.timeout):
			return xerror.NewText("goroutine pool job queue blocked with %s", p.timeout)
		case <-ctx.Done():
			return xerror.NewText("goroutine pool job push context done")
		case p.jobQueue <- job:
		}
	}
	return nil
}

func (p *GoroutinePool) SetSize(n int) {
	lWorkers := len(p.workers)
	if lWorkers == n {
		return
	}

	if lWorkers > n {
		// Asynchronously stop all workers > N
		for i := n; i < lWorkers; i++ {
			p.workers[i].stop()
		}

		// Synchronously wait for all workers > N to stop
		for i := n; i < lWorkers; i++ {
			p.workers[i].join()
		}

		// Remove stopped workers from slice
		p.workers = p.workers[:n]

		return
	}

	// Add extra workers if N > len(workers)
	workers := make([]*worker, n)
	copy(workers, p.workers)
	for i := lWorkers; i < n; i++ {
		w := newWorker()
		w.Start(p.jobQueue)
		workers[i] = w
	}
	p.workers = workers
}

// IsClosed 协程竞争池是否已关闭
func (p *GoroutinePool) IsClosed() bool {
	return p.closeFlag.Get() == 1
}

// Close 关闭协程竞争池——同步等所有 worker 把 jobQueue 残余 job drain
// 跑完后退出。幂等。
//
// 关停语义（2026-06 改）：
//
//   - 先 closeFlag CAS 0→1 让 Push 短路返 "pool closed"（防 Close 期
//     间新 job 进来；下游 caller 还应自行保证关停期不再 Push）
//   - SetSize(0) 同步等所有 worker drain 完 jobQueue + 退出
//     （worker.Start 收到 closeChan 后切 drain 模式，详见其注释）
//   - **不**关 jobQueue：原实现 close(jobQueue) 与 race 进入 send 路径
//     的 Push 同时发生会触发 send-on-closed-channel panic；不关则 race
//     进入的 Push 仅 silent 丢 job（不 panic）—— 上层 caller 通过先停
//     接受新请求再 Close 的关停顺序兜住，下游 dataserver 走 module 框
//     架的 OnClose 顺序保证。
//
// 历史背景：原实现 SetSize(0) 内 worker 收 closeChan 立即退出**不
// drain** + close(jobQueue) 让 jobQueue 残余 job 永不被取——同时
// close(jobQueue) 还与 Push 撞 send-on-closed-channel panic race。
// 本次改 worker drain 模式同时撤销 close(jobQueue)，两个问题一起解决。
func (p *GoroutinePool) Close() {
	if p.closeFlag.CompareAndSwap(0, 1) {
		p.SetSize(0)
	}
}
