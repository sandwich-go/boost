package xpool

import (
	"errors"
	"github.com/sandwich-go/boost/xerror"
	"github.com/sandwich-go/boost/xsync"
	"github.com/sandwich-go/boost/xtime"
	"time"
)

type worker struct {
	jobChan chan Job

	// closeChan can be closed in order to cleanly shut down this worker.
	closeChan chan struct{}
	// closedChan is closed by the run() goroutine when it exits.
	closedChan chan struct{}
}

func (w *worker) Start(jobQueue chan Job) {
	defer func() {
		close(w.closedChan)
	}()

	go func() {
		var job Job
		for {
			select {
			case job = <-jobQueue:
				if job == nil {
					return
				}
				job()
			case <-w.closeChan:
				return
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
	queuedJobs xsync.AtomicInt64
	jobQueue   chan Job
	workers    []*worker
	closeFlag  xsync.AtomicInt32
	timeout    time.Duration
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
func (p *GoroutinePool) Push(job Job) error {
	if p.IsClosed() {
		return errors.New("pool closed")
	}
	if p.timeout == 0 {
		p.jobQueue <- job
	} else {
		select {
		case <-poolTimeWheel.After(p.timeout):
			return xerror.NewText("goroutine pool job queue blocked with %s", p.timeout)
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

	// Add extra workers if N > len(workers)
	for i := lWorkers; i < n; i++ {
		w := newWorker()
		w.Start(p.jobQueue)
		p.workers = append(p.workers, w)
	}

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
}

// IsClosed 协程竞争池是否已关闭
func (p *GoroutinePool) IsClosed() bool {
	return p.closeFlag.Get() == 1
}

// Close 关闭协程竞争池
func (p *GoroutinePool) Close() {
	if p.closeFlag.CompareAndSwap(0, 1) {
		p.SetSize(0)
		close(p.jobQueue)
	}
}
