package common

// Logger is an interface for passing a custom logger
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
}
