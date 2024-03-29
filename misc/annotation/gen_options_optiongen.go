// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package annotation

import (
	"sync/atomic"
	"unsafe"
)

// Options should use NewOptions to initialize it
type Options struct {
	MagicPrefix string       `xconf:"magic_prefix" usage:"只有包含 MagicPrefix 的行，才能萃取到注释"`
	LowerKey    bool         `xconf:"lower_key" usage:"key是否为转化为小写"`
	Descriptors []Descriptor `xconf:"descriptors" usage:"描述数组"`
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

// ApplyOption apply multiple new option and return the old ones
// sample:
// old := cc.ApplyOption(WithTimeout(time.Second))
// defer cc.ApplyOption(old...)
func (cc *Options) ApplyOption(opts ...Option) []Option {
	var previous []Option
	for _, opt := range opts {
		previous = append(previous, opt(cc))
	}
	return previous
}

// Option option func
type Option func(cc *Options) Option

// WithMagicPrefix 只有包含 MagicPrefix 的行，才能萃取到注释
func WithMagicPrefix(v string) Option {
	return func(cc *Options) Option {
		previous := cc.MagicPrefix
		cc.MagicPrefix = v
		return WithMagicPrefix(previous)
	}
}

// WithLowerKey key是否为转化为小写
func WithLowerKey(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.LowerKey
		cc.LowerKey = v
		return WithLowerKey(previous)
	}
}

// WithDescriptors 描述数组
func WithDescriptors(v ...Descriptor) Option {
	return func(cc *Options) Option {
		previous := cc.Descriptors
		cc.Descriptors = v
		return WithDescriptors(previous...)
	}
}

// InstallOptionsWatchDog the installed func will called when NewOptions  called
func InstallOptionsWatchDog(dog func(cc *Options)) { watchDogOptions = dog }

// watchDogOptions global watch dog
var watchDogOptions func(cc *Options)

// newDefaultOptions new default Options
func newDefaultOptions() *Options {
	cc := &Options{}

	for _, opt := range [...]Option{
		WithMagicPrefix("annotation@"),
		WithLowerKey(true),
		WithDescriptors(nil...),
	} {
		opt(cc)
	}

	return cc
}

// AtomicSetFunc used for XConf
func (cc *Options) AtomicSetFunc() func(interface{}) { return AtomicOptionsSet }

// atomicOptions global *Options holder
var atomicOptions unsafe.Pointer

// onAtomicOptionsSet global call back when  AtomicOptionsSet called by XConf.
// use OptionsInterface.ApplyOption to modify the updated cc
// if passed in cc not valid, then return false, cc will not set to atomicOptions
var onAtomicOptionsSet func(cc OptionsInterface) bool

// InstallCallbackOnAtomicOptionsSet install callback
func InstallCallbackOnAtomicOptionsSet(callback func(cc OptionsInterface) bool) {
	onAtomicOptionsSet = callback
}

// AtomicOptionsSet atomic setter for *Options
func AtomicOptionsSet(update interface{}) {
	cc := update.(*Options)
	if onAtomicOptionsSet != nil && !onAtomicOptionsSet(cc) {
		return
	}
	atomic.StorePointer(&atomicOptions, (unsafe.Pointer)(cc))
}

// AtomicOptions return atomic *OptionsVisitor
func AtomicOptions() OptionsVisitor {
	current := (*Options)(atomic.LoadPointer(&atomicOptions))
	if current == nil {
		defaultOne := newDefaultOptions()
		if watchDogOptions != nil {
			watchDogOptions(defaultOne)
		}
		atomic.CompareAndSwapPointer(&atomicOptions, nil, (unsafe.Pointer)(defaultOne))
		return (*Options)(atomic.LoadPointer(&atomicOptions))
	}
	return current
}

// all getter func
func (cc *Options) GetMagicPrefix() string       { return cc.MagicPrefix }
func (cc *Options) GetLowerKey() bool            { return cc.LowerKey }
func (cc *Options) GetDescriptors() []Descriptor { return cc.Descriptors }

// OptionsVisitor visitor interface for Options
type OptionsVisitor interface {
	GetMagicPrefix() string
	GetLowerKey() bool
	GetDescriptors() []Descriptor
}

// OptionsInterface visitor + ApplyOption interface for Options
type OptionsInterface interface {
	OptionsVisitor
	ApplyOption(...Option) []Option
}
