package logger

import (
	"os"

	"github.com/arif-x/sqlx-mysql-boilerplate/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

var logger = &Logger{}

func SetUpLogger() {
	logger = &Logger{logrus.New()}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(os.Stdout)

	if config.AppCfg().Debug {
		logger.SetLevel(logrus.DebugLevel)
	}
}

func GetLogger() *Logger {
	return logger
}
