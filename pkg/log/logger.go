package log

// Logging interface to create standard calls like: log.Debug("Woot")
// Avoid all other packages from declaring their own loggers
// This strategy enables simple change of the backend logger
type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

var logger Logger

// Pass your logger object to this function
// After the logger is setup it will be available across your packages
// If SetLogger is not used Debug will not create output
func SetLogger(givenLogger Logger) {
	logger = givenLogger
}

func Debug(message ...interface{}) {
	if logger != nil {
		logger.Debug(message...)
	}
}

func Warn(message ...interface{}) {
	logger.Warn(message...)
}

func Fatal(message ...interface{}) {
	logger.Fatal(message...)
}
