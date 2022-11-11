package xtime

import (
	"fmt"
	"github.com/sandwich-go/boost/internal/log"
	"sync"
	"time"

	"github.com/sandwich-go/boost/xsync"
	"github.com/sandwich-go/boost/xtime/cron"
)

const (
	DefaultTimerDomain = "timer-domain-system"
)

type Dispatcher interface {
	AfterFunc(d time.Duration, cb func()) *Timer
	CronFunc(cronExpr *cron.Expression, cb func()) *Cron

	RemoveAllTimer()
	RemoveAllTimerInDomain(domain string)
	AfterFuncInDomain(td time.Duration, cb func(), domain string) *Timer

	TimerNotify() <-chan *Timer
	Close()
	Start()
}

// dispatcher one dispatcher per goroutine (goroutine not safe)
type dispatcher struct {
	ChanTimer     chan *Timer // one receiver, N senders
	timerMutex    sync.RWMutex
	runningTimers sync.Map
	queueLen      int
	stopChan      chan struct{}
	closeFlag     xsync.AtomicInt32
}

// NewDispatcher 构造新的Dispatcher
// Note:
//      - Dispatcher目前主要服务于pkg/skeleton，用于Actor类型对象的有序动作处理
//      - Dispatcher通过TimerNotify将事件通知到外层，并不执行注册的的回调方法，需逻辑层接管事件通知触发回调逻辑
func NewDispatcher(l int) Dispatcher {
	d := new(dispatcher)
	d.queueLen = l
	d.closeFlag.Set(1)
	d.Start()
	return d
}

// for skeleton restart the dispatcher
func (d *dispatcher) Start() {
	if d.closeFlag.CompareAndSwap(1, 0) {
		d.stopChan = make(chan struct{})
		d.ChanTimer = make(chan *Timer, d.queueLen)
	}
}

func (d *dispatcher) TimerNotify() <-chan *Timer { return d.ChanTimer }
func (d *dispatcher) Close() {
	if d.closeFlag.CompareAndSwap(0, 1) {
		close(d.stopChan)
		// close #174
		d.timerMutex.Lock()
		close(d.ChanTimer)
		d.timerMutex.Unlock()

		//clear all timers
		for t := range d.ChanTimer {
			t.Cb()
		}
		d.RemoveAllTimer()
	}
}

func (d *dispatcher) RemoveTimer(t *Timer) { d.runningTimers.Delete(t) }

func (d *dispatcher) RemoveAllTimer() {
	d.runningTimers.Range(func(key, value interface{}) bool {
		t := key.(*Timer)
		t.Stop()
		d.RemoveTimer(t)
		return true
	})
}
func (d *dispatcher) RemoveAllTimerInDomain(domain string) {
	d.runningTimers.Range(func(key, value interface{}) bool {
		t := key.(*Timer)
		if t.domain != domain {
			return true
		}
		t.Stop()
		d.RemoveTimer(t)
		return true
	})
}

func (d *dispatcher) AfterFunc(td time.Duration, cb func()) *Timer {
	return d.AfterFuncInDomain(td, cb, DefaultTimerDomain)
}
func (d *dispatcher) AfterFuncInDomain(td time.Duration, cb func(), domain string) *Timer {
	t := new(Timer)
	t.cb = cb
	t.domain = domain
	t.t = time.AfterFunc(td, func() {
		// callback from another goroutine
		select {
		// FIRST read from no buffer chan, even closed, will return false
		case <-d.stopChan:
			return
		default:
			// close #174 (走到这里时，Close被执行了，这里的ChanTimer可能被close了)
			d.timerMutex.RLock()
			if d.closeFlag.Get() == 0 {
				d.ChanTimer <- t
				// 这里需要删除，否则runningTimers有泄漏
				d.RemoveTimer(t)
			}
			d.timerMutex.RUnlock()
		}
	})
	d.runningTimers.Store(t, struct{}{})
	log.Debug(fmt.Sprintf("Timer dispatcher add AfterFuncInDomain:%s after:%s", domain, td))
	return t
}

func (d *dispatcher) CronFunc(cronExpr *cron.Expression, callBack func()) *Cron {
	c := new(Cron)

	now := time.Now()
	nextTime := cronExpr.Next(now)
	if nextTime.IsZero() {
		return c
	}

	// callback
	var cb func()
	cb = func() {
		defer callBack()
		now := time.Now()
		nextTime := cronExpr.Next(now)
		if nextTime.IsZero() {
			return
		}
		c.t = d.AfterFunc(nextTime.Sub(now), cb)
	}

	c.t = d.AfterFunc(nextTime.Sub(now), cb)
	return c
}
