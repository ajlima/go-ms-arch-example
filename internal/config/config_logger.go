package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

func ConfigureLogger(log *logrus.Logger, fileName string, logLevel string) *logrus.Logger {
	// Configure log file
	if fileName != "stdout" {
		logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Panicf("Cannot create a log file: %s", fileName)
		}
		log.Out = logFile
	} else {
		log.Out = os.Stdout
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Panicf("Cannot configure log level to %s", logLevel)
	}
	log.Level = level

	return log
}
