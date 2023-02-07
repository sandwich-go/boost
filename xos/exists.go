package xos

import (
	"errors"
	"os"
)

// ExistsTreatErrorAsExist Exists函数组获取是否存在时，如果发生了非os.ErrNotExist错误，视作存在或者不存在，减轻逻辑层判断负担
var ExistsTreatErrorAsExist = true

// Exists 指定的文件或者目录是否存在
// 如果发生了非 os.ErrNotExist 错误，则认为存在
func Exists(fileOrDirPath string) bool {
	_, err := os.Stat(fileOrDirPath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return ExistsTreatErrorAsExist
}

// ExistsFile 给定的 fileName 是否存在且是一个文件,如果 fileName 存在但是是一个目录也会返回 false
// 如果发生了非 os.ErrNotExist 错误，则认为存在
func ExistsFile(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return ExistsTreatErrorAsExist
	}
	return !info.IsDir()
}

// ExistsDir 给定的 filePath 是否存在且是一个目录,如果 filePath 存在但是是一个文件也会返回错误
// 如果发生了非 os.ErrNotExist 错误，则认为存在
func ExistsDir(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return ExistsTreatErrorAsExist
	}
	return info.IsDir()
}
