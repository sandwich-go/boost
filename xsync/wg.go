package xsync

import (
	"context"
	"fmt"
	"github.com/sandwich-go/boost/internal/log"
	"sync"
	"time"
)

// WaitGroupTimeout is a wrapper of sync.WaitGroup that
// supports wait with timeout.
type WaitGroupTimeout struct {
	sync.WaitGroup
}

func (w *WaitGroupTimeout) WaitTimeout(timeout time.Duration) bool {
	return WaitTimeout(&w.WaitGroup, timeout)
}

// WaitContext performs a timed wait on a given sync.WaitGroup
func WaitContext(waitGroup *sync.WaitGroup, ctx context.Context) bool {
	success := make(chan struct{})
	go func() {
		defer func() {
			if reason := recover(); reason != nil {
				log.Error(fmt.Sprintf("wait context panic, reason: %v", reason))
			}
		}()
		defer close(success)
		waitGroup.Wait()
	}()
	select {
	case <-success: // completed normally
		return false
	case <-ctx.Done():
		return true
	}
}

// WaitTimeout performs a timed wait on a given sync.WaitGroup
func WaitTimeout(waitGroup *sync.WaitGroup, timeout time.Duration) bool {
	ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancelFunc()
	return WaitContext(waitGroup, ctx)
}
