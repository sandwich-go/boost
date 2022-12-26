package xos

import (
	"os"
	"path/filepath"
	"strings"
)

// IsHiddenOrInHiddenDir 判断指定文件是否为隐藏文件，或者包含在隐藏目录中
func IsHiddenOrInHiddenDir(path string) (bool, error) {
	ss := strings.Split(path, string(os.PathSeparator))
	var s string
	for _, v := range ss {
		s = filepath.Join(s, v)
		is, err := IsHidden(s)
		if err != nil || is {
			return is, err
		}
	}
	return false, nil
}
