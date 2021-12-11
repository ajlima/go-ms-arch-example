package util

import (
	"time"

	"github.com/sirupsen/logrus"
)

func TrackTime(logger *logrus.Logger, start time.Time, message string) {
	logger.Printf(message, time.Since(start))
}
