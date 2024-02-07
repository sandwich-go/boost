package xpool

import (
	"testing"
	"time"
)

type Event struct {
	Reference

	Name   string
	Action string
	At     time.Time
}

func TestName(t *testing.T) {
	pool := NewReferencePool(func(r Reference) Reference {
		e := &Event{Reference: r}
		return e
	}, func(i interface{}) {
		t.Log("Reference Reset...")
		e := i.(*Event)
		e.Name = ""
		e.Action = ""
		e.At = time.Time{}
	})

	e := pool.Get().(*Event)
	e.Decr()
	e = pool.Get().(*Event)
	e.Decr()
	t.Log(pool.Stats())
}
