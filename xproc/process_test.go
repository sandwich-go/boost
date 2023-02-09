package xproc

import (
	"bytes"
	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xstrings"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func TestProcess(t *testing.T) {
	Convey("", t, func() {
		tmpFile := filepath.Join(os.TempDir(), "test.sh")
		err := xos.FilePutContents(tmpFile, []byte("echo \"GOT ME $1\""))
		So(err, ShouldBeNil)
		stdOut, errSteErr := ShellRun(tmpFile, WithArgs("1.2.0"))
		So(stdOut, ShouldEqual, "GOT ME 1.2.0")
		So(errSteErr, ShouldBeNil)

		stdOut1 := new(bytes.Buffer)
		stdErr1 := new(bytes.Buffer)
		p := NewProcessShellCmdWithOptions(tmpFile, NewProcessOptions(
			WithArgs("1.3.0"),
			WithStdout(stdOut1),
			WithStderr(stdErr1),
		))
		err = p.Run()
		So(err, ShouldBeNil)
		So(xstrings.Trim(stdOut1.String()), ShouldEqual, "GOT ME 1.3.0")
		So(stdErr1.Len(), ShouldBeZeroValue)
	})
}
