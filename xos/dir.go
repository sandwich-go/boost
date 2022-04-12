package xos

import (
	"os"
	"path/filepath"

	"github.com/sandwich-go/boost/xslice"
)

// IsEmpty 检测目录是否为空
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
