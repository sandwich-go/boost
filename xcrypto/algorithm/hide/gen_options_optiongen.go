// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package hide

// Options should use NewOptions to initialize it
type Options struct {
	Suffix          string
	PrefixKeep      int
	SuffixKeep      int
	HideLenMin      int
	HideReplaceWith rune
	// 0 等长，否则按照指定的长度替换
	HideReplaceLen int
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

// WithSuffix option func for filed Suffix
func WithSuffix(v string) Option {
	return func(cc *Options) {
		cc.Suffix = v
	}
}

// WithPrefixKeep option func for filed PrefixKeep
func WithPrefixKeep(v int) Option {
	return func(cc *Options) {
		cc.PrefixKeep = v
	}
}

// WithSuffixKeep option func for filed SuffixKeep
func WithSuffixKeep(v int) Option {
	return func(cc *Options) {
		cc.SuffixKeep = v
	}
}

// WithHideLenMin option func for filed HideLenMin
func WithHideLenMin(v int) Option {
	return func(cc *Options) {
		cc.HideLenMin = v
	}
}

// WithHideReplaceWith option func for filed HideReplaceWith
func WithHideReplaceWith(v rune) Option {
	return func(cc *Options) {
		cc.HideReplaceWith = v
	}
}

// WithHideReplaceLen option func for filed HideReplaceLen
func WithHideReplaceLen(v int) Option {
	return func(cc *Options) {
		cc.HideReplaceLen = v
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
		WithSuffix("hash"),
		WithPrefixKeep(3),
		WithSuffixKeep(3),
		WithHideLenMin(3),
		WithHideReplaceWith('*'),
		WithHideReplaceLen(0),
	} {
		opt(cc)
	}

	return cc
}

// all getter func
func (cc *Options) GetSuffix() string        { return cc.Suffix }
func (cc *Options) GetPrefixKeep() int       { return cc.PrefixKeep }
func (cc *Options) GetSuffixKeep() int       { return cc.SuffixKeep }
func (cc *Options) GetHideLenMin() int       { return cc.HideLenMin }
func (cc *Options) GetHideReplaceWith() rune { return cc.HideReplaceWith }
func (cc *Options) GetHideReplaceLen() int   { return cc.HideReplaceLen }

// OptionsVisitor visitor interface for Options
type OptionsVisitor interface {
	GetSuffix() string
	GetPrefixKeep() int
	GetSuffixKeep() int
	GetHideLenMin() int
	GetHideReplaceWith() rune
	GetHideReplaceLen() int
}

// OptionsInterface visitor + ApplyOption interface for Options
type OptionsInterface interface {
	OptionsVisitor
	ApplyOption(...Option)
}
