package retry

import (
	"math"
	"math/rand"
	"time"

	"github.com/sandwich-go/boost/xerror"
)

type RetryableFunc func(attempt uint) error

func Do(retryableFunc RetryableFunc, opts ...Option) error {
	cc := NewOptions(opts...)
	if err := cc.Context.Err(); err != nil {
		return err
	}
	var limit = cc.Limit
	if limit <= 0 {
		limit = 1
	}
	var lastErr error
	var ea xerror.Array
	for attempt := uint(0); attempt < limit; attempt++ {
		err := retryableFunc(attempt)
		if err == nil {
			return nil
		}
		if cc.LastErrorOnly {
			lastErr = err
		} else {
			ea.Push(unpackUnrecoverable(err))
		}

		if !cc.RetryIf(err) {
			break
		}

		cc.OnRetry(attempt, err)

		// if this is last attempt - don't wait
		if attempt == cc.Limit-1 {
			break
		}

		delayTime := cc.DelayType(attempt, err, cc)
		if cc.MaxDelay > 0 && delayTime > cc.MaxDelay {
			delayTime = cc.MaxDelay
		}

		select {
		case <-time.After(delayTime):
		case <-cc.Context.Done():
			return cc.Context.Err()
		}
	}
	if cc.LastErrorOnly {
		return lastErr
	}
	return ea.Err()
}

type unrecoverableError struct {
	error
}

// Unrecoverable wraps an error in `unrecoverableError` struct
func Unrecoverable(err error) error {
	return unrecoverableError{err}
}

// IsRecoverable checks if error is an instance of `unrecoverableError`
func IsRecoverable(err error) bool {
	_, isUnrecoverable := err.(unrecoverableError)
	return !isUnrecoverable
}

func unpackUnrecoverable(err error) error {
	if unrecoverable, isUnrecoverable := err.(unrecoverableError); isUnrecoverable {
		return unrecoverable.error
	}

	return err
}

type IfFunc func(error) bool
type OnRetryFunc func(n uint, err error)
type DelayTypeFunc func(n uint, err error, opts *Options) time.Duration

// BackOffDelay 指数级增长重试的推迟时长
func BackOffDelay(n uint, _ error, opts *Options) time.Duration {
	// 1 << 63 would overflow signed int64 (time.Duration), thus 62.
	const max uint = 62

	if opts.MaxBackOffNInner == 0 {
		if opts.Delay <= 0 {
			opts.Delay = 1
		}
		opts.MaxBackOffNInner = max - uint(math.Floor(math.Log2(float64(opts.Delay))))
	}

	if n > opts.MaxBackOffNInner {
		n = opts.MaxBackOffNInner
	}

	return opts.Delay << n
}

// FixedDelay 固定延时
func FixedDelay(_ uint, _ error, config *Options) time.Duration { return config.Delay }

// RandomDelay 随机延时
func RandomDelay(_ uint, _ error, config *Options) time.Duration {
	return time.Duration(rand.Int63n(int64(config.MaxJitter)))
}

// CombineDelay is a DelayType the combines all of the specified delays into a new DelayTypeFunc
func CombineDelay(delays ...DelayTypeFunc) DelayTypeFunc {
	const maxInt64 = uint64(math.MaxInt64)

	return func(n uint, err error, config *Options) time.Duration {
		var total uint64
		for _, delay := range delays {
			total += uint64(delay(n, err, config))
			if total > maxInt64 {
				total = maxInt64
			}
		}
		return time.Duration(total)
	}
}
