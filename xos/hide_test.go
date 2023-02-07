package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func TestHide(t *testing.T) {
	Convey("hide", t, func() {
		var tmpFile = filepath.Join(os.TempDir(), "test_dir")
		_, err := IsHidden(tmpFile)
		So(err, ShouldNotBeNil)
		So(MkdirAll(tmpFile), ShouldBeNil)
		MustFilePutContents(tmpFile, []byte("a"))
		is, err1 := IsHidden(tmpFile)
		So(err1, ShouldBeNil)
		So(is, ShouldBeFalse)
		var newTmpFile string
		newTmpFile, err = Hide(tmpFile)
		So(err, ShouldBeNil)
		is, err = IsHidden(newTmpFile)
		So(err, ShouldBeNil)
		So(is, ShouldBeTrue)
		tmpFile, err = UnHide(newTmpFile)
		So(err, ShouldBeNil)
		is, err = IsHidden(tmpFile)
		So(err, ShouldBeNil)
		So(is, ShouldBeFalse)

		So(os.RemoveAll(tmpFile), ShouldBeNil)
	})
}
