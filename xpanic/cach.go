package xpanic

import (
	"runtime"
)

// Panic is a snapshot of a panic, containing both the panic's reason and the
// system stack.
type Panic struct {
	// Reason is the value supplied to the recover function.
	Reason interface{}
	Stack  []byte
}

// Catch recovers from panic. It should be used as a deferred call.
//
// If the supplied panic callback is nil, the panic will be silently discarded.
// Otherwise, the callback will be invoked with the panic's information.
func Catch(cb func(p *Panic)) {
	if reason := recover(); reason != nil && cb != nil {
		cb(&Panic{
			Reason: reason,
			Stack: func() []byte {
				const size = 4096
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]

				return buf
			}(),
		})
	}
}

// Do executes f. If a panic occurs during execution, the supplied callback will
// be called with the panic's information.
//
// If the panic callback is nil, the panic will be caught and discarded silently.
func Do(f func(), cb func(p *Panic)) {
	defer Catch(cb)
	f()
}
