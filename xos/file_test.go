package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func TestFile(t *testing.T) {
	Convey("file", t, func() {
		var file = "file.go"
		var dirFullPath0 = filepath.Join(os.TempDir(), "a")
		So(Mkdir(dirFullPath0), ShouldBeNil)
		So(FileCopyToDir(file, dirFullPath0), ShouldBeNil)
		So(ExistsFile(filepath.Join(dirFullPath0, file)), ShouldBeTrue)
		So(Ext(file), ShouldEqual, ".go")
		content0, err0 := FileGetContents(filepath.Join(dirFullPath0, file))
		So(err0, ShouldBeNil)
		content1, err1 := FileGetContents(file)
		So(err1, ShouldBeNil)
		So(content0, ShouldResemble, content1)
		var bFile = filepath.Join(dirFullPath0, "b")
		MustFilePutContents(bFile, []byte("a"))
		So(FilePutContents(bFile, []byte("b")), ShouldBeNil)
		content1, err1 = FileGetContents(bFile)
		So(err1, ShouldBeNil)
		So(content1, ShouldResemble, []byte("b"))
		writer, cancel := MustGetFileWriter(bFile, true)
		_, err1 = writer.Write([]byte("a"))
		So(err1, ShouldBeNil)
		cancel()
		content1, err1 = FileGetContents(bFile)
		So(err1, ShouldBeNil)
		So(content1, ShouldResemble, []byte("ab"))
		So(os.RemoveAll(dirFullPath0), ShouldBeNil)
	})
}
