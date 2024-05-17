package xtime

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/xpanic"

	"github.com/sandwich-go/boost/xsync"
	"github.com/sandwich-go/boost/xtime/cron"
)

const (
	DefaultTimerDomain = "timer-domain-system"
)

// Lifecycle manages the start and close lifecycle methods.
type Lifecycle interface {
	Start()
	Close()
}

// TimerDispatcher manages timer-related functionalities.
type TimerDispatcher interface {
	AfterFunc(d time.Duration, cb func()) *SafeTimer
	AfterFuncWithOwnershipTransfer(td time.Duration, cb func()) *DanglingTimer
	CronFunc(cronExpr *cron.Expression, cb func()) *Cron
	// TickFunc 注册tick回调，key防止cb重复注册
	TickFunc(key interface{}, cb TickFunc)

	RemoveAllTimer()
	RemoveAllTimerInDomain(domain string)
	AfterFuncInDomain(td time.Duration, cb func(), domain string) *SafeTimer
	AfterFuncWithOwnershipTransferInDomain(td time.Duration, cb func(), domain string) *DanglingTimer

	TimerNotify() <-chan Timer
}

type Dispatcher interface {
	TimerDispatcher
	Lifecycle
}

type TickerDispatcher interface {
	TickerC() <-chan time.Time
	TriggerTickFuncs(context.Context)
}

type TickFunc func(context.Context)

// dispatcher one dispatcher per goroutine (goroutine not safe)
type dispatcher struct {
	ChanTimer     chan Timer // one receiver, N senders
	timerMutex    sync.RWMutex
	runningTimers sync.Map
	queueLen      int
	stopChan      chan struct{}
	closeFlag     xsync.AtomicInt32
	cc            *Options
	ticker        *time.Ticker
	tickMutex     sync.RWMutex
	tickFuncs     []*tickHandler
}

const (
	stateRunning = 0
	stateClosed  = 1
)

type tickHandler struct {
	key interface{}
	cb  TickFunc
}

// NewDispatcher 构造新的Dispatcher
// Note:
//   - Dispatcher目前主要服务于pkg/skeleton，用于Actor类型对象的有序动作处理
//   - Dispatcher通过TimerNotify将事件通知到外层，并不执行注册的的回调方法，需逻辑层接管事件通知触发回调逻辑
func NewDispatcher(l int, opts ...Option) Dispatcher {
	d := new(dispatcher)
	d.cc = NewOptions(opts...)
	d.queueLen = l
	d.closeFlag.Set(stateClosed)
	d.Start()
	return d
}

// Start for skeleton restart the dispatcher
func (d *dispatcher) Start() {
	if d.closeFlag.CompareAndSwap(stateClosed, stateRunning) {
		d.stopChan = make(chan struct{})
		d.ChanTimer = make(chan Timer, d.queueLen)
		if d.cc.TickDuration > 0 {
			d.ticker = time.NewTicker(d.cc.TickDuration)
			d.cc.TickCount.Inc()
			if d.cc.TickHostingMode {
				go d.tickHosting()
			}
		}
	}
}

func (d *dispatcher) tickHosting() {
	for {
		select {
		case <-d.stopChan:
			return
		case <-d.ticker.C:
			d.TriggerTickFuncs(context.Background())
		}
	}
}

func (d *dispatcher) TickFunc(key interface{}, cb TickFunc) {
	d.tickMutex.Lock()
	defer d.tickMutex.Unlock()
	for _, h := range d.tickFuncs {
		if h.key == key {
			log.Warn(fmt.Sprintf("multi tick func, key:%s", key))
			return
		}
	}
	d.tickFuncs = append(d.tickFuncs, &tickHandler{
		key: key,
		cb:  cb,
	})
}

// TickerC 返回ticker 的 chan，用于外部协程接管ticker执行协程
// d := NewDispatcher(
//
//		64,
//		WithTickDuration(time.Millisecond*50),
//		WithTickHostingMode(false),
//	)
//	d.Start()
//	td := d.(TickerDispatcher)
//	for {
//		select {
//		case <-td.TickerC():
//			td.TriggerTickFuncs(context.Background())
//		}
//	}
func (d *dispatcher) TickerC() <-chan time.Time {
	if d.ticker != nil {
		if d.cc.TickHostingMode {
			log.Warn("To indicate that the ticker is externally handled, you can set the TickHostingMode to false")
			return nil
		}
		return d.ticker.C
	}
	return nil
}

func (d *dispatcher) TriggerTickFuncs(ctx context.Context) {
	if d.closeFlag.Get() != stateRunning {
		return
	}
	d.tickMutex.RLock()
	defer d.tickMutex.RUnlock()
	for _, h := range d.tickFuncs {
		xpanic.Try(func() {
			h.cb(ctx)
		}).Catch(func(err xpanic.E) {
			log.Error(fmt.Sprintf("panic in tick funcs, reason:%v", err))
		})
	}
}

func (d *dispatcher) TimerNotify() <-chan Timer { return d.ChanTimer }
func (d *dispatcher) Close() {
	if d.closeFlag.CompareAndSwap(stateRunning, stateClosed) {
		close(d.stopChan)
		if d.ticker != nil {
			d.ticker.Stop()
			d.cc.TickCount.Dec()
		}
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

func (d *dispatcher) RemoveTimer(t Timer) { d.runningTimers.Delete(t) }

func (d *dispatcher) RemoveAllTimer() {
	d.runningTimers.Range(func(key, value interface{}) bool {
		t := key.(Timer)
		t.stop()
		d.RemoveTimer(t)
		return true
	})
}
func (d *dispatcher) RemoveAllTimerInDomain(domain string) {
	d.runningTimers.Range(func(key, value interface{}) bool {
		t := key.(Timer)
		if t.GetDomain() != domain {
			return true
		}
		t.stop()
		d.RemoveTimer(t)
		return true
	})
}

// AfterFuncWithOwnershipTransfer 自己管理 *Timer 声明周期, 使用完毕后需要手动调用 xtime Stop() 方法释放资源, 以免内存泄露
// 可以使用 Reset() 方法重置定时器
func (d *dispatcher) AfterFuncWithOwnershipTransfer(td time.Duration, cb func()) *DanglingTimer {
	return d.AfterFuncWithOwnershipTransferInDomain(td, cb, DefaultTimerDomain)
}

func (d *dispatcher) AfterFuncWithOwnershipTransferInDomain(td time.Duration, cb func(), domain string) *DanglingTimer {
	t := new(DanglingTimer)
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
			if d.closeFlag.Get() == stateRunning {
				d.ChanTimer <- t
			}
			d.timerMutex.RUnlock()
		}
	})
	log.Debug(fmt.Sprintf("Timer dispatcher add AfterFuncInDomain:%s after:%s", domain, td))
	return t
}

// AfterFunc 框架管理 Timer, dispatcher close时自动释放, 无需手动Stop
// Reset() 方法无效
func (d *dispatcher) AfterFunc(td time.Duration, cb func()) *SafeTimer {
	return d.AfterFuncInDomain(td, cb, DefaultTimerDomain)
}
func (d *dispatcher) AfterFuncInDomain(td time.Duration, cb func(), domain string) *SafeTimer {
	t := new(SafeTimer)
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
			if d.closeFlag.Get() == stateRunning {
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
