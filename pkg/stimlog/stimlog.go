package stimlog

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cornelk/hashmap"
	logurs "github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

// StimLogger this struct is a generic logger used by stim packages
type StimLogger interface {
	Debug(...interface{})
	Verbose(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
	SetLogger(Logger)
	SetLevel(Level)
	SetDateFormat(string)
	AddLogFile(string, Level) error
	ForceFlush(bool)
}

// Level is the Level of logging set in stim
type Level int32

const (
	defaultLevel Level = -1
	//FatalLevel this is used to log an error that will cause fatal problems in the program
	FatalLevel Level = 0
	//WarnLevel is logging for interesting events that need to be known about but are not crazy
	WarnLevel    Level = 20
	InfoLevel    Level = 30
	VerboseLevel Level = 40
	//DebugLevel is used to debugging certain calls in Stim to see what is going on, usually only used for development
	DebugLevel Level = 50
)

type logFile struct {
	path     string
	logLevel Level
	fp       *os.File
}

type logMessage struct {
	logLevel Level
	msg      string
}

type stimLogger struct {
	setLogger    Logger
	currentLevel Level
	highestLevel Level
	dateFMT      string
	logfiles     hashmap.HashMap
	logQueue     chan *logMessage
	forceFlush   bool
}

var logger *stimLogger

const debugMsg = "[ DEBUG ]"
const warnMsg = "[ WARN  ]"
const fatalMsg = "[ FATAL ]"
const infoMsg = "[ INFO  ]"
const verboseMsg = "[VERBOSE]"
const dateFMT = "2006-01-02 15:04:05.9999999"
const subSTR = "{}"

//GetLogger gets a logger for logging in stim.
func GetLogger() StimLogger {
	if logger == nil {
		mu := sync.Mutex{}
		mu.Lock()
		if logger == nil {
			lg := logurs.New()
			logger = &stimLogger{
				// setLogger:    lg,
				currentLevel: WarnLevel,
				highestLevel: WarnLevel,
				dateFMT:      dateFMT,
				logQueue:     make(chan *logMessage, 20),
				logfiles:     hashmap.HashMap{},
			}
			//We set logurs to debug since we are handling the filtering
			lg.SetLevel(logurs.DebugLevel)
			logger.AddLogFile("STDOUT", defaultLevel)
			go logger.writeLogQueue()
			mu.Unlock()
		}
	}
	return logger
}

func (stimLogger *stimLogger) AddLogFile(file string, logLevel Level) error {
	var fp *os.File
	var err error
	if file == "STDOUT" {
		fp = os.Stdout
	} else if file == "STDERR" {
		fp = os.Stderr
	} else {
		fp, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0750)
		if err != nil {
			return err
		}
		fs, err := fp.Stat()
		if err != nil {
			return err
		}
		fp.Seek(fs.Size(), 0)
	}
	if logLevel > stimLogger.highestLevel {
		stimLogger.highestLevel = logLevel
	}
	stimLogger.logfiles.Set(file, &logFile{path: file, logLevel: logLevel, fp: fp})
	return nil
}

func (stimLogger *stimLogger) writeLogQueue() {
	for {
		logmsg := <-stimLogger.logQueue
		stimLogger.writeLogs(logmsg)
	}
}

func (stimLogger *stimLogger) writeLogs(lm *logMessage) {
	for kv := range stimLogger.logfiles.Iter() {
		lgr := kv.Value.(*logFile)
		if lgr.logLevel >= lm.logLevel || (lgr.logLevel == defaultLevel && stimLogger.currentLevel >= lm.logLevel) {
			lgr.fp.WriteString(lm.msg)
			if stimLogger.forceFlush {
				lgr.fp.Sync()
			}
		}
	}
}

func (stimLogger *stimLogger) formatAndLog(ll Level, level string, args ...interface{}) {
	stimLogger.logQueue <- stimLogger.formatString(ll, level, args...)
}

