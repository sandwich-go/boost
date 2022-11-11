package boost

import "github.com/sandwich-go/boost/internal/log"

// InstallLogger 安装logger
func InstallLogger(logger log.Logger) {
	log.GlobalLogger = logger
}

func LogErrorAndEatError(err error) {
	if err != nil {
		log.Error(err.Error())
	}
}
