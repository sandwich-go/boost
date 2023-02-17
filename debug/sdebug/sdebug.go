package sdebug

import (
	"github.com/sandwich-go/boost/xcmd"
)

const (
	cmdEnvKeyForDebug = "sandwich_debug"
)

var isDebugEnabled = false

func init() {
	xcmd.AddFlag(cmdEnvKeyForDebug, xcmd.DefaultStringFalse)
	SetEnabled(xcmd.IsTrue(xcmd.GetOptWithEnv(cmdEnvKeyForDebug)))
}

// Enabled 获取是否开启了debug
func Enabled() bool {
	return isDebugEnabled
}

// SetEnabled 设置是否开启debug
func SetEnabled(enable bool) { isDebugEnabled = enable }
