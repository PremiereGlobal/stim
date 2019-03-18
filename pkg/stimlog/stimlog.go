package stimlog

import (
	"github.com/readytalk/stim/pkg/log"
	logurs "github.com/sirupsen/logrus"
)

type StimLogger struct {
	setLogger    log.Logger
	currentLevel log.Level
}

var cl *StimLogger = &StimLogger{
	setLogger:    logurs.New(),
	currentLevel: log.WarnLevel,
}

func GetLogger() *StimLogger {
	return cl
}

// SetLogger takes a structured logger to interface with.
// After the logger is setup it will be available across your packages
// If SetLogger is not used Debug will not create output
func (stimLogger *StimLogger) SetLogger(givenLogger log.Logger) {
	stimLogger.setLogger = givenLogger
}

func (stimLogger *StimLogger) SetLevel(level log.Level) {
	stimLogger.currentLevel = level
}

// Debug logs a message at level Debug on the standard logger.
func (stimLogger *StimLogger) Debug(message ...interface{}) {
	if stimLogger.currentLevel >= log.DebugLevel {
		if stimLogger.setLogger != nil {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

// Warn logs a message at level Warn on the standard logger.
func (stimLogger *StimLogger) Warn(message ...interface{}) {
	if stimLogger.currentLevel >= log.WarnLevel {
		if stimLogger.setLogger != nil {
			stimLogger.setLogger.Warn(message...)
		}
	}
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (stimLogger *StimLogger) Fatal(message ...interface{}) {
	if stimLogger.currentLevel >= log.FatalLevel {
		if stimLogger.setLogger != nil {
			stimLogger.setLogger.Fatal(message...)
		}
	}
}
