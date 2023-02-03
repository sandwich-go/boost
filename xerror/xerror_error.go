package xerror

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

var (
	ErrorCodeOk              int32 = 0
	ErrorCodeUnsetAsDefault  int32 = 1
	goRootForFilter                = runtime.GOROOT()
	CodeHandlerForNotAPICode       = func(err error) int32 {
		if err == nil {
			return ErrorCodeOk
		}
		return ErrorCodeUnsetAsDefault
	}
)

func init() {
	goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
}

var _ APICode = &Error{}
var _ apiStack = &Error{}
var _ apiCause = &Error{}

// Error 实现error接口，返回所有的错误信息，text: error具体数据
func (cc *Error) Error() string {
	if cc == nil {
		return ""
	}
	errStr := cc.text
	if cc.err != nil {
		if cc.text != "" {
			errStr += ": "
		}
		errStr += cc.err.Error()
	}
	return errStr
}

// Code 返回错误码，如果Error为nil，则根据CodeHandlerForNotAPICode逻辑返回
func (cc *Error) Code() int32 {
	if cc == nil {
		// 注意这里使用nil作为error接口的输入值,不要使用cc引起接口类型不为空但是值为空的问题
		return CodeHandlerForNotAPICode(nil)
	}
	return cc.code
}

// Cause 返回错误的起因
func (cc *Error) Cause() error {
	if cc == nil {
		return nil
	}
	loop := cc
	for loop != nil {
		if loop.err != nil {
			if e, ok := loop.err.(*Error); ok {
				loop = e
			} else if ac, ok0 := loop.err.(apiCause); ok0 {
				return ac.Cause()
			} else {
				return loop.err
			}
		} else {
			// 直接返回最初的error对象
			return loop
		}
	}
	return nil
}

// Format formats the frame according to the fmt.Formatter interface.
// %s,%v 全部错误信息
// %-s,%-v 当前错误信息
// %+v,%+s 错误字符串 + stack信息
func (cc *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('-'):
			if cc.text != "" {
				_, _ = io.WriteString(s, cc.text)
			} else {
				_, _ = io.WriteString(s, cc.Error())
			}
		case s.Flag('+'):
			_, _ = io.WriteString(s, cc.Error()+"\n"+cc.Stack())
		default:
			_, _ = io.WriteString(s, cc.Error())
		}
	}
}
