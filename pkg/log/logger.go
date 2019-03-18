package log

import (
	"fmt"
)

// Logger is a interface for creating standard logging calls.
// This will enable depended code  log.Debug("Woot")
// Avoid all other packages from declaring their own loggers
// This strategy enables simple change of the backend logger
type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

type Level uint32

const (
	WarnLevel  Level = 10
	FatalLevel Level = 20
	DebugLevel Level = 50
)

func Debug(message ...interface{}) {
	fmt.Println(message)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(message ...interface{}) {
	fmt.Println(message)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(message ...interface{}) {
	fmt.Println(message)
}
