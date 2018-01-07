package logger

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// Log is the singleton logger for the platform
type Log struct {
	*logrus.Logger
}

var log *Log

// New initializes the singleton logger
func New(prod, debug bool) error {
	if log != nil {
		return errors.New("logger already initialised")
	}
	log = &Log{logrus.New()}
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	if prod {
		log.Infoln("Running in PRODUCTION mode")
	}
	if debug {
		log.Infoln("Running in DEBUG mode")
		log.Level = logrus.DebugLevel
	}
	return nil
}

// Get returns the singleton instance of the logger
func Get() *Log {
	if log == nil {
		panic("logger not initialized")
	}

	return log
}
