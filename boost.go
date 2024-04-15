package boost

import (
	"fmt"

	"github.com/sandwich-go/boost/internal/log"
)

// InstallLogger 安装 log.Logger
func InstallLogger(logger log.Logger) {
	log.Default = logger
}

// LogErrorAndEatError 输出 err
func LogErrorAndEatError(err error) {
	if err != nil {
		log.Error(err.Error())
	}
}

// LogDebug 输出 info
func LogDebug(msg string)                       { log.Debug(msg) }
func LogDebugf(format string, a ...interface{}) { LogDebug(fmt.Sprintf(format, a...)) }

// LogInfo 输出 info
func LogInfo(msg string)                       { log.Info(msg) }
func LogInfof(format string, a ...interface{}) { LogInfo(fmt.Sprintf(format, a...)) }

// LogWarn 输出 info
func LogWarn(msg string)                       { log.Warn(msg) }
func LogWarnf(format string, a ...interface{}) { LogWarn(fmt.Sprintf(format, a...)) }

// LogError 输出 info
func LogError(msg string)                       { log.Error(msg) }
func LogErrorf(format string, a ...interface{}) { LogError(fmt.Sprintf(format, a...)) }

// LogFatal 输出 info
func LogFatal(msg string)                       { log.Fatal(msg) }
func LogFatalf(format string, a ...interface{}) { LogFatal(fmt.Sprintf(format, a...)) }
