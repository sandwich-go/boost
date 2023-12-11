package retry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/sandwich-go/boost/xerror"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDoFirstOk(t *testing.T) {
	Convey(`first call return ok`, t, func() {
		var retrySum uint
		err := Do(
			func(uint) error { return nil },
			WithOnRetry(func(n uint, err error) { retrySum += n }),
		)
		So(err, ShouldBeNil)
		So(retrySum, ShouldEqual, 0)
	})
	Convey(`do once when limit is zero`, t, func() {
		var doSum uint
		err := Do(
			func(uint) error {
				doSum += 1
				return nil
			},
			WithLimit(0),
		)

		So(err, ShouldBeNil)
		So(doSum, ShouldEqual, 1)
	})

	Convey(`retry if`, t, func() {
		var retryCount uint
		err := Do(
			func(uint) error {
				if retryCount >= 2 {
					return errors.New("special")
				} else {
					return errors.New("test")
				}
			},
			WithOnRetry(func(n uint, err error) { retryCount++ }),
			WithRetryIf(func(err error) bool {
				return err.Error() != "special"
			}),
			WithDelay(time.Nanosecond),
		)
		So(err, ShouldNotBeNil)
		var errWillBe xerror.Array
		errWillBe.Push(errors.New("test"))
		errWillBe.Push(errors.New("test"))
		errWillBe.Push(errors.New("special"))
		So(err.Error(), ShouldEqual, errWillBe.Error())
	})

	Convey(`default sleep`, t, func() {
		start := time.Now()
		err := Do(
			func(uint) error { return errors.New("test") },
			WithLimit(3),
		)
		So(err, ShouldNotBeNil)
		dur := time.Since(start)
		So(dur, ShouldBeGreaterThan, 3*newDefaultOptions().Delay)
	})
	Convey(`fixed sleep`, t, func() {
		start := time.Now()
		err := Do(
			func(uint) error { return errors.New("test") },
			WithLimit(3),
			WithDelayType(FixedDelay),
		)
		So(err, ShouldNotBeNil)
		dur := time.Since(start)
		So(dur, ShouldBeLessThan, 4*newDefaultOptions().Delay)
	})
	Convey(`last error only`, t, func() {
		var retrySum uint
		err := Do(
			func(uint) error { return fmt.Errorf("%d", retrySum) },
			WithOnRetry(func(n uint, err error) { retrySum += 1 }),
			WithDelay(time.Nanosecond),
			WithLastErrorOnly(true),
		)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "9")
	})
	Convey(`unrecoverable error`, t, func() {
		attempts := 0
		expectedErr := errors.New("error")
		err := Do(
			func(uint) error {
				attempts++
				return Unrecoverable(expectedErr)
			},
			WithLimit(2),
			WithLastErrorOnly(true),
		)
		So(err, ShouldNotBeNil)
		So(attempts, ShouldEqual, 1)
	})

	Convey(`max delay`, t, func() {
		start := time.Now()
		err := Do(
			func(uint) error { return errors.New("test") },
			WithLimit(5),
			WithDelay(10*time.Millisecond),
			WithMaxDelay(50*time.Millisecond),
		)
		dur := time.Since(start)
		So(err, ShouldNotBeNil)
		So(dur, ShouldBeLessThan, 205*time.Millisecond) // 重试5次,4个间隔*50ms
		So(dur, ShouldBeGreaterThan, 150*time.Millisecond)
	})

	Convey(`with context`, t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		retrySum := 0
		start := time.Now()
		err := Do(
			func(uint) error { return errors.New("test") },
			WithOnRetry(func(n uint, err error) { retrySum += 1 }),
			WithContext(ctx),
		)
		dur := time.Since(start)
		So(err, ShouldNotBeNil)
		So(dur, ShouldBeLessThan, newDefaultOptions().Delay)
	})

	Convey(`with context cancel in retry progress`, t, func() {
		ctx, cancel := context.WithCancel(context.Background())

		retrySum := 0
		err := Do(
			func(uint) error { return errors.New("test") },
			WithOnRetry(func(n uint, err error) {
				retrySum += 1
				if retrySum > 1 {
					cancel()
				}
			}),
			WithContext(ctx),
		)
		So(err, ShouldNotBeNil)
		So(retrySum, ShouldEqual, 2)
	})

	Convey(`just run`, t, func() {
		var (
			ClientPickerRetryLimit    uint = 10
			ClientPickerRetryMaxDelay      = time.Duration(500) * time.Millisecond
		)
		last := time.Now()
		_ = Do(
			func(attempt uint) (errPick error) {
				fmt.Println("attempt ", attempt, time.Since(last))
				last = time.Now()
				return fmt.Errorf("attempt %d", attempt)
			},
			WithLimit(ClientPickerRetryLimit),
			WithMaxDelay(ClientPickerRetryMaxDelay))
	})

	Convey(`backoff delay`, t, func() {
		for _, c := range []struct {
			label         string
			delay         time.Duration
			expectedMaxN  int
			n             uint
			expectedDelay time.Duration
		}{
			{
				label:         "negative-delay",
				delay:         -1,
				expectedMaxN:  62,
				n:             2,
				expectedDelay: 4,
			},
			{
				label:         "zero-delay",
				delay:         0,
				expectedMaxN:  62,
				n:             65,
				expectedDelay: 1 << 62,
			},
			{
				label:         "one-second",
				delay:         time.Second,
				expectedMaxN:  33,
				n:             62,
				expectedDelay: time.Second << 33,
			},
		} {
			cc := Options{
				Delay: c.delay,
			}
			delay := BackOffDelay(c.n, nil, &cc)
			So(c.expectedMaxN, ShouldEqual, cc.MaxBackOffNInner)
			So(c.expectedDelay, ShouldEqual, delay)
		}
	})
}

func TestRetryDelay(t *testing.T) {
	lastMilli := time.Now().UnixMilli()
	Do(func(attempt uint) error {
		tt := time.Now().UnixMilli()
		log.Println(tt, tt-lastMilli)
		lastMilli = tt
		return errors.New("some error")
	}, WithDelay(time.Millisecond*100), WithLimit(3))
}
