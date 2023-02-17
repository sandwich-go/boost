package boost

import "github.com/sandwich-go/boost/internal/log"

// InstallLogger 安装logger
func InstallLogger(logger log.Logger) {
	log.GlobalLogger = logger
}

// LogErrorAndEatError 输出 err
func LogErrorAndEatError(err error) {
	if err != nil {
		log.Error(err.Error())
	}
}

// LogDebug 输出 info
func LogDebug(msg string) { log.Debug(msg) }

// LogInfo 输出 info
func LogInfo(msg string) { log.Info(msg) }

// LogIWarn 输出 info
func LogIWarn(msg string) { log.Warn(msg) }

// LogError 输出 info
func LogError(msg string) { log.Error(msg) }

// LogFatal 输出 info
func LogFatal(msg string) { log.Fatal(msg) }
