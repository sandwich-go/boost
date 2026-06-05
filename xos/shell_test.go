package xos

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShell(t *testing.T) {
	Convey("shell", t, func() {
		So(GetShell(), ShouldNotBeEmpty)
		So([]string{"/c", "-c"}, ShouldContain, GetShellOption())
	})
}
