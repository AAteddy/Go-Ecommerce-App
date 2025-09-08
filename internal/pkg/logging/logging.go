package logging

import (
	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger for structured logging.
type Logger struct {
	*logrus.Logger
}

// Init creates a new Logger instance with JSON formatting.
func Init() *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	return &Logger{log}
}

// WithFields adds fields to the logger.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{l.Logger.WithFields(fields).Logger}
}
