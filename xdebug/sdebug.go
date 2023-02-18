package xdebug

var isDebugEnabled = false

func Enabled() bool {
	return isDebugEnabled
}

func SetEnabled(enable bool) { isDebugEnabled = enable }
