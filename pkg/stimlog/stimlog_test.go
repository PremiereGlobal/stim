package stimlog

import (
	"io/ioutil"
	"os"
	"testing"

	"gotest.tools/assert"
)

var LOGLEVELS = []Level{WarnLevel, InfoLevel, VerboseLevel, DebugLevel, TraceLevel}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestSingleton(t *testing.T) {
	sa := make([]StimLogger, 0)
	sc := make(chan StimLogger)
	for i := 0; i < 100; i++ {
		go func() {
			sc <- GetLogger()
		}()
	}
	for i := 0; i < 100; i++ {
		sa = append(sa, <-sc)
	}
	for _, sl := range sa {
		assert.Equal(t, GetLogger(), sl, "Loggers are not Equal!")
	}
}

func TestSimpleLog(t *testing.T) {
	resetLogger()
	sl := GetLogger()
	sl.EnableTimeLogging(false)
	sl.EnableLevelLogging(false)
	sl.RemoveLogFile("STDOUT")
	for _, LL := range LOGLEVELS {
		tmpfile, err := ioutil.TempFile("", ".TESTLOG")
		defer os.Remove(tmpfile.Name())
		check(err)
		sl.SetLevel(LL)
		sl.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
		sl.Warn("Warn {} {}", "Message", LL)
		sl.Info("Info {} {}", "Message", LL)
		sl.Verbose("Verbose {} {}", "Message", LL)
		sl.Debug("Debug {} {}", "Message", LL)
		sl.Trace("Trace {} {}", "Message", LL)
		data, err := ioutil.ReadFile(tmpfile.Name())
		check(err)
		if LL == WarnLevel {
			assert.Equal(t, "Warn Message 20\n", string(data), "Loggers are not Equal!")
		} else if LL == InfoLevel {
			assert.Equal(t, "Warn Message 30\nInfo Message 30\n", string(data), "Loggers are not Equal!")
		} else if LL == VerboseLevel {
			assert.Equal(t, "Warn Message 40\nInfo Message 40\nVerbose Message 40\n", string(data), "Loggers are not Equal!")
		} else if LL == DebugLevel {
			assert.Equal(t, "Warn Message 50\nInfo Message 50\nVerbose Message 50\nDebug Message 50\n", string(data), "Loggers are not Equal!")
		} else {
			assert.Equal(t, "Warn Message 60\nInfo Message 60\nVerbose Message 60\nDebug Message 60\nTrace Message 60\n", string(data), "Loggers are not Equal!")
		}
	}
}

func TestSimpleLogNoFlush(t *testing.T) {
	resetLogger()
	sl := GetLogger()
	sl.EnableTimeLogging(false)
	sl.EnableLevelLogging(false)
	sl.RemoveLogFile("STDOUT")
	sl.ForceFlush(false)
	for _, LL := range LOGLEVELS {
		tmpfile, err := ioutil.TempFile("", "TESTLOG")
		// defer os.Remove(tmpfile.Name())
		check(err)
		sl.SetLevel(LL)
		sl.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
		sl.Warn("Warn {} {}", "Message", LL)
		sl.Info("Info {} {}", "Message", LL)
		sl.Verbose("Verbose {} {}", "Message", LL)
		sl.Debug("Debug {} {}", "Message", LL)
		sl.Trace("Trace {} {}", "Message", LL)
		sl.ForceFlush(true)
		data, err := ioutil.ReadFile(tmpfile.Name())
		check(err)
		if LL == WarnLevel {
			assert.Equal(t, "Warn Message 20\n", string(data), "Loggers are not Equal!")
		} else if LL == InfoLevel {
			assert.Equal(t, "Warn Message 30\nInfo Message 30\n", string(data), "Loggers are not Equal!")
		} else if LL == VerboseLevel {
			assert.Equal(t, "Warn Message 40\nInfo Message 40\nVerbose Message 40\n", string(data), "Loggers are not Equal!")
		} else if LL == DebugLevel {
			assert.Equal(t, "Warn Message 50\nInfo Message 50\nVerbose Message 50\nDebug Message 50\n", string(data), "Loggers are not Equal!")
		} else {
			assert.Equal(t, "Warn Message 60\nInfo Message 60\nVerbose Message 60\nDebug Message 60\nTrace Message 60\n", string(data), "Loggers are not Equal!")
		}
	}
}

func TestSimpleLogPrefix(t *testing.T) {
	resetLogger()
	sl := GetLoggerWithPrefix("PREFIX")
	sl.EnableTimeLogging(false)
	sl.EnableLevelLogging(false)
	sl.RemoveLogFile("STDOUT")
	for _, LL := range LOGLEVELS {
		tmpfile, err := ioutil.TempFile("", "TESTLOG")
		defer os.Remove(tmpfile.Name())
		check(err)
		sl.SetLevel(LL)
		sl.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
		sl.Warn("Warn {} {}", "Message", LL)
		sl.Info("Info {} {}", "Message", LL)
		sl.Verbose("Verbose {} {}", "Message", LL)
		sl.Debug("Debug {} {}", "Message", LL)
		sl.Trace("Trace {} {}", "Message", LL)
		data, err := ioutil.ReadFile(tmpfile.Name())
		check(err)
		if LL == WarnLevel {
			assert.Equal(t, "PREFIX:Warn Message 20\n", string(data), "Loggers are not Equal!")
		} else if LL == InfoLevel {
			assert.Equal(t, "PREFIX:Warn Message 30\nPREFIX:Info Message 30\n", string(data), "Loggers are not Equal!")
		} else if LL == VerboseLevel {
			assert.Equal(t, "PREFIX:Warn Message 40\nPREFIX:Info Message 40\nPREFIX:Verbose Message 40\n", string(data), "Loggers are not Equal!")
		} else if LL == DebugLevel {
			assert.Equal(t, "PREFIX:Warn Message 50\nPREFIX:Info Message 50\nPREFIX:Verbose Message 50\nPREFIX:Debug Message 50\n", string(data), "Loggers are not Equal!")
		} else {
			assert.Equal(t, "PREFIX:Warn Message 60\nPREFIX:Info Message 60\nPREFIX:Verbose Message 60\nPREFIX:Debug Message 60\nPREFIX:Trace Message 60\n", string(data), "Loggers are not Equal!")
		}
	}
}
