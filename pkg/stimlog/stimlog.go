package stimlog

import (
	"sync"

	"github.com/readytalk/stim/pkg/log"
	logurs "github.com/sirupsen/logrus"
)

// StimLogger this struct is a generic logger used by stim packages
type StimLogger struct {
	setLogger    log.Logger
	currentLevel Level
}

// Level is the Level of logging set in stim
type Level uint32

const (
	//FatalLevel this is used to log an error that will cause fatal problems in the program
	FatalLevel Level = 0
	//WarnLevel is logging for interesting events that need to be known about but are not crazy
	WarnLevel Level = 20
	//DebugLevel is used to debugging certain calls in Stim to see what is going on, usually only used for development
	DebugLevel Level = 50
)

var logger *StimLogger

//GetLogger gets a logger for logging in stim.
func GetLogger() *StimLogger {
	if logger == nil {
		mu := sync.Mutex{}
		mu.Lock()
		if logger == nil {
			lg := logurs.New()
			logger = &StimLogger{
				setLogger:    lg,
				currentLevel: WarnLevel,
			}
			//We set logurs to debug since we are handling the filtering
			lg.SetLevel(logurs.DebugLevel)
		}
	}
	return logger
}

// SetLogger takes a structured logger to interface with.
// After the logger is setup it will be available across your packages
// If SetLogger is not used Debug will not create output
func (stimLogger *StimLogger) SetLogger(givenLogger log.Logger) {
	stimLogger.setLogger = givenLogger
}

// SetLevel sets the StimLogger log level.
func (stimLogger *StimLogger) SetLevel(level Level) {
	stimLogger.currentLevel = level
}

// Debug logs a message at level Debug on the standard logger.
func (stimLogger *StimLogger) Debug(message ...interface{}) {
	if stimLogger.currentLevel >= DebugLevel {
		if stimLogger.setLogger != nil {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

// Warn logs a message at level Warn on the standard logger.
func (stimLogger *StimLogger) Warn(message ...interface{}) {
	if stimLogger.currentLevel >= WarnLevel {
		if stimLogger.setLogger != nil {
			stimLogger.setLogger.Warn(message...)
		}
	}
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (stimLogger *StimLogger) Fatal(message ...interface{}) {
	if stimLogger.currentLevel >= FatalLevel {
		if stimLogger.setLogger != nil {
			stimLogger.setLogger.Fatal(message...)
		}
	}
}
