package log

type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

var logger Logger

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
