package lru

import (
	"log"
	"time"

	"github.com/sandwich-go/boost/module"
	"github.com/sandwich-go/boost/xmath"
	"github.com/sandwich-go/boost/xpool"
	"github.com/sandwich-go/boost/xtime"
)

type ClearnWorker interface {
	Clean(d time.Duration, id string, f func())
}

var DefaultCleanWorker ClearnWorker = NewWorkerPerEngine()

func NewWorkerPerEngine() ClearnWorker {
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
				log.Println("shutdown notify")
				return
			}
		}
	}()
}

func NewWorkerHashPool(numWorkers int, jobQueueLen int, timeout time.Duration) ClearnWorker {
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
			w.pool.PushJob(id, cleanFunc)
		})
	}
	xtime.AfterFunc(d, func() {
		w.pool.PushJob(id, cleanFunc)
	})
}
