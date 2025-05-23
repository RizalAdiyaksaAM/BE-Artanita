package logger

import (
	"strings"

	"gorm.io/gorm/logger"
)

func GetDBLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	}

	return logger.Silent
}
