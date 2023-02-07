package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	Convey("exists", t, func() {
		var file, dir = "a.go", "b"
		So(Exists(file), ShouldBeFalse)
		So(ExistsFile(file), ShouldBeFalse)

		So(Exists(dir), ShouldBeFalse)
		So(ExistsDir(dir), ShouldBeFalse)

		var dirFullPath = filepath.Join(os.TempDir(), dir)
		var fileFullPath = filepath.Join(dirFullPath, file)
		So(MkdirAll(fileFullPath), ShouldBeNil)
		So(Copy("exists.go", fileFullPath), ShouldBeNil)

		So(Exists(fileFullPath), ShouldBeTrue)
		So(ExistsFile(fileFullPath), ShouldBeTrue)

		So(Exists(dirFullPath), ShouldBeTrue)
		So(ExistsDir(dirFullPath), ShouldBeTrue)

		So(os.RemoveAll(dirFullPath), ShouldBeNil)
	})
}
