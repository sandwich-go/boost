package xexec

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xstrings"
)

func envGet(key string, def ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v
}

// getShell returns the shell command depending on current working operation system.
// It returns "cmd.exe" for windows, and "bash" or "sh" for others.
func getShell() string {
	switch runtime.GOOS {
	case "windows":
		return SearchBinary("cmd.exe")
	default:
		// Check the default binary storage path.
		if xos.FileExists("/bin/bash") {
			return "/bin/bash"
		}
		if xos.FileExists("/bin/sh") {
			return "/bin/sh"
		}
		// Else search the env PATH.
		path := SearchBinary("bash")
		if path == "" {
			path = SearchBinary("sh")
		}
		return path
	}
}

func SearchBinary(file string) string {
	// Check if it's absolute path of exists at current working directory.
	if xos.FileExists(file) {
		return file
	}
	return SearchBinaryPath(file)
}

func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

// SearchBinaryPath searches the binary <file> in PATH environment.
func SearchBinaryPath(file string) string {
	array := ([]string)(nil)
	switch runtime.GOOS {
	case "windows":
		envPath := envGet("PATH", envGet("Path"))
		if strings.Contains(envPath, ";") {
			array = xstrings.SplitAndTrim(envPath, ";")
		} else if strings.Contains(envPath, ":") {
			array = xstrings.SplitAndTrim(envPath, ":")
		}
		if Ext(file) != ".exe" {
			file += ".exe"
		}
	default:
		array = xstrings.SplitAndTrim(envGet("PATH"), ":")
	}
	if len(array) > 0 {
		path := ""
		for _, v := range array {
			path = v + string(filepath.Separator) + file
			if xos.FileExists(path) {
				return path
			}
		}
	}
	return ""
}
