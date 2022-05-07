package xproc

import (
	"testing"

	"github.com/sandwich-go/boost/xos"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShellRun(t *testing.T) {
	Convey("shell run", t, func() {
		tmpFile := "/tmp/test.sh"
		xos.FilePutContents(tmpFile, []byte("echo \"GOT ME $1\""))
		stdOut, errSteErr := ShellRun(tmpFile, WithArgs("1.2.0"))
		So(stdOut, ShouldEqual, "GOT ME 1.2.0")
		So(errSteErr, ShouldBeNil)
	})
}
