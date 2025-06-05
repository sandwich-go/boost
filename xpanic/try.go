package xpanic

import (
	"fmt"
	"runtime/debug"
)

// https://github.com/manucorporat/try/blob/master/try.go

const rethrow_panic = "_____rethrow"

type (
	E         interface{}
	exception struct {
		finally func()
		Error   E
		Stack   []byte
	}
)

func (e exception) String() string {
	return fmt.Sprintf("%v\n%s", e.Error, string(e.Stack))
}

func Throw() {
	panic(rethrow_panic)
}

func Try(f func()) (e exception) {
	e = exception{nil, nil, nil}
	// catch error in
	defer func() {
		if e.Error = recover(); e.Error != nil {
			e.Stack = debug.Stack()
		}
	}()
	f()
	return
}

func (e exception) Catch(f func(err E)) {
	if e.Error != nil {
		defer func() {
			// call finally
			if e.finally != nil {
				e.finally()
			}

			// rethrow exceptions
			if err := recover(); err != nil {
				if err == rethrow_panic {
					err = e.Error
				}
				panic(err)
			}
		}()
		f(e)
	} else if e.finally != nil {
		e.finally()
	}
}

func (e exception) Finally(f func()) (e2 exception) {
	if e.finally != nil {
		panic("finally was only set")
	}
	e2 = e
	e2.finally = f
	return
}
