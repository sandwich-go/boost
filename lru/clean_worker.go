package lru

import (
	"fmt"
	"time"

	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/module"
	"github.com/sandwich-go/boost/xmath"
	"github.com/sandwich-go/boost/xpool"
	"github.com/sandwich-go/boost/xtime"
)

// jitterPercent lru 清理任务的间隔抖动百分比；与历史 xmath.Disturb(d, 10)
// 行为对齐：每轮间隔在 [d*0.9, d*1.1] 之间。
const jitterPercent = 10

type CleanWorker interface {
	Clean(d time.Duration, id string, f func())
}

var DefaultCleanWorker CleanWorker = NewWorkerPerEngine()

func NewWorkerPerEngine() CleanWorker {
	return &workerPerEngine{}
}

type workerPerEngine struct {
}

func (w *workerPerEngine) Clean(d time.Duration, _ string, f func()) {
	// 走 xtime.PeriodicWithShutdown：单 task 整个生命周期只 alloc 一个底层
	// timer，每轮 reschedule 由 Periodic 内部复用同一个 *time.Timer
	// （或在 timewheel build 下复用同一个 tick 闭包）。相比手写
	// AfterFunc + 自 reschedule 的旧实现，避免稳态产生 per-tick 的 timer
	// / closure 分配。
	//
	// shutdown 信号关闭后下一轮自动跳过，等价于原实现 select 退出语义。
	xtime.PeriodicWithShutdown(d, jitterPercent, module.ShutdownNotify(), f)
}

func NewWorkerHashPool(numWorkers int, jobQueueLen int, timeout time.Duration) CleanWorker {
	return &workerHashPool{
		numWorkers: numWorkers,
		pool:       xpool.NewHashGoroutinePool(numWorkers, jobQueueLen, timeout),
	}
}

type workerHashPool struct {
	numWorkers int
	pool       *xpool.HashGoroutinePool
}

func (w *workerHashPool) Clean(d time.Duration, id string, f func()) {
	var cleanFunc func()
	cleanFunc = func() {
		f()

		select {
		case <-module.ShutdownNotify():
			return
		default:
		}

		xtime.AfterFunc(xmath.Disturb(d, 10), func() {
			if err := w.pool.PushJob(id, cleanFunc); err != nil {
				log.Error(fmt.Sprintf("lru worker push job error:%v", err))
			}
		})
	}
	xtime.AfterFunc(d, func() {
		if err := w.pool.PushJob(id, cleanFunc); err != nil {
			log.Error(fmt.Sprintf("lru worker push job error:%v", err))
		}
	})
}
