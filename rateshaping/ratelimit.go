package rateshaping // import "go.uber.org/ratelimit"

import (
	"time"

	"github.com/sandwich-go/boost/z"
)

// Note: This file is inspired by:
// https://github.com/prashantv/go-bench/blob/master/ratelimit

// Limiter is used to rate-limit some process, possibly across goroutines.
// The process is expected to call Take() before every iteration, which
// may block to throttle the goroutine.
type Limiter interface {
	// Take should block to make sure that the RPS is met.
	Take() time.Time
}

// New returns a Limiter that will limit to the given RPS.
func New(rate int, opts ...Option) Limiter {
	return newAtomicInt64Based(rate, opts...)
}

type unlimited struct{}

// NewUnlimited returns a RateLimiter that is not limited.
func NewUnlimited() Limiter {
	return unlimited{}
}

func (unlimited) Take() time.Time {
	return z.Now()
}
