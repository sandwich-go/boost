// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package retry

import (
	"context"
	"time"
)

// Options should use NewOptions to initialize it
type Options struct {
	// annotation@Limit(comment="最大尝试次数")
	Limit uint
	// annotation@Delay(comment="固定延迟")
	Delay time.Duration
	// annotation@MaxJitter(comment="延迟最大抖动")
	MaxJitter time.Duration
	// annotation@OnRetry(comment="每次重试会先调用此方法")
	OnRetry func(n uint, err error) /*do nothing now*/
	// annotation@RetryIf(comment="何种error进行重试")
	RetryIf func(err error) bool
	// annotation@DelayType(comment="何种error进行重试")
	DelayType DelayTypeFunc
	// annotation@LastErrorOnly(comment="是否只返回最后遇到的error")
	LastErrorOnly bool
	// annotation@Context(comment="context，可以设定超时等")
	Context context.Context
	// annotation@MaxDelay(comment="最大延迟时间")
	MaxDelay         time.Duration
	MaxBackOffNInner uint
}

// NewOptions new Options
func NewOptions(opts ...Option) *Options {
	cc := newDefaultOptions()
	for _, opt := range opts {
		opt(cc)
	}
	if watchDogOptions != nil {
		watchDogOptions(cc)
	}
	return cc
}

// ApplyOption apply multiple new option
func (cc *Options) ApplyOption(opts ...Option) {
	for _, opt := range opts {
		opt(cc)
	}
}

// Option option func
type Option func(cc *Options)

// WithLimit 最大尝试次数
func WithLimit(v uint) Option {
	return func(cc *Options) {
		cc.Limit = v
	}
}

// WithDelay 固定延迟
func WithDelay(v time.Duration) Option {
	return func(cc *Options) {
		cc.Delay = v
	}
}

// WithMaxJitter 延迟最大抖动
func WithMaxJitter(v time.Duration) Option {
	return func(cc *Options) {
		cc.MaxJitter = v
	}
}

// WithOnRetry 每次重试会先调用此方法
func WithOnRetry(v func(n uint, err error)) Option {
	return func(cc *Options) {
		cc.OnRetry = v
	}
}

// WithRetryIf 何种error进行重试
func WithRetryIf(v func(err error) bool) Option {
	return func(cc *Options) {
		cc.RetryIf = v
	}
}

// WithDelayType 何种error进行重试
func WithDelayType(v DelayTypeFunc) Option {
	return func(cc *Options) {
		cc.DelayType = v
	}
}

// WithLastErrorOnly 是否只返回最后遇到的error
func WithLastErrorOnly(v bool) Option {
	return func(cc *Options) {
		cc.LastErrorOnly = v
	}
}

// WithContext context，可以设定超时等
func WithContext(v context.Context) Option {
	return func(cc *Options) {
		cc.Context = v
	}
}

// WithMaxDelay 最大延迟时间
func WithMaxDelay(v time.Duration) Option {
	return func(cc *Options) {
		cc.MaxDelay = v
	}
}

// InstallOptionsWatchDog the installed func will called when NewOptions  called
func InstallOptionsWatchDog(dog func(cc *Options)) { watchDogOptions = dog }

// watchDogOptions global watch dog
var watchDogOptions func(cc *Options)

// newDefaultOptions new default Options
func newDefaultOptions() *Options {
	cc := &Options{
		MaxBackOffNInner: 0,
	}

	for _, opt := range [...]Option{
		WithLimit(10),
		WithDelay(100 * time.Millisecond),
		WithMaxJitter(100 * time.Millisecond),
		WithOnRetry(func(n uint, err error) {
		}),
		WithRetryIf(IsRecoverable),
		WithDelayType(CombineDelay(BackOffDelay, RandomDelay)),
		WithLastErrorOnly(false),
		WithContext(context.Background()),
		WithMaxDelay(0),
	} {
		opt(cc)
	}

	return cc
}

// all getter func
func (cc *Options) GetLimit() uint                      { return cc.Limit }
func (cc *Options) GetDelay() time.Duration             { return cc.Delay }
func (cc *Options) GetMaxJitter() time.Duration         { return cc.MaxJitter }
func (cc *Options) GetOnRetry() func(n uint, err error) { return cc.OnRetry }
func (cc *Options) GetRetryIf() func(err error) bool    { return cc.RetryIf }
func (cc *Options) GetDelayType() DelayTypeFunc         { return cc.DelayType }
func (cc *Options) GetLastErrorOnly() bool              { return cc.LastErrorOnly }
func (cc *Options) GetContext() context.Context         { return cc.Context }
func (cc *Options) GetMaxDelay() time.Duration          { return cc.MaxDelay }
func (cc *Options) GetMaxBackOffNInner() uint           { return cc.MaxBackOffNInner }

// OptionsVisitor visitor interface for Options
type OptionsVisitor interface {
	GetLimit() uint
	GetDelay() time.Duration
	GetMaxJitter() time.Duration
	GetOnRetry() func(n uint, err error)
	GetRetryIf() func(err error) bool
	GetDelayType() DelayTypeFunc
	GetLastErrorOnly() bool
	GetContext() context.Context
	GetMaxDelay() time.Duration
	GetMaxBackOffNInner() uint
}

// OptionsInterface visitor + ApplyOption interface for Options
type OptionsInterface interface {
	OptionsVisitor
	ApplyOption(...Option)
}
