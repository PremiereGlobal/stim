package log

// Logger is a interface for creating standard logging calls.
// This will enable depended code  log.Debug("Woot")
// Avoid all other packages from declaring their own loggers
// This strategy enables simple change of the backend logger
type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

var logger Logger

// SetLogger takes a structured logger to interface with.
// After the logger is setup it will be available across your packages
// If SetLogger is not used Debug will not create output
func SetLogger(givenLogger Logger) {
	logger = givenLogger
}

// Debug logs a message at level Debug on the standard logger.
func Debug(message ...interface{}) {
	if logger != nil {
		logger.Debug(message...)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(message ...interface{}) {
	logger.Warn(message...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(message ...interface{}) {
	logger.Fatal(message...)
}
