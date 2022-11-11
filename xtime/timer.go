package xtime

import "time"

type Timer struct {
	t      *time.Timer
	domain string
	cb     func()
}

func (t *Timer) Stop() {
	t.t.Stop()
	t.cb = nil
}
func (t *Timer) Reset(d time.Duration) bool { return t.t.Reset(d) }

func (t *Timer) Cb() {
	defer func() {
		t.cb = nil
	}()
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
