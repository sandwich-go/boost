package xpanic

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPanicWhen(t *testing.T) {
	Convey("panic when", t, func() {
		var err = errors.New("error")
		So(func() {
			WhenErrorAsFmtFirst(err, "%w, %d", 1)
		}, ShouldPanic)
		So(func() {
			WhenErrorAsFmtFirst(nil, "%w, %d", 1)
		}, ShouldNotPanic)

		Try(func() {
			WhenErrorAsFmtFirst(err, "%w, %d", 1)
		}).Catch(func(err E) {
			So(err.(error).Error(), ShouldEqual, "error, 1")
		})

		So(func() { WhenError(err) }, ShouldPanic)
		So(func() { WhenError(nil) }, ShouldNotPanic)
		So(func() { WhenTrue(true, "%d", 1) }, ShouldPanic)
		So(func() { WhenTrue(false, "%d", 1) }, ShouldNotPanic)
	})
}
