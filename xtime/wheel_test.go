package xtime

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWheel(t *testing.T) {
	Convey("print timeWheel interval effect", t, func() {
		w := NewWheel(time.Second, 3)
		defer w.Stop()

		time.Sleep(500 * time.Millisecond)
		t1 := time.Now()

		go func() {
			select {
			case <-w.After(1 * time.Second):
				t.Logf("expected 1s, got %s", time.Since(t1))
			}
		}()

		time.Sleep(490 * time.Millisecond)
		t2 := time.Now()

		go func() {
			select {
			case <-w.After(0 * time.Second):
				t.Logf("expected 1s, got %s", time.Since(t2))
			}
		}()

		for {
			select {
			case <-w.After(3 * time.Second):
				t.Logf("expected 2s, got %s", time.Since(t2))
				return
			}
		}
	})
}