func (stimLogger *stimLogger) formatString(ll Level, level string, args ...interface{}) *logMessage {
	var msg string
	switch args[0].(type) {
	case string:
		msg = args[0].(string)
	default:
		msg = fmt.Sprintf("%v", args[0])
	}
	subs := strings.Split(msg, subSTR)
	var sb strings.Builder
	sb.WriteString(time.Now().Format(dateFMT))
	sb.WriteString("\t")
	sb.WriteString(level)
	sb.WriteString("\t")
	for i, v := range subs {
		v = strings.Replace(v, "{{", "{", -1)
		v = strings.Replace(v, "}}", "}", -1)
		sb.WriteString(v)
		if i < len(args)-1 {
			sb.WriteString(fmt.Sprintf("%v", args[i+1]))
		}
	}
	sb.WriteString("\n")
	return &logMessage{msg: sb.String(), logLevel: ll}
}

func (stimLogger *stimLogger) ForceFlush(ff bool) {
	stimLogger.forceFlush = ff
}

func (stimLogger *stimLogger) SetDateFormat(string) {

}

// SetLogger takes a structured logger to interface with.
// After the logger is setup it will be available across your packages
// If SetLogger is not used Debug will not create output
func (stimLogger *stimLogger) SetLogger(givenLogger Logger) {
	stimLogger.setLogger = givenLogger
}

// SetLevel sets the StimLogger log level.
func (stimLogger *stimLogger) SetLevel(level Level) {
	stimLogger.currentLevel = level
}

// Debug logs a message at level Debug on the standard logger.
func (stimLogger *stimLogger) Debug(message ...interface{}) {
	if stimLogger.highestLevel >= DebugLevel {
		if stimLogger.setLogger == nil {
			if stimLogger.forceFlush {
				stimLogger.writeLogs(stimLogger.formatString(DebugLevel, debugMsg, message...))
			} else {
				stimLogger.formatAndLog(DebugLevel, debugMsg, message...)
			}
		} else {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

// Debug logs a message at level Debug on the standard logger.
func (stimLogger *stimLogger) Verbose(message ...interface{}) {
	if stimLogger.highestLevel >= VerboseLevel {
		if stimLogger.setLogger == nil {
			if stimLogger.forceFlush {
				stimLogger.writeLogs(stimLogger.formatString(VerboseLevel, verboseMsg, message...))
			} else {
				stimLogger.formatAndLog(VerboseLevel, verboseMsg, message...)
			}
		} else {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

// Warn logs a message at level Warn on the standard logger.
func (stimLogger *stimLogger) Warn(message ...interface{}) {
	if stimLogger.highestLevel >= WarnLevel {
		if stimLogger.setLogger == nil {
			if stimLogger.forceFlush {
				stimLogger.writeLogs(stimLogger.formatString(WarnLevel, warnMsg, message...))
			} else {
				stimLogger.formatAndLog(WarnLevel, warnMsg, message...)
			}
		} else {
			stimLogger.setLogger.Warn(message...)
		}
	}
}

// Warn logs a message at level Warn on the standard logger.
func (stimLogger *stimLogger) Info(message ...interface{}) {
	if stimLogger.highestLevel >= InfoLevel {
		if stimLogger.setLogger == nil {
			if stimLogger.forceFlush {
				stimLogger.writeLogs(stimLogger.formatString(InfoLevel, infoMsg, message...))
			} else {
				stimLogger.formatAndLog(InfoLevel, infoMsg, message...)
			}
		} else {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (stimLogger *stimLogger) Fatal(message ...interface{}) {
	if stimLogger.highestLevel >= FatalLevel {
		if stimLogger.setLogger == nil {
			stimLogger.writeLogs(stimLogger.formatString(FatalLevel, fatalMsg, message...))
		} else {
			stimLogger.setLogger.Fatal(message...)
		}
		os.Exit(5)
	}
}
