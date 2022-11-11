// wheel implements a time wheel algorithm that is suitable for large numbers of timers.
// It is based on golang channel broadcast mechanism.

package xtime

import (
	"sync"
	"time"
)

var (
	timerMap = NewWheelMap()
	accuracy = 20 // means max 1/20 deviation
)

func newWheelWithDuration(t time.Duration) *Wheel {
	return NewWheel(t/time.Duration(accuracy), accuracy+1)
}

// WheelAfter 根据Duration复用时间轮
// Note:
//      内部会根据Duration创建时间轮， 相同Duration可以共用，这样带来的副作用就是如果时间不固定则会创建特别的多的时间轮
func WheelAfter(t time.Duration) <-chan struct{} {
	w, _ := timerMap.GetOrSetFuncLock(t, newWheelWithDuration)
	return w.After(t)
}

// SetAccuracy sets the accuracy for the timewheel.
// low accuracy usually have better performance.
func SetAccuracy(a int) {
	accuracy = a
}

type Wheel struct {
	sync.Mutex
	interval   time.Duration
	ticker     *time.Ticker
	quit       chan struct{}
	maxTimeout time.Duration
	cs         []chan struct{}
	pos        int
}

func NewWheel(interval time.Duration, buckets int) *Wheel {
	w := &Wheel{
		interval:   interval,
		quit:       make(chan struct{}),
		pos:        0,
		maxTimeout: interval * (time.Duration(buckets - 1)),
		cs:         make([]chan struct{}, buckets),
		ticker:     time.NewTicker(interval),
	}
	for i := range w.cs {
		w.cs[i] = make(chan struct{})
	}
	go w.run()
	return w
}

func (w *Wheel) Stop() {
	close(w.quit)
}

// After 误差在一个interval内
// timeline : ---w.pos-1<--{x}-->call After()<--{y}-->w.pos-----
// x + y == interval, y 即是误差
// Note:
//      如果超过时间轮的最大值则使用最大值作为Timeout时间
func (w *Wheel) After(timeout time.Duration) <-chan struct{} {
	if timeout > w.maxTimeout {
		timeout = w.maxTimeout
	} else if timeout < time.Second {
		timeout = time.Second
	}

	w.Lock()
	index := (w.pos + int(timeout/w.interval)) % len(w.cs)
	b := w.cs[index]
	w.Unlock()
	return b
}

func (w *Wheel) run() {
	for {
		select {
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

func (w *Wheel) onTicker() {
	w.Lock()
	lastC := w.cs[w.pos]
	w.cs[w.pos] = make(chan struct{})
	w.pos = (w.pos + 1) % len(w.cs)
	w.Unlock()
	close(lastC)
}
