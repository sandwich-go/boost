package xerror

import "os"

var IsErrorWithStack = false

// Error struct
type Error struct {
	err        error  // 底层错误
	text       string // 错误信息
	code       int32  // 错误码
	logic      bool   // 是否为逻辑层异常，主要为sandwich 的queue模式服务
	callStack  stack  // 堆栈信息，内部数据自动生成
	skip       int    // error skip
	setTimeout bool   // 是否设置超时, 如果明确设置超时状态，则不再根据底层错误判断
	timeout    bool   // 是否为超时错误
}

// Logic 是否为逻辑层错误
func (cc *Error) Logic() bool { return cc.logic }

// Timeout 是否为超时错误 os.IsTimeout
func (cc *Error) Timeout() bool {
	if !cc.setTimeout && cc.err != nil {
		return cc.timeout || os.IsTimeout(cc.err)
	}
	return cc.timeout
}

// SetTimeout 设定为超时
func (cc *Error) SetTimeout() *Error {
	cc.setTimeout = true
	cc.timeout = true
	return cc
}

// UnsetTimeout 设定为非超时
func (cc *Error) UnsetTimeout() *Error {
	cc.setTimeout = true
	cc.timeout = false
	return cc
}

// WithStack 设置堆栈
func (cc *Error) WithStack() *Error {
	cc.callStack = callers()
	return cc
}

// SetLogic 设定为逻辑层错误
func (cc *Error) SetLogic() *Error {
	cc.logic = true
	return cc
}

// UnsetLogic 设定为非逻辑层错误
func (cc *Error) UnsetLogic() *Error {
	cc.logic = false
	return cc
}

// Unwrap 兼容 errors.Unwrap
func (cc *Error) Unwrap() error { return cc.err }

// New 新建 Error 对象
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
func WithSkip(v int) ErrorOption    { return func(cc *Error) { cc.skip = v } }
func WithTimeout(v bool) ErrorOption {
	return func(cc *Error) {
		cc.setTimeout = true
		cc.timeout = v
	}
}

// WithStack option func for stack
func WithStack() ErrorOption {
	return func(cc *Error) {
		if cc.callStack != nil {
			return
		}
		cc.callStack = callers(1)
	}
}
