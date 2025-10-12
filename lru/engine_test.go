package lru

import (
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEngine(t *testing.T) {
	Convey("lru should work ok", t, func() {
		for _, v := range []struct {
			lazy bool
		}{
			{false},
			//{true},
		} {
			mx := sync.Mutex{}
			ttl := 10 * time.Millisecond
			var wg sync.WaitGroup
			wg.Add(1)
			handler := func(s string) {
				t.Log("expire,", s)
				wg.Done()
			}

			var e *Engine[string]
			if v.lazy {
				e = NewLazyEngine[string](ttl, &mx, handler)
			} else {
				e = NewEngine[string](ttl, &mx, handler)
			}

			mx.Lock()
			n := e.Add("1")
			So(e.size(), ShouldEqual, 1)
			e.Remove(n)
			e.Remove(n)
			So(e.size(), ShouldEqual, 0)
			n = e.Add("1")
			e.Promote(n)
			So(e.size(), ShouldEqual, 1)
			mx.Unlock()

			wg.Wait()
			mx.Lock()
			So(e.size(), ShouldEqual, 0)
			mx.Unlock()
		}
	})
}
