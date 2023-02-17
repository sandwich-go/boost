package xpanic

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCatch(t *testing.T) {
	Convey(`Catch will suppress a panic.`, t, func() {
		var pv *Panic
		So(func() {
			defer Catch(func(p *Panic) {
				pv = p
			})
			panic("Everybody panic!")
		}, ShouldNotPanic)

		So(pv, ShouldNotBeNil)
		So(pv.Reason, ShouldEqual, "Everybody panic!")
	})
}

// Example is a very simple example of how to use Catch to recover from a panic
// and log its stack trace.
func Example() {
	Do(func() {
		fmt.Println("Doing something...")
		panic("Something wrong happened!")
	}, func(p *Panic) {
		fmt.Println("Caught a panic:", p.Reason)
	})
	// Output: Doing something...
	// Caught a panic: Something wrong happened!
}
