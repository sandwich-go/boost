//go:build windows
// +build windows

package xos

import (
	"os"
	"syscall"
)

func getFileAttrs(path string) (uint32, *uint16, error) {
	utf16PtrPath, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, nil, err
	}
	var attrs uint32
	attrs, err = syscall.GetFileAttributes(utf16PtrPath)
	return attrs, utf16PtrPath, err
}

// IsHidden 判断是否为隐藏文件
func IsHidden(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	var attrs uint32
	attrs, _, err = getFileAttrs(path)
	if err != nil {
		return false, err
	}
	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN > 0, nil
}

func hide(path string, hidden bool) error {
	isHidden, err := IsHidden(path)
	if err != nil || hidden == isHidden {
		return err
	}
	var attrs uint32
	var utf16PtrPath *uint16
	attrs, utf16PtrPath, err = getFileAttrs(path)
	if err != nil {
		return err
	}

	var newAttrs uint32
	if hidden {
		if attrs&syscall.FILE_ATTRIBUTE_HIDDEN > 0 {
			return nil
		}
		newAttrs = attrs | syscall.FILE_ATTRIBUTE_HIDDEN
	} else {
		if attrs&syscall.FILE_ATTRIBUTE_HIDDEN == 0 {
			return nil
		}
		newAttrs = attrs - (attrs & syscall.FILE_ATTRIBUTE_HIDDEN)
	}
	return syscall.SetFileAttributes(utf16PtrPath, newAttrs)
}

// Hide 隐藏指定文件，返回被隐藏文件的名称
func Hide(path string) (string, error) {
	return path, hide(path, true)
}

// ForceHide 隐藏指定文件，返回被隐藏文件的名称
func ForceHide(path string) (string, error) {
	return Hide(path)
}

// UnHide 取消隐藏指定文件，返回被取消隐藏文件的名称
func UnHide(path string) (string, error) {
	return path, hide(path, false)
}

// ForceUnHide 取消隐藏指定文件，返回被取消隐藏文件的名称
func ForceUnHide(path string) (string, error) {
	return UnHide(path)
}
