package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestShell(t *testing.T) {
	Convey("shell", t, func() {
		So(GetShell(), ShouldNotBeEmpty)
		So([]string{"/c", "-c"}, ShouldContain, GetShellOption())
	})
}
