package xtime

import (
	"sync"
	"time"
)

type Timer interface {
	Cb()
	stop()
	GetDomain() string
}

type SafeTimer struct {
	t      *time.Timer
	domain string
	cb     func()
}

func (t *SafeTimer) stop() {
	t.t.Stop()
	t.cb = nil
}

func (t *SafeTimer) Cb() {
	defer func() {
		t.cb = nil
	}()
	if t.cb != nil {
		t.cb()
	}
}

func (t *SafeTimer) GetDomain() string {
	return t.domain
}

type DanglingTimer struct {
	t      *time.Timer
	lock   sync.RWMutex // 框架管理的timer不需要加锁，业务自己管理的timer需要加锁
	domain string
	cb     func()
}

func (t *DanglingTimer) stop() {
	t.Stop()
}

func (t *DanglingTimer) Stop() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.t.Stop()
	t.cb = nil
}

func (t *DanglingTimer) Reset(d time.Duration) bool { return t.t.Reset(d) }

func (t *DanglingTimer) Cb() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.cb != nil {
		t.cb()
	}
}

func (t *DanglingTimer) GetDomain() string {
	return t.domain
}

type Cron struct{ t Timer }

func (c *Cron) Stop() {
	if c.t != nil {
		c.t.stop()
	}
}
