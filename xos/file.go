package xos

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sandwich-go/boost/xpanic"
	"github.com/sandwich-go/boost/xslice"
)

// Ext 返回后缀
func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

// FileExists 给定的filename是否存在且是一个文件
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
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

// FileGetContents 获取文件内容
func FileGetContents(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// MustFilePutContents 写入文件，如果发生错误则panic
func MustFilePutContents(filename string, content []byte) {
	dirName := filepath.Dir(filename)
	xpanic.PanicIfErrorAsFmtFirst(os.MkdirAll(dirName, os.ModePerm), "got error:%w while MkdirAll with:%s", dirName)
	xpanic.PanicIfErrorAsFmtFirst(ioutil.WriteFile(filename, content, 0644), "got error:%w while WriteFile with:%s", filename)
}

// FilePutContents 写入文件
func FilePutContents(filename string, content []byte) error {
	dirName := filepath.Dir(filename)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, content, 0644)
}
