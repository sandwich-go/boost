package xerror

var IsErrorWithStack = false

// Error struct
type Error struct {
	err       error  // 底层错误
	text      string // 错误信息
	code      int32  // 错误码
	logic     bool   // 是否为逻辑层异常，主要为sandwich 的queue模式服务
	callStack stack  // 堆栈信息，内部数据自动生成
	skip      int    // error skip
}

// Logic 是否为逻辑层错误
func (cc *Error) Logic() bool { return cc.logic }

// SetLogic 设定为逻辑层错误
func (cc *Error) WithStack() *Error {
	cc.callStack = callers()
	return cc
}

// SetLogic 设定为逻辑层错误
func (cc *Error) SetLogic() *Error {
	cc.logic = true
	return cc
}

// UnsetLogic 设定为非逻辑层凑无
func (cc *Error) UnsetLogic() *Error {
	cc.logic = false
	return cc
}

func New(opts ...ErrorOption) *Error {
	e := &Error{callStack: nil}
	for _, opt := range opts {
		opt(e)
	}
	if IsErrorWithStack && e.callStack == nil {
		e.callStack = callers(e.skip)
	}
	return e
}

type ErrorOption func(cc *Error)

func WithErr(v error) ErrorOption   { return func(cc *Error) { cc.err = v } }
func WithText(v string) ErrorOption { return func(cc *Error) { cc.text = v } }
func WithCode(v int32) ErrorOption  { return func(cc *Error) { cc.code = v } }
func WithLogic(v bool) ErrorOption  { return func(cc *Error) { cc.logic = v } }

// WithText option func for text
func WithStack() ErrorOption {
	return func(cc *Error) {
		if cc.callStack != nil {
			return
		}
		cc.callStack = callers(1)
	}
}
