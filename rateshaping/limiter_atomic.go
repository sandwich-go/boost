package rateshaping

import (
	"sync/atomic"
	"time"
	_ "unsafe"

	"github.com/sandwich-go/boost/z"
)

type atomicInt64Limiter struct {
	//lint:ignore U1000 Padding is unused but it is crucial to maintain performance
	// of this rate limiter in case of collocation with other frequently accessed memory.
	prepadding [64]byte // cache line size = 64; created to avoid false sharing.
	state      int64    // unix nanoseconds of the next permissions issue.
	//lint:ignore U1000 like prepadding.
	postpadding [56]byte // cache line size - state size = 64 - 8; created to avoid false sharing.

	perRequest time.Duration
	maxSlack   time.Duration
}

// newAtomicBased returns a new atomic based limiter.
func newAtomicInt64Based(rate int, opts ...Option) *atomicInt64Limiter {
	// TODO consider moving config building to the implementation
	// independent code.
	config := NewOptions(opts...)
	perRequest := config.per / time.Duration(rate)
	l := &atomicInt64Limiter{
		perRequest: perRequest,
		maxSlack:   time.Duration(config.Slack) * perRequest,
	}
	atomic.StoreInt64(&l.state, 0)
	return l
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *atomicInt64Limiter) Take() time.Time {
	var (
		newTimeOfNextPermissionIssue int64
		now                          int64
	)
	for {
		now = int64(z.MonoOffset())
		timeOfNextPermissionIssue := atomic.LoadInt64(&t.state)

		switch {
		case timeOfNextPermissionIssue == 0 || (t.maxSlack == 0 && now-timeOfNextPermissionIssue > int64(t.perRequest)):
			// if this is our first call or t.maxSlack == 0 we need to shrink issue time to now
			newTimeOfNextPermissionIssue = now
		case t.maxSlack > 0 && now-timeOfNextPermissionIssue > int64(t.maxSlack)+int64(t.perRequest):
			// a lot of nanoseconds passed since the last Take call
			// we will limit max accumulated time to maxSlack
			newTimeOfNextPermissionIssue = now - int64(t.maxSlack)
		default:
			// calculate the time at which our permission was issued
			newTimeOfNextPermissionIssue = timeOfNextPermissionIssue + int64(t.perRequest)
		}

		if atomic.CompareAndSwapInt64(&t.state, timeOfNextPermissionIssue, newTimeOfNextPermissionIssue) {
			break
		}
	}

	sleepDuration := time.Duration(newTimeOfNextPermissionIssue - now)
	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
		return z.NowWithOffset(time.Duration(newTimeOfNextPermissionIssue))
	}
	// return now if we don't sleep as atomicLimiter does
	return z.NowWithOffset(time.Duration(now))
}
