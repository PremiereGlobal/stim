package stimlog

import (
	"fmt"
	"os"
	"strconv"
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
	ltime    time.Time
	msg      string
	args     []interface{}
	wg       *sync.WaitGroup
}

type fullStimLogger struct {
	setLogger    Logger
	currentLevel Level
	highestLevel Level
	dateFMT      string
	logfiles     hashmap.HashMap
	logQueue     chan *logMessage
	forceFlush   bool
	logLevel     bool
	logTime      bool
	// wqc          *sync.Cond
}

var logger *fullStimLogger

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
	prefixLogger = nil
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
			logger = &fullStimLogger{
				currentLevel: InfoLevel,
				highestLevel: InfoLevel,
				dateFMT:      dateFMT,
				logQueue:     make(chan *logMessage, 1000),
				logfiles:     hashmap.HashMap{},
				forceFlush:   true,
				logLevel:     true,
				logTime:      true,
			}
			logger.AddLogFile("STDOUT", defaultLevel)
			go logger.writeLogQueue()
		}
	}
	return logger
}

//GetLoggerWithPrefix gets a logger for logging in stim with a prefix.
func GetLoggerWithPrefix(prefix string) StimLogger {
	baseLogger := GetLogger()
	if prefix == "" {
		return baseLogger
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
	prefixLogger[prefix] = &stimPrefixLogger{stimLogger: logger, prefix: prefix}
	return prefixLogger[prefix]
}

//GetLogLevel gets the current highest set log level for stimlog
func (stimLogger *fullStimLogger) GetLogLevel() Level {
	return stimLogger.highestLevel
}

//EnableLevelLogging enables/disables logging of the level (WARN/DEBUG, etc)
func (stimLogger *fullStimLogger) EnableLevelLogging(b bool) {
	stimLogger.logLevel = b
}

//EnableTimeLogging enables/disables logging of the timestamp
func (stimLogger *fullStimLogger) EnableTimeLogging(b bool) {
	stimLogger.logTime = b
}

//RemoveLogFile removes logging of a file (can be STDOUT/STDERR too)
func (stimLogger *fullStimLogger) RemoveLogFile(file string) {
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

//AddLogFile adds logging of a file (can be STDOUT/STDERR too)
func (stimLogger *fullStimLogger) AddLogFile(file string, logLevel Level) error {
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

func (stimLogger *fullStimLogger) writeLogQueue() {
	syncNone := time.Duration(time.Hour * 100)
	syncSoon := time.Duration(time.Millisecond * 100)
	syncDelay := syncNone
	for {
		select {
		case lm := <-stimLogger.logQueue:
			stimLogger.writeLogs(lm.logLevel, lm.wg != nil, stimLogger.formatLogMessage(lm))
			if lm.wg != nil {
				lm.wg.Done()
				syncDelay = syncNone
			} else {
				if !stimLogger.forceFlush {
					syncDelay = syncSoon
				}
			}
		case <-time.After(syncDelay):
			for kv := range stimLogger.logfiles.Iter() {
				lgr := kv.Value.(*logFile)
				lgr.fp.Sync()
			}
			syncDelay = syncNone
		}
	}
}

func (stimLogger *fullStimLogger) formatLogMessage(lm *logMessage) string {
	var sb strings.Builder
	if stimLogger.logTime {
		sb.WriteString(lm.ltime.Format(dateFMT))
		sb.WriteString("\t")
	}
	if stimLogger.logLevel {
		if lm.logLevel == FatalLevel {
			sb.WriteString(fatalMsg)
		} else if lm.logLevel == WarnLevel {
			sb.WriteString(warnMsg)
		} else if lm.logLevel == InfoLevel {
			sb.WriteString(infoMsg)
		} else if lm.logLevel == VerboseLevel {
			sb.WriteString(verboseMsg)
		} else if lm.logLevel == DebugLevel {
			sb.WriteString(debugMsg)
		} else if lm.logLevel == TraceLevel {
			sb.WriteString(traceMsg)
		} else {
			sb.WriteString(strconv.FormatInt(int64(lm.logLevel), 10))
		}
		sb.WriteString("\t")
	}

	subs := strings.Split(lm.msg, subSTR)

	for i, v := range subs {
		v = strings.Replace(strings.Replace(v, "{{", "{", -1), "}}", "}", -1)
		sb.WriteString(v)
		if i < len(lm.args) {
			sb.WriteString(fmt.Sprintf("%v", lm.args[i]))
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

func (stimLogger *fullStimLogger) writeLogs(logLevel Level, sync bool, msg string) {
	for kv := range stimLogger.logfiles.Iter() {
		lgr := kv.Value.(*logFile)
		if lgr.logLevel >= logLevel || (lgr.logLevel == defaultLevel && stimLogger.currentLevel >= logLevel) {
			lgr.fp.WriteString(msg)
			if stimLogger.forceFlush || sync {
				lgr.fp.Sync()
			}
		}
	}
}

func (stimLogger *fullStimLogger) wrapMessage(ll Level, wg *sync.WaitGroup, args ...interface{}) *logMessage {
	st := time.Now()
	var msg string
	switch args[0].(type) {
	case string:
		msg = args[0].(string)
	default:
		msg = fmt.Sprintf("%v", args[0])
	}
	return &logMessage{
		ltime:    st,
		msg:      msg,
		args:     args[1:],
		logLevel: ll,
		wg:       wg,
	}
}

//Forces logger to write all current logs, will block till done
func (stimLogger *fullStimLogger) Flush() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	stimLogger.logQueue <- stimLogger.wrapMessage(TraceLevel, wg, "Flush")
	wg.Wait()
}

//ForceFlush enables/disables forcing sync on logfiles after every write
func (stimLogger *fullStimLogger) ForceFlush(ff bool) {
	stimLogger.forceFlush = ff
	if ff {
		stimLogger.Flush()
	}
}

//SetDateFormat allows you to set how the date/time is formated
func (stimLogger *fullStimLogger) SetDateFormat(nf string) {
	stimLogger.dateFMT = nf
}

// SetLogger takes a structured logger to interface with.
// After the logger is setup it will be available across your packages
// If SetLogger is not used Debug will not create output
func (stimLogger *fullStimLogger) SetLogger(givenLogger Logger) {
	stimLogger.setLogger = givenLogger
}

// SetLevel sets the StimLogger log level.
func (stimLogger *fullStimLogger) SetLevel(level Level) {
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
func (stimLogger *fullStimLogger) Debug(message ...interface{}) {
	if stimLogger.highestLevel >= DebugLevel {
		if stimLogger.setLogger == nil {
			if stimLogger.forceFlush {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				stimLogger.logQueue <- stimLogger.wrapMessage(DebugLevel, wg, message...)
				wg.Wait()
			} else {
				stimLogger.logQueue <- stimLogger.wrapMessage(DebugLevel, nil, message...)
			}
		} else {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

//Verbose logs a message at level Verbose on the standard logger.
func (stimLogger *fullStimLogger) Verbose(message ...interface{}) {
	if stimLogger.highestLevel >= VerboseLevel {
		if stimLogger.setLogger == nil {
			if stimLogger.forceFlush {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				stimLogger.logQueue <- stimLogger.wrapMessage(VerboseLevel, wg, message...)
				wg.Wait()
			} else {
				stimLogger.logQueue <- stimLogger.wrapMessage(VerboseLevel, nil, message...)
			}
		} else {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

// Warn logs a message at level Warn on the standard logger.
func (stimLogger *fullStimLogger) Warn(message ...interface{}) {
	if stimLogger.highestLevel >= WarnLevel {
		if stimLogger.setLogger == nil {
			if stimLogger.forceFlush {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				stimLogger.logQueue <- stimLogger.wrapMessage(WarnLevel, wg, message...)
				wg.Wait()
			} else {
				stimLogger.logQueue <- stimLogger.wrapMessage(WarnLevel, nil, message...)
			}
		} else {
			stimLogger.setLogger.Warn(message...)
		}
	}
}

// Trace logs a message at level Warn on the standard logger.
func (stimLogger *fullStimLogger) Trace(message ...interface{}) {
	if stimLogger.highestLevel >= TraceLevel {
		if stimLogger.setLogger == nil {
			if stimLogger.forceFlush {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				stimLogger.logQueue <- stimLogger.wrapMessage(TraceLevel, wg, message...)
				wg.Wait()
			} else {
				stimLogger.logQueue <- stimLogger.wrapMessage(TraceLevel, nil, message...)
			}
		} else {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

// Info logs a message at level Info on the standard logger.
func (stimLogger *fullStimLogger) Info(message ...interface{}) {
	if stimLogger.highestLevel >= InfoLevel {
		if stimLogger.setLogger == nil {

			if stimLogger.forceFlush {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				stimLogger.logQueue <- stimLogger.wrapMessage(InfoLevel, wg, message...)
				wg.Wait()
			} else {
				stimLogger.logQueue <- stimLogger.wrapMessage(InfoLevel, nil, message...)
			}

		} else {
			stimLogger.setLogger.Debug(message...)
		}
	}
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (stimLogger *fullStimLogger) Fatal(message ...interface{}) {
	if stimLogger.highestLevel >= FatalLevel {
		if stimLogger.setLogger == nil {
			wg := &sync.WaitGroup{}
			wg.Add(1)
			stimLogger.logQueue <- stimLogger.wrapMessage(InfoLevel, wg, message...)
			wg.Wait()
		} else {
			stimLogger.setLogger.Fatal(message...)
		}
		os.Exit(5)
	}
}

type stimPrefixLogger struct {
	stimLogger *fullStimLogger
	prefix     string
}

func (spl *stimPrefixLogger) prefixLog(i ...interface{}) []interface{} {
	s := fmt.Sprintf("%v", i[0])
	var sb strings.Builder
	sb.WriteString(spl.prefix)
	sb.WriteString(":")
	sb.WriteString(s)
	i[0] = sb.String()
	return i
}
func (spl *stimPrefixLogger) Trace(i ...interface{}) {
	if spl.stimLogger.GetLogLevel() >= TraceLevel {
		spl.stimLogger.Trace(spl.prefixLog(i...)...)
	}
}
func (spl *stimPrefixLogger) Debug(i ...interface{}) {
	if spl.stimLogger.GetLogLevel() >= DebugLevel {
		spl.stimLogger.Debug(spl.prefixLog(i...)...)
	}
}
func (spl *stimPrefixLogger) Verbose(i ...interface{}) {
	if spl.stimLogger.GetLogLevel() >= VerboseLevel {
		spl.stimLogger.Verbose(spl.prefixLog(i...)...)
	}
}
func (spl *stimPrefixLogger) Info(i ...interface{}) {
	if spl.stimLogger.GetLogLevel() >= InfoLevel {
		spl.stimLogger.Info(spl.prefixLog(i...)...)
	}
}
func (spl *stimPrefixLogger) Warn(i ...interface{}) {
	if spl.stimLogger.GetLogLevel() >= WarnLevel {
		spl.stimLogger.Warn(spl.prefixLog(i...)...)
	}
}
func (spl *stimPrefixLogger) Fatal(i ...interface{}) {
	if spl.stimLogger.GetLogLevel() >= FatalLevel {
		spl.stimLogger.Fatal(spl.prefixLog(i...)...)
	}
}
func (spl *stimPrefixLogger) GetLogLevel() Level { return spl.stimLogger.GetLogLevel() }
