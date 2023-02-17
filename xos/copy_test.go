package xos

import (
	"github.com/sandwich-go/boost/xhash/md5"
	"github.com/sandwich-go/boost/xmap"
	. "github.com/smartystreets/goconvey/convey"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestCopy(t *testing.T) {
	Convey("copy a file", t, func() {
		var file = "binary.go"
		var dest = filepath.Join(os.TempDir(), file)
		err := Copy(file, dest)
		So(err, ShouldBeNil)
		So(ExistsFile(dest), ShouldBeTrue)
		var md50, md51 string
		md50, err = md5.File(file)
		So(err, ShouldBeNil)
		md51, err = md5.File(dest)
		So(err, ShouldBeNil)
		So(md50, ShouldEqual, md51)
		So(os.Remove(dest), ShouldBeNil)
	})

	Convey("copy a dir", t, func() {
		var files []string
		var md5s = make(map[string]string)
		var dir = "../xencoding/protobuf/test_perf"
		err := FilePathWalkFollowLink(dir, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, info.Name())
				md5s[info.Name()], _ = md5.File(path)
			}
			return nil
		})
		sort.Strings(files)
		So(err, ShouldBeNil)
		var destFiles []string
		var destMd5s = make(map[string]string)
		var dest = filepath.Join(os.TempDir(), "test_perf")
		err = Copy(dir, dest)
		So(err, ShouldBeNil)
		So(Exists(dest), ShouldBeTrue)
		err = FilePathWalkFollowLink(dir, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				destFiles = append(destFiles, info.Name())
				destMd5s[info.Name()], _ = md5.File(path)
			}
			return nil
		})
		sort.Strings(destFiles)
		So(files, ShouldResemble, destFiles)
		So(len(md5s), ShouldEqual, len(destMd5s))
		So(xmap.EqualStringStringMap(md5s, destMd5s), ShouldBeTrue)
		So(os.RemoveAll(dest), ShouldBeNil)
	})
}
