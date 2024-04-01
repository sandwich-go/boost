package xtime

import (
	"context"
	"testing"
	"time"
)

func TestTick(t *testing.T) {
	d := NewDispatcher(0, WithTickDuration(time.Millisecond*5))
	count := 0
	stop := make(chan struct{})
	d.TickFunc("123", func(_ context.Context) {
		count += 1
		if count >= 2 {
			close(stop)
		}
	})
	d.Start()
	select {
	case <-stop:
		t.Log("trigger tick twice")
	case <-time.After(time.Millisecond * 12):
		t.Fatal("ticker not work")
	}
}

func TestTickExternalHost(t *testing.T) {
	d := NewDispatcher(0, WithTickDuration(time.Millisecond*5), WithTickHostingMode(false))
	count := 0
	stop := make(chan struct{})
	d.TickFunc("123", func(_ context.Context) {
		count += 1
		if count >= 2 {
			close(stop)
		}
	})
	d.Start()
	td := d.(TickerDispatcher)

	for {
		select {
		case <-td.TickerC():
			td.TriggerTickFuncs(context.Background())
		case <-stop:
			t.Log("trigger tick twice")
			return
		case <-time.After(time.Millisecond * 12):
			t.Fatal("ticker not work")
		}
	}

}
