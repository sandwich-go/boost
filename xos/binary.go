package xos

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sandwich-go/boost/xstrings"
)

// SearchBinary 查找二进制文件
func SearchBinary(file string) string {
	// Check if it's absolute path of exists at current working directory.
	if FileExists(file) {
		return file
	}
	return SearchBinaryPath(file)
}

// SearchBinaryPath searches the binary <file> in PATH environment.
func SearchBinaryPath(file string) string {
	array := ([]string)(nil)
	switch runtime.GOOS {
	case "windows":
		envPath := EnvGet("PATH", EnvGet("Path"))
		if strings.Contains(envPath, ";") {
			array = xstrings.SplitAndTrim(envPath, ";")
		} else if strings.Contains(envPath, ":") {
			array = xstrings.SplitAndTrim(envPath, ":")
		}
		if Ext(file) != ".exe" {
			file += ".exe"
		}
	default:
		array = xstrings.SplitAndTrim(EnvGet("PATH"), ":")
	}
	if len(array) > 0 {
		path := ""
		for _, v := range array {
			path = v + string(filepath.Separator) + file
			if FileExists(path) {
				return path
			}
		}
	}
	return ""
}
