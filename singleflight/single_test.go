package singleflight

import (
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
)

func TestSingleFlight(t *testing.T) {
	Convey("test single flight", t, func() {
		var do1 = func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return 1, nil
		}
		var do2 = func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return 2, nil
		}

		var key = "key"
		var g = New()
		var wg sync.WaitGroup
		wg.Add(2)
		var ret1, ret2 interface{}
		go func() {
			ret1, _ = g.Do(key, do1)
			wg.Done()
		}()
		go func() {
			ret2, _ = g.Do(key, do2)
			wg.Done()
		}()
		wg.Wait()
		So(ret1, ShouldEqual, ret2)
	})
}
