package xos

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var goRunOnce sync.Once
var isGoRun bool

//IsGoRun returns true if the binary is run from a go run command.
func IsGoRun() bool {
	goRunOnce.Do(
		func() {
			ex, _ := os.Executable()
			exPath := filepath.Dir(ex)
			isGoRun = strings.Contains(exPath, "go-build") || strings.Contains(ex, "_go_build_") || strings.Contains(exPath, "/private/var/folders/")
		},
	)
	return isGoRun
}
