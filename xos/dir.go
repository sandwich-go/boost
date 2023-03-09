package xos

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/sandwich-go/boost/xslice"
)

func MkdirAll(path string) error {
	return Mkdir(filepath.Dir(path))
}

func Mkdir(dir string) error {
	if FileExists(dir) {
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
		file, err := os.Open(path)
		if err != nil {
			return true
		}
		defer file.Close()
		names, err := file.Readdirnames(-1)
		if err != nil {
			return true
		}
		return len(names) == 0
	} else {
		return stat.Size() == 0
	}
}

// RemoveSubDirsUnderDir 删除指定目录下的子目录
func RemoveSubDirsUnderDir(dir string, filter func(dir string) bool) error {
	if !DirExists(dir) {
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

func removEmptyDir(path string, info os.FileInfo) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	if len(files) != 0 {
		return
	}

	err = os.Remove(path)
	if err != nil {
		panic(err)
	}
}

// RemoveEmptyDirs 删除空目录
func RemoveEmptyDirs(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			removEmptyDir(path, info)
		}
		return nil
	})
}

// RemoveDirs 删除指定目录
func RemoveDirs(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
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
			if xslice.ContainString(ext, filepath.Ext(path)) {
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
				if xslice.ContainString(ext, filepath.Ext(path)) {
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
			if xslice.ContainString(ext, filepath.Ext(path)) {
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
	return filepath.Walk(GetActuallyDir(root), walkFn)
}

// GetActuallyDir 获取真实目录，如果root为一个link则寻找link的真实目录
func GetActuallyDir(root string) string {
	dirInfo, err := os.Lstat(root)
	if err != nil {
		return root
	}
	if dirInfo.Mode()&os.ModeSymlink != 0 {
		dirLinkTo, err := os.Readlink(root)
		if err != nil {
			return root
		}
		return dirLinkTo
	}
	return root
}
