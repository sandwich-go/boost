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
	ch := make(chan struct{}, 1)
	go func() {
		xtime.AfterFunc(d, func() {
			ch <- struct{}{}
		})
		for {
			select {
			case <-ch:
				f()
				xtime.AfterFunc(xmath.Disturb(d, 10), func() {
					ch <- struct{}{}
				})
			case <-module.ShutdownNotify():
				log.Info("shutdown notify")
				return
			}
		}
	}()
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
