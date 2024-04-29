package xos

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sandwich-go/boost/xpanic"
)

const PathSeparator = string(os.PathSeparator)

var goRunOnce sync.Once
var isGoRun bool

// IsGoRun returns true if the binary is run from a go run command.
func IsGoRun() bool {
	goRunOnce.Do(
		func() {
			ex, _ := os.Executable()
			exPath := filepath.Dir(ex)
			isGoRun = strings.Contains(exPath, "go-build") || strings.Contains(ex, "go_build_") || strings.Contains(exPath, "/private/var/folders/")
		},
	)
	return isGoRun
}

// MustGetBinaryFilePath returns binary path if the binary is run from a go run command.
func MustGetBinaryFilePath() (ret string) {
	ex, err := os.Executable()
	xpanic.WhenErrorAsFmtFirst(err, "Executable got error:%w")
	realPath, err0 := filepath.EvalSymlinks(ex)
	xpanic.WhenErrorAsFmtFirst(err0, "EvalSymlinks got error:%w")
	return realPath
}
