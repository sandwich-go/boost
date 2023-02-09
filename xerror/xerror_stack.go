package xerror

import (
	"bytes"
	"fmt"
	"runtime"
)

func (cc *Error) Caller(skip int) (file, funcName string, line int) {
	if cc == nil {
		return
	}
	for _, p := range cc.callStack {
		if skip == 0 {
			if fn := runtime.FuncForPC(p - 1); fn == nil {
				funcName = "unknown"
			} else {
				funcName = fn.Name()
				file, line = fn.FileLine(p - 1)
			}
			return
		}
		skip--
	}
	return
}

// Stack returns the stack callers as string.
func (cc *Error) Stack() string {
	if cc == nil {
		return ""
	}
	var (
		curr   = cc
		index  = 1
		buffer = bytes.NewBuffer(nil)
	)
	for curr != nil {
		_, _ = buffer.WriteString(fmt.Sprintf("%d: %-v\n", index, curr))
		index++
		formatSubStack(curr.callStack, buffer)
		if curr.err == nil {
			break
		}
		if e, ok := curr.err.(*Error); ok {
			curr = e
		} else {
			_, _ = buffer.WriteString(fmt.Sprintf("%d. %s\n", index, curr.err.Error()))
			index++
			break
		}
	}
	return buffer.String()
}
