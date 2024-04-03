package xtime

import (
	"sync"
	"time"
)

type Timer struct {
	t      *time.Timer
	lock   *sync.RWMutex // 框架管理的timer不需要加锁，业务自己管理的timer需要加锁
	domain string
	cb     func()
}

func (t *Timer) Stop() {
	if t.lock != nil {
		t.lock.Lock()
		defer t.lock.Unlock()
	}
	t.t.Stop()
	t.cb = nil
}

func (t *Timer) Reset(d time.Duration) bool { return t.t.Reset(d) }

func (t *Timer) Cb() {
	if t.lock != nil {
		t.lock.Lock()
		defer t.lock.Unlock()
	} else {
		defer func() {
			t.cb = nil
		}()
	}
	if t.cb != nil {
		t.cb()
	}
}

type Cron struct{ t *Timer }

func (c *Cron) Stop() {
	if c.t != nil {
		c.t.Stop()
	}
}
