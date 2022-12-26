//go:build !windows
// +build !windows

package xos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsHidden 判断是否为隐藏文件
func IsHidden(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if strings.HasPrefix(f.Name(), ".") {
		return true, nil
	}
	return false, nil
}

func hide(path string, hidden, force bool) (string, error) {
	isHidden, err := IsHidden(path)
	if err != nil || hidden == isHidden {
		return path, err
	}
	var newPath string
	if hidden {
		newPath = filepath.Join(filepath.Dir(path), "."+filepath.Base(path))
	} else {
		newPath = filepath.Join(filepath.Dir(path), strings.TrimPrefix(filepath.Base(path), "."))
	}
	if !force {
		_, err = os.Stat(newPath)
		if err == nil {
			return path, fmt.Errorf("\"%s\" already exists\nUse the `ForceHide/ForceUnHide` to skip this check", newPath)
		}
	}
	err = os.Rename(path, newPath)
	if err != nil {
		return path, err
	}
	return newPath, nil
}

// Hide 隐藏指定文件，返回被隐藏文件的名称
func Hide(path string) (string, error) {
	return hide(path, true, false)
}

// ForceHide 强制隐藏指定文件，返回被隐藏文件的名称，若已存在相同名称的隐藏文件，则会被覆盖
func ForceHide(path string) (string, error) {
	return hide(path, true, true)
}

// UnHide 取消隐藏指定文件，返回被取消隐藏文件的名称
func UnHide(path string) (string, error) {
	return hide(path, false, false)
}

// ForceUnHide 强制取消隐藏指定文件，返回被取消隐藏文件的名称，若已存在相同名称的取消隐藏文件，则会被覆盖
func ForceUnHide(path string) (string, error) {
	return hide(path, false, true)
}
