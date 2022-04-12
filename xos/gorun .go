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

func MustGetBinaryFilePath() (ret string) {
	ex, err := os.Executable()
	xpanic.PanicIfErrorAsFmtFirst(err, "Executable got error:%w")
	realPath, err := filepath.EvalSymlinks(ex)
	xpanic.PanicIfErrorAsFmtFirst(err, "EvalSymlinks got error:%w")
	return realPath
}
