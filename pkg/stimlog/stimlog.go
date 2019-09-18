package stimlog

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cornelk/hashmap"
)

type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

// StimLogger this struct is a generic logger used by stim packages
type StimLogger interface {
	Trace(...interface{})
	Debug(...interface{})
	Verbose(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
	GetLogLevel() Level
}
type StimLoggerConfig interface {
	SetLogger(Logger)
	SetLevel(Level)
	SetDateFormat(string)
	AddLogFile(string, Level) error
	RemoveLogFile(string)
	ForceFlush(bool)
	Flush()
	EnableLevelLogging(bool)
	EnableTimeLogging(bool)
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
	TraceLevel Level = 60
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

type FullStimLogger struct {
	setLogger    Logger
	currentLevel Level
	highestLevel Level
	dateFMT      string
	logfiles     hashmap.HashMap
	logQueue     []*logMessage
	forceFlush   bool
	logLevel     bool
	logTime      bool
	wqc          *sync.Cond
}

var logger *FullStimLogger

var prefixLogger map[string]StimLogger

const traceMsg = "[ TRACE ]"
const debugMsg = "[ DEBUG ]"
const warnMsg = "[ WARN  ]"
const fatalMsg = "[ FATAL ]"
const infoMsg = "[ INFO  ]"
const verboseMsg = "[VERBOSE]"
const dateFMT = "2006-01-02 15:04:05.9999999"
const subSTR = "{}"

var stimLoggerCreateLock sync.Mutex = sync.Mutex{}

func resetLogger() {
	logger = nil
}

func GetLoggerConfig() StimLoggerConfig {
	GetLogger()
	return logger
}

//GetLogger gets a logger for logging in stim.
func GetLogger() StimLogger {
	if logger == nil {
		stimLoggerCreateLock.Lock()
		defer stimLoggerCreateLock.Unlock()
		if logger == nil {
			logger = &FullStimLogger{
				currentLevel: InfoLevel,
				highestLevel: InfoLevel,
				dateFMT:      dateFMT,
				logQueue:     make([]*logMessage, 0),
				logfiles:     hashmap.HashMap{},
				forceFlush:   true,
				logLevel:     true,
				logTime:      true,
				// wql:          l,
				wqc: sync.NewCond(&sync.Mutex{}),
			}
			//We set logurs to debug since we are handling the filtering
			logger.AddLogFile("STDOUT", defaultLevel)
			go logger.writeLogQueue()
		}
	}
	return logger
}

//GetLogger gets a logger for logging in stim.
func GetLoggerWithPrefix(prefix string) StimLogger {
	if prefix == "" {
		return GetLogger()
	}
	if prefixLogger == nil {
		stimLoggerCreateLock.Lock()
		if prefixLogger == nil {
			prefixLogger = make(map[string]StimLogger)
		}
		stimLoggerCreateLock.Unlock()
	}
	stimLoggerCreateLock.Lock()
	defer stimLoggerCreateLock.Unlock()
	if sl, ok := prefixLogger[prefix]; ok {
		return sl
	}
	prefixLogger[prefix] = &StimPrefixLogger{stimLogger: GetLogger(), prefix: prefix}
	return prefixLogger[prefix]
}

func (stimLogger *FullStimLogger) GetLogLevel() Level {
	return stimLogger.highestLevel
}

func (stimLogger *FullStimLogger) EnableLevelLogging(b bool) {
	stimLogger.logLevel = b
}

func (stimLogger *FullStimLogger) EnableTimeLogging(b bool) {
	stimLogger.logTime = b
}

func (stimLogger *FullStimLogger) RemoveLogFile(file string) {
	_, ok := stimLogger.logfiles.Get(file)
	if ok {
		highestLL := defaultLevel
		stimLogger.logfiles.Del(file)
		for kv := range stimLogger.logfiles.Iter() {
			lgr := kv.Value.(*logFile)
			if lgr.logLevel > highestLL {
				highestLL = lgr.logLevel
			}
		}
		if highestLL > defaultLevel {
			stimLogger.highestLevel = highestLL
		}
	}
}

func (stimLogger *FullStimLogger) AddLogFile(file string, logLevel Level) error {
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

func (stimLogger *FullStimLogger) writeLogQueue() {
	stimLogger.wqc.L.Lock()
	defer stimLogger.wqc.L.Unlock()
	for {
		if len(stimLogger.logQueue) > 0 {
			for len(stimLogger.logQueue) > 0 {
				wl := stimLogger.logQueue[0]
				stimLogger.writeLogs(wl)
				stimLogger.logQueue = stimLogger.logQueue[1:]
			}
		} else {
			stimLogger.wqc.Wait()
		}

	}
}

func (stimLogger *FullStimLogger) writeLogs(lm *logMessage) {
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

func (stimLogger *FullStimLogger) formatAndLog(ll Level, level string, args ...interface{}) {
	stimLogger.wqc.L.Lock()
	defer stimLogger.wqc.L.Unlock()
	stimLogger.logQueue = append(stimLogger.logQueue, stimLogger.formatString(ll, level, args...))
	stimLogger.wqc.Broadcast()
}

func (stimLogger *FullStimLogger) formatString(ll Level, level string, args ...interface{}) *logMessage {
	var msg string
	switch args[0].(type) {
	case string:
		msg = args[0].(string)
	default:
		msg = fmt.Sprintf("%v", args[0])
	}
	subs := strings.Split(msg, subSTR)
	var sb strings.Builder
	if stimLogger.logTime {
		sb.WriteString(time.Now().Format(dateFMT))
		sb.WriteString("\t")
	}
	if stimLogger.logLevel {
		sb.WriteString(level)
		sb.WriteString("\t")
	}
	for i, v := range subs {
		v = strings.Replace(strings.Replace(v, "{{", "{", -1), "}}", "}", -1)
		sb.WriteString(v)
		if i < len(args)-1 {
			sb.WriteString(fmt.Sprintf("%v", args[i+1]))
		}
	}
	sb.WriteString("\n")
	return &logMessage{msg: sb.String(), logLevel: ll}
}

func (stimLogger *FullStimLogger) Flush() {
	for {
		stimLogger.wqc.L.Lock()
		l := len(stimLogger.logQueue)
		if l == 0 {
			stimLogger.wqc.L.Unlock()
			return
		} else {
			time.Sleep(time.Millisecond)
		}
		stimLogger.wqc.L.Unlock()
	}
}

func (stimLogger *FullStimLogger) ForceFlush(ff bool) {
	stimLogger.forceFlush = ff
	if ff {
		stimLogger.Flush()
	}
}

func (stimLogger *FullStimLogger) SetDateFormat(nf string) {
	stimLogger.dateFMT = nf
}

// SetLogger takes a structured logger to interface with.
// After the logger is setup it will be available across your packages
// If SetLogger is not used Debug will not create output
func (stimLogger *FullStimLogger) SetLogger(givenLogger Logger) {
	stimLogger.setLogger = givenLogger
}

// SetLevel sets the StimLogger log level.
func (stimLogger *FullStimLogger) SetLevel(level Level) {
	stimLogger.currentLevel = level
	hl := level
	for kv := range stimLogger.logfiles.Iter() {
		lgr := kv.Value.(*logFile)
		if lgr.logLevel > hl {
			hl = lgr.logLevel
		}
	}
	stimLogger.highestLevel = hl
}

// Debug logs a message at level Debug on the standard logger.
func (stimLogger *FullStimLogger) Debug(message ...interface{}) {
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
func (stimLogger *FullStimLogger) Verbose(message ...interface{}) {
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
func (stimLogger *FullStimLogger) Warn(message ...interface{}) {
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

// Trace logs a message at level Warn on the standard logger.
func (stimLogger *FullStimLogger) Trace(message ...interface{}) {
	if stimLogger.highestLevel >= TraceLevel {
		if stimLogger.setLogger == nil {
			if stimLogger.forceFlush {
				stimLogger.writeLogs(stimLogger.formatString(TraceLevel, traceMsg, message...))
			} else {
				stimLogger.formatAndLog(TraceLevel, traceMsg, message...)
			}
		} else {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

// Warn logs a message at level Warn on the standard logger.
func (stimLogger *FullStimLogger) Info(message ...interface{}) {
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
func (stimLogger *FullStimLogger) Fatal(message ...interface{}) {
	if stimLogger.highestLevel >= FatalLevel {
		if stimLogger.setLogger == nil {
			stimLogger.writeLogs(stimLogger.formatString(FatalLevel, fatalMsg, message...))
		} else {
			stimLogger.setLogger.Fatal(message...)
		}
		os.Exit(5)
	}
}

type StimPrefixLogger struct {
	stimLogger StimLogger
	prefix     string
}

func (spl *StimPrefixLogger) prefixLog(ll Level, i ...interface{}) []interface{} {
	if spl.GetLogLevel() >= ll {
		s := fmt.Sprintf("%v", i[0])
		var sb strings.Builder
		sb.WriteString(spl.prefix)
		sb.WriteString(":")
		sb.WriteString(s)
		i[0] = sb.String()
	}
	return i
}
func (spl *StimPrefixLogger) Trace(i ...interface{}) {
	spl.stimLogger.Trace(spl.prefixLog(TraceLevel, i...)...)
}
func (spl *StimPrefixLogger) Debug(i ...interface{}) {
	spl.stimLogger.Debug(spl.prefixLog(DebugLevel, i...)...)
}
func (spl *StimPrefixLogger) Verbose(i ...interface{}) {
	spl.stimLogger.Verbose(spl.prefixLog(VerboseLevel, i...)...)
}
func (spl *StimPrefixLogger) Info(i ...interface{}) {
	spl.stimLogger.Info(spl.prefixLog(InfoLevel, i...)...)
}
func (spl *StimPrefixLogger) Warn(i ...interface{}) {
	spl.stimLogger.Warn(spl.prefixLog(WarnLevel, i...)...)
}
func (spl *StimPrefixLogger) Fatal(i ...interface{}) {
	spl.stimLogger.Fatal(spl.prefixLog(FatalLevel, i...)...)
}
func (spl *StimPrefixLogger) GetLogLevel() Level { return spl.stimLogger.GetLogLevel() }
