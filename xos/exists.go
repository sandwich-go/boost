package xos

import (
	"errors"
	"os"
)

// Exists 指定的文件或者目录是否存在
func Exists(fileOrDirPath string) (bool, error) {
	_, err := os.Stat(fileOrDirPath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// FileExists 给定的filenName是否存在且是一个文件,如果filenName存在但是是一个目录也会返回false
func FileExists(filenName string) (bool, error) {
	info, err := os.Stat(filenName)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

// DirExists 给定的filePath是否存在且是一个目录,如果filePath存在但是是一个文件也会返回错误
func DirExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
