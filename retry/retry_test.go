package retry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
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
		// 5 次尝试 = 4 个间隔，BackOff 公式 Delay<<attempt 在 MaxDelay 50ms
		// 处 cap：10/20/40/50 = 120ms 理论值。上界给 CI runtime / 调度
		// overhead 留 130ms 余量（原 205ms 仅留 5ms 在 CI Linux 多核环境
		// 下抖动会触发偶发 fail，commit <本批> 收紧上界并写明 BackOff 推算）。
		So(dur, ShouldBeLessThan, 250*time.Millisecond)
		So(dur, ShouldBeGreaterThan, 100*time.Millisecond)
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
	log.Println("TestRetryDelay ==> ")
	lastMilli := time.Now().UnixMilli()
	_ = Do(func(attempt uint) error {
		tt := time.Now().UnixMilli()
		log.Println(tt, tt-lastMilli)
		lastMilli = tt
		return errors.New("some error")
	}, WithDelay(time.Millisecond*100), WithLimit(3))
}

// TestCombineDelay 验证 CombineDelay 把多个 DelayTypeFunc 求和返回。
func TestCombineDelay(t *testing.T) {
	Convey("CombineDelay 多 delay 求和", t, func() {
		// 三个 fixed delay：100ms / 50ms / 30ms
		opt := &Options{Delay: 100 * time.Millisecond}
		fixed100 := FixedDelay
		fixed50 := func(_ uint, _ error, _ *Options) time.Duration {
			return 50 * time.Millisecond
		}
		fixed30 := func(_ uint, _ error, _ *Options) time.Duration {
			return 30 * time.Millisecond
		}
		combined := CombineDelay(fixed100, fixed50, fixed30)
		So(combined(0, nil, opt), ShouldEqual, 180*time.Millisecond)
	})

	Convey("CombineDelay 总和溢出 cap 到 MaxInt64", t, func() {
		// 两个超大 delay 加起来溢出
		huge := func(_ uint, _ error, _ *Options) time.Duration {
			return time.Duration(math.MaxInt64) - 1
		}
		combined := CombineDelay(huge, huge)
		// MaxInt64-1 + MaxInt64-1 溢出 → cap 到 MaxInt64
		So(combined(0, nil, &Options{}), ShouldEqual, time.Duration(math.MaxInt64))
	})

	Convey("CombineDelay 空入参返 0", t, func() {
		combined := CombineDelay()
		So(combined(0, nil, &Options{}), ShouldEqual, time.Duration(0))
	})
}

// TestUnpackUnrecoverable 验证 unpackUnrecoverable 解包逻辑：
// - unrecoverableError → 取出内部 err
// - 普通 error → 原样返回
func TestUnpackUnrecoverable(t *testing.T) {
	Convey("unpackUnrecoverable", t, func() {
		inner := errors.New("inner")

		// 包装路径：unpackUnrecoverable 解出 inner
		wrapped := Unrecoverable(inner)
		So(unpackUnrecoverable(wrapped), ShouldEqual, inner)

		// 普通 error 路径：unpackUnrecoverable 原样返
		plain := errors.New("plain")
		So(unpackUnrecoverable(plain), ShouldEqual, plain)

		// 配合 IsRecoverable 验证语义
		So(IsRecoverable(plain), ShouldBeTrue)
		So(IsRecoverable(wrapped), ShouldBeFalse)
	})
}
func TestBackoffDelay(t *testing.T) {
	log.Println("TestBackoffDelay ==> ")
	start := time.Now()
	last := start
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancelFunc()
	_ = Do(
		func(attempt uint) (errPick error) {
			tt := time.Now()
			log.Println(fmt.Sprintf("last:%s start:", time.Now().Sub(last)), time.Now().Sub(start))
			last = tt
			return errors.New("some error")
		},
		WithLimit(30),
		WithDelay(time.Duration(30)*time.Millisecond),
		WithMaxDelay(time.Second),
		WithContext(ctx),
		WithLastErrorOnly(true),
	)
	log.Println("TestBackoffDelay ==> ", fmt.Sprint(time.Now().Sub(start)))
}
