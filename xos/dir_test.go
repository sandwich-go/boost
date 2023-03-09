//go:build !windows
// +build !windows

package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func TestDir(t *testing.T) {
	Convey("IsEmpty", t, func() {
		temDir := os.TempDir()
		temDirTesting, err := os.MkdirTemp(temDir, "issue")
		So(err, ShouldBeNil)
		So(IsEmpty(temDirTesting), ShouldBeTrue)
		filename := ".config"
		err = os.WriteFile(filepath.Join(temDirTesting, filename), []byte(filename), 0666)
		So(err, ShouldBeNil)
		So(IsEmpty(temDirTesting), ShouldBeFalse)
	})
}
