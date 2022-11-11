package paniccatcher

// Panic is a snapshot of a panic, containing both the panic's reason and the
// system stack.
type Panic struct {
	// Reason is the value supplied to the recover function.
	Reason interface{}
}

// Catch recovers from panic. It should be used as a deferred call.
//
// If the supplied panic callback is nil, the panic will be silently discarded.
// Otherwise, the callback will be invoked with the panic's information.
func Catch(cb func(p *Panic)) {
	if reason := recover(); reason != nil && cb != nil {
		cb(&Panic{
			Reason: reason,
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
