package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestDir(t *testing.T) {
	Convey("test dir", t, func() {
		var dest = filepath.Join(os.TempDir(), "test_perf")
		_ = os.RemoveAll(dest)

		So(Exists(dest), ShouldBeFalse)

		So(Mkdir(dest), ShouldBeNil)
		So(Exists(dest), ShouldBeTrue)

		So(os.RemoveAll(dest), ShouldBeNil)

		So(MkdirAll(filepath.Join(dest, "a.go")), ShouldBeNil)
		So(Exists(dest), ShouldBeTrue)

		So(os.RemoveAll(dest), ShouldBeNil)

		So(IsEmpty(dest), ShouldBeTrue)
		So(Mkdir(dest), ShouldBeNil)
		So(Exists(dest), ShouldBeTrue)
		So(IsEmpty(dest), ShouldBeTrue)

		var destFile = filepath.Join(dest, "a.go")
		So(FilePutContents(destFile, nil), ShouldBeNil)
		So(IsEmpty(destFile), ShouldBeTrue)

		// 不存在的子目录，不会报错
		So(RemoveSubDirsUnderDir(filepath.Join(dest, "sub"), nil), ShouldBeNil)

		var destDir1, destDir2 = filepath.Join(dest, "a"), filepath.Join(dest, "b")
		So(Mkdir(destDir1), ShouldBeNil)
		So(Mkdir(destDir2), ShouldBeNil)
		So(FilePutContents(filepath.Join(destDir1, "a.go"), nil), ShouldBeNil)
		So(FilePutContents(filepath.Join(destDir1, "b.go"), nil), ShouldBeNil)

		So(RemoveSubDirsUnderDir(dest, func(dir string) bool { return dir == destDir1 }), ShouldBeNil)
		So(Exists(destDir1), ShouldBeFalse)

		So(Mkdir(destDir1), ShouldBeNil)
		So(FilePutContents(filepath.Join(destDir1, "a.go"), nil), ShouldBeNil)
		So(FilePutContents(filepath.Join(destDir1, "b.go"), nil), ShouldBeNil)
		So(RemoveEmptyDirs(destDir1), ShouldNotBeNil)
		So(Exists(destDir1), ShouldBeTrue)
		So(RemoveDirs(destDir1), ShouldBeNil)
		So(Exists(destDir1), ShouldBeTrue)
		So(IsEmpty(destDir1), ShouldBeTrue)

		So(FilePutContents(filepath.Join(destDir1, "a.go"), nil), ShouldBeNil)
		So(FilePutContents(filepath.Join(destDir1, "b.go"), nil), ShouldBeNil)
		RemoveFilesUnderDir(destDir1, func(f string) bool { return f == filepath.Join(destDir1, "a.go") })
		var fileList []string
		So(filepath.Walk(destDir1, FileWalkFunc(&fileList, ".txt")), ShouldBeNil)
		So(len(fileList), ShouldBeZeroValue)
		So(filepath.Walk(destDir1, FileWalkFunc(&fileList, ".go")), ShouldBeNil)
		So(len(fileList), ShouldEqual, 1)
		So(fileList[0], ShouldEqual, filepath.Join(destDir1, "b.go"))
		fileList = nil
		So(filepath.Walk(destDir1, FileWalkFuncWithIncludeFilter(&fileList, func(f string) bool { return f == filepath.Join(destDir1, "a.go") }, ".go")), ShouldBeNil)
		So(len(fileList), ShouldBeZeroValue)
		fileList = nil
		So(filepath.Walk(destDir1, FileWalkFuncWithIncludeFilter(&fileList, func(f string) bool { return f == filepath.Join(destDir1, "b.go") }, ".go")), ShouldBeNil)
		So(len(fileList), ShouldEqual, 1)
		So(fileList[0], ShouldEqual, filepath.Join(destDir1, "b.go"))
		fileList = nil
		So(filepath.Walk(destDir1, FileWalkFuncWithExcludeFilter(&fileList, func(f string) bool { return f == filepath.Join(destDir1, "b.go") }, ".go")), ShouldBeNil)
		So(len(fileList), ShouldBeZeroValue)
		fileList = nil
		So(filepath.Walk(destDir1, FileWalkFuncWithExcludeFilter(&fileList, func(f string) bool { return f == filepath.Join(destDir1, "a.go") }, ".go")), ShouldBeNil)
		So(len(fileList), ShouldEqual, 1)
		So(fileList[0], ShouldEqual, filepath.Join(destDir1, "b.go"))
		fileList = nil

		So(FilePathWalkFollowLink(destDir1, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				fileList = append(fileList, path)
			}
			return err
		}), ShouldBeNil)
		So(len(fileList), ShouldEqual, 1)
		So(fileList[0], ShouldEqual, filepath.Join(destDir1, "b.go"))

		So(FilePutContents(filepath.Join(destDir1, "c.txt"), nil), ShouldBeNil)

		fileNames, err := ReadDir(destDir1)
		So(err, ShouldBeNil)
		So(len(fileNames), ShouldEqual, 2)

		fileNames, err = ReadDirWithExt(destDir1, ".go")
		So(err, ShouldBeNil)
		So(len(fileNames), ShouldEqual, 1)
		So(fileNames[0], ShouldEqual, "b.go")

		So(IsDirWriteable(destDir1), ShouldBeNil)
		So(TouchDirAll(filepath.Join(destDir1, "c.txt")), ShouldNotBeNil)
		So(TouchDirAll(destDir1), ShouldBeNil)
		So(CreateDirAll(destDir1), ShouldNotBeNil)

		So(os.RemoveAll(dest), ShouldBeNil)
	})
}
