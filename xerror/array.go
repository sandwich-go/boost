package xerror

import (
	"fmt"
	"strings"
)

type Array struct {
	errors     []error
	formatFunc ErrorFormatFunc
}

// Push 根据SetFormatFunc方法构建错误信息，默认为ListFormatFunc
func (e Array) Error() string {
	fn := e.formatFunc
	if fn == nil {
		fn = ListFormatFunc
	}
	return fn(e.errors)
}

// Push 推入一个错误信息，err如果为nil则丢弃
func (e *Array) Push(err error) {
	if err == nil {
		return
	}
	e.errors = append(e.errors, err)
}

// LastErr 返回最后一个错误信息，如果没有错误则返回nil
func (e *Array) LastErr() error {
	if e == nil {
		return nil
	}
	if len(e.errors) == 0 {
		return nil
	}
	return e.errors[len(e.errors)-1]
}

// Err 返回标准error对象，如果错误列表为空则返回nil
func (e *Array) Err() error {
	if e == nil {
		return nil
	}
	if len(e.errors) == 0 {
		return nil
	}
	return e
}

func (e *Array) String() string                  { return fmt.Sprintf("*%#v", *e) }
func (e *Array) WrappedErrors() []error          { return e.errors }
func (e *Array) SetFormatFunc(f ErrorFormatFunc) { e.formatFunc = f }

type ErrorFormatFunc func([]error) string

func DotFormatFunc(es []error) string {
	var errStr = make([]string, 0)
	for i := 0; i < len(es); i++ {
		errStr = append(errStr, es[i].Error())
	}
	return strings.Join(errStr, ",")
}

func ListFormatFunc(es []error) string {
	points := make([]string, len(es))
	for i, err := range es {
		points[i] = fmt.Sprintf("#%d: %s", i+1, err)
	}
	return fmt.Sprintf(
		"%d errors occurred:\n%s",
		len(es), strings.Join(points, "\n"))
}
