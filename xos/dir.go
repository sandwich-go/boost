package xos

import (
	"fmt"
	"github.com/sandwich-go/boost"
	"github.com/sandwich-go/boost/xerror"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/sandwich-go/boost/xslice"
)

// MkdirAll 创建 path 的目录
func MkdirAll(path string) error {
	return Mkdir(filepath.Dir(path))
}

// Mkdir 创建目录
func Mkdir(dir string) error {
	if ExistsFile(dir) {
		return nil
	}
	return os.MkdirAll(dir, 0775)
}

// IsEmpty 检测目录是否为空,如目录下存在隐藏文件也会认为是非空目录
func IsEmpty(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return true
	}
	if stat.IsDir() {
		file, err0 := os.Open(path)
		if err0 != nil {
			return true
		}
		defer func() {
			_ = file.Close()
		}()
		names, err1 := file.Readdirnames(-1)
		if err1 != nil {
			return true
		}
		return len(names) == 0
	} else {
		return stat.Size() == 0
	}
}

// RemoveSubDirsUnderDir 删除指定目录下的子目录
func RemoveSubDirsUnderDir(dir string, filter func(dir string) bool) error {
	if !ExistsDir(dir) {
		return nil
	}
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	fileList := make([]string, 0)
	for _, f := range fs {
		if f.IsDir() {
			fd := filepath.Join(dir, f.Name())
			if filter == nil || filter(fd) {
				fileList = append(fileList, fd)
			}
		}
	}
	for _, filePath := range fileList {
		subFileList := make([]string, 0)
		_ = filepath.Walk(filePath, FileWalkFunc(&subFileList))
		// 排序，是为了先删除文件，再删除空目录
		sort.Sort(sort.Reverse(sort.StringSlice(subFileList)))
		for _, subFilePath := range subFileList {
			_ = os.Remove(subFilePath)
		}
	}
	return nil
}

func removeEmptyDir(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	if len(files) != 0 {
		return fmt.Errorf("not empty dir, %s", path)
	}
	return os.Remove(path)
}

// RemoveEmptyDirs 删除空目录
// 若目录不为空，则报错
func RemoveEmptyDirs(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = removeEmptyDir(path)
		}
		return err
	})
}

// RemoveDirs 删除指定目录下所有内容
func RemoveDirs(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer func() {
		_ = d.Close()
	}()
	names, err0 := d.Readdirnames(-1)
	if err0 != nil {
		return err0
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveFilesUnderDir 删除目录下的文件
func RemoveFilesUnderDir(pathStr string, filter func(filePath string) bool) {
	fileList := make([]string, 0)
	_ = filepath.Walk(pathStr, FileWalkFuncWithIncludeFilter(&fileList, filter))
	for _, filePath := range fileList {
		_ = os.Remove(filePath)
	}
}

// FileWalkFunc 只有指定的ext为合法文件
func FileWalkFunc(files *[]string, ext ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if len(ext) > 0 {
			if xslice.StringsContain(ext, filepath.Ext(path)) {
				*files = append(*files, path)
			}
		} else {
			*files = append(*files, path)
		}
		return err
	}
}

// FileWalkFuncWithIncludeFilter  filepath.WalkFunc 通过include过滤合法的文件,如指定了ext则只有ext扩展类型的文件合法
func FileWalkFuncWithIncludeFilter(files *[]string, include func(f string) bool, ext ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return err
		}
		if include == nil || include(path) {
			if len(ext) > 0 {
				if xslice.StringsContain(ext, filepath.Ext(path)) {
					*files = append(*files, path)
				}
			} else {
				*files = append(*files, path)
			}
		}
		return err
	}
}

// FileWalkFuncWithExcludeFilter filepath.WalkFunc 通过excluded过滤不合法文件,如指定了ext则只有ext扩展类型的文件合法
func FileWalkFuncWithExcludeFilter(files *[]string, excluded func(f string) bool, ext ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if excluded != nil && excluded(path) {
			return err
		}
		if len(ext) > 0 {
			if xslice.StringsContain(ext, filepath.Ext(path)) {
				*files = append(*files, path)
			}
		} else {
			*files = append(*files, path)
		}
		return err
	}
}

// FilePathWalkFollowLink 遍历目录
func FilePathWalkFollowLink(root string, walkFn filepath.WalkFunc) error {
	dir, err := GetActuallyDir(root)
	if err != nil {
		return err
	}
	return filepath.Walk(dir, walkFn)
}

// GetActuallyDir 获取真实目录，如果root为一个link则寻找link的真实目录
func GetActuallyDir(root string) (string, error) {
	dirInfo, err := os.Lstat(root)
	if err != nil {
		return root, err
	}
	if dirInfo.Mode()&os.ModeSymlink != 0 {
		dirLinkTo, err0 := os.Readlink(root)
		if err0 != nil {
			return root, err0
		}
		return dirLinkTo, nil
	}
	return root, nil
}

// ReadDir returns the filenames in the given directory in sorted order.
func ReadDir(root string) ([]string, error) {
	return ReadDirWithExt(root, "")
}

// ReadDirWithExt returns the filenames in the given directory in sorted order.
func ReadDirWithExt(root, ext string) ([]string, error) {
	dir, err := os.Open(root)
	if err != nil {
		return nil, err
	}
	defer func() {
		boost.LogErrorAndEatError(dir.Close())
	}()
	names, err0 := dir.Readdirnames(-1)
	if err0 != nil {
		return nil, err0
	}
	sort.Strings(names)

	if len(ext) > 0 {
		tss := make([]string, 0, len(names))
		for _, v := range names {
			if filepath.Ext(v) == ext {
				tss = append(tss, v)
			}
		}
		names = tss
	}
	return names, nil
}

// IsDirWriteable checks if dir is writable by writing and removing a file
// to dir. It returns nil if dir is writable.
func IsDirWriteable(dir string) error {
	f := filepath.Join(dir, ".touch")
	if err := ioutil.WriteFile(f, []byte(""), 0600); err != nil {
		return err
	}
	return os.Remove(f)
}

// TouchDirAll is similar to os.MkdirAll. It creates directories with 0700 permission if any directory
// does not exists. TouchDirAll also ensures the given directory is writable.
func TouchDirAll(dir string) error {
	// If path is already a directory, MkdirAll does nothing
	// and returns nil.
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		// if mkdirAll("a/text") and "text" is not
		// a directory, this will return syscall.ENOTDIR
		return err
	}
	return IsDirWriteable(dir)
}

// CreateDirAll is similar to TouchDirAll but returns error
// if the deepest directory was not empty.
func CreateDirAll(dir string) error {
	err := TouchDirAll(dir)
	if err == nil {
		var ns []string
		ns, err = ReadDir(dir)
		if err != nil {
			return err
		}
		if len(ns) != 0 {
			err = xerror.NewText("expected %q to be empty, got %q", dir, ns)
		}
	}
	return err
}
