package stimlog

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

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
	slc := GetLoggerConfig()
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)
	slc.RemoveLogFile("STDOUT")
	sl := GetLogger()
	for _, LL := range LOGLEVELS {
		tmpfile, err := ioutil.TempFile("", ".TESTLOG")
		defer os.Remove(tmpfile.Name())
		check(err)
		slc.SetLevel(LL)
		slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
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
	slc := GetLoggerConfig()
	sl := GetLogger()
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)
	slc.RemoveLogFile("STDOUT")
	slc.ForceFlush(false)
	for _, LL := range LOGLEVELS {
		tmpfile, err := ioutil.TempFile("", "TESTLOG")
		// defer os.Remove(tmpfile.Name())
		check(err)
		slc.SetLevel(LL)
		slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
		sl.Warn("Warn {} {}", "Message", LL)
		sl.Info("Info {} {}", "Message", LL)
		sl.Verbose("Verbose {} {}", "Message", LL)
		sl.Debug("Debug {} {}", "Message", LL)
		sl.Trace("Trace {} {}", "Message", LL)
		slc.ForceFlush(true)
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
			assert.Equal(t, "Warn Message 60\nInfo Message 60\nVerbose Message 60\nDebug Message 60\nTrace Message 60\nFlush\n", string(data), "Loggers are not Equal!")
		}
	}
}

func TestLogPrefixFirst(t *testing.T) {
	resetLogger()
	sl := GetLoggerWithPrefix("PREFIX")
	slc := GetLoggerConfig()
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)
	slc.RemoveLogFile("STDOUT")
	slc.SetLevel(TraceLevel)
	tmpfile, err := ioutil.TempFile("", "TESTLOG")
	defer os.Remove(tmpfile.Name())
	check(err)
	slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
	sl.Warn("Warn {}", "Message")
	sl.Info("Info {}", "Message")
	sl.Verbose("Verbose {}", "Message")
	sl.Debug("Debug {}", "Message")
	sl.Trace("Trace {}", "Message")
	data, err := ioutil.ReadFile(tmpfile.Name())
	check(err)
	assert.Equal(t, "PREFIX:Warn Message\nPREFIX:Info Message\nPREFIX:Verbose Message\nPREFIX:Debug Message\nPREFIX:Trace Message\n", string(data), "Loggers are not Equal!")
}

func TestSimpleLogLevelNoFlush(t *testing.T) {
	resetLogger()
	slc := GetLoggerConfig()
	sl := GetLogger()
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(true)
	slc.RemoveLogFile("STDOUT")
	slc.ForceFlush(false)
	for _, LL := range LOGLEVELS {
		tmpfile, err := ioutil.TempFile("", "TESTLOG")
		// defer os.Remove(tmpfile.Name())
		check(err)
		slc.SetLevel(LL)
		slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
		sl.Warn("Warn {} {}", "Message", LL)
		sl.Info("Info {} {}", "Message", LL)
		sl.Verbose("Verbose {} {}", "Message", LL)
		sl.Debug("Debug {} {}", "Message", LL)
		sl.Trace("Trace {} {}", "Message", LL)
		slc.Flush()
		data, err := ioutil.ReadFile(tmpfile.Name())
		check(err)
		if LL == WarnLevel {
			assert.Equal(t, "[ WARN  ]\tWarn Message 20\n", string(data), "Loggers are not Equal!")
		} else if LL == InfoLevel {
			assert.Equal(t, "[ WARN  ]\tWarn Message 30\n[ INFO  ]\tInfo Message 30\n", string(data), "Loggers are not Equal!")
		} else if LL == VerboseLevel {
			assert.Equal(t, "[ WARN  ]\tWarn Message 40\n[ INFO  ]\tInfo Message 40\n[VERBOSE]\tVerbose Message 40\n", string(data), "Loggers are not Equal!")
		} else if LL == DebugLevel {
			assert.Equal(t, "[ WARN  ]\tWarn Message 50\n[ INFO  ]\tInfo Message 50\n[VERBOSE]\tVerbose Message 50\n[ DEBUG ]\tDebug Message 50\n", string(data), "Loggers are not Equal!")
		} else {
			assert.Equal(t, "[ WARN  ]\tWarn Message 60\n[ INFO  ]\tInfo Message 60\n[VERBOSE]\tVerbose Message 60\n[ DEBUG ]\tDebug Message 60\n[ TRACE ]\tTrace Message 60\n[ TRACE ]\tFlush\n", string(data), "Loggers are not Equal!")
		}
	}
}

func TestSimpleLogDiffLevelLoggers(t *testing.T) {
	resetLogger()
	slc := GetLoggerConfig()
	sl := GetLogger()
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)
	slc.RemoveLogFile("STDOUT")
	slc.ForceFlush(true)

	warnfile, err := ioutil.TempFile("", "TESTLOG")
	defer os.Remove(warnfile.Name())
	check(err)

	debugfile, err := ioutil.TempFile("", "TESTLOG")
	defer os.Remove(debugfile.Name())
	check(err)

	slc.AddLogFile(warnfile.Name(), WarnLevel)
	slc.AddLogFile(debugfile.Name(), DebugLevel)

	slc.SetLevel(WarnLevel)

	sl.Warn("{} WARN", WarnLevel)
	sl.Debug("{} DEBUG", DebugLevel)

	data, err := ioutil.ReadFile(warnfile.Name())
	check(err)
	assert.Equal(t, "20 WARN\n", string(data), "Loggers are not Equal!")

	data, err = ioutil.ReadFile(debugfile.Name())
	check(err)
	assert.Equal(t, "20 WARN\n50 DEBUG\n", string(data), "Loggers are not Equal!")

	slc.RemoveLogFile(warnfile.Name())

	sl.Warn("{} WARN", WarnLevel)
	sl.Debug("{} DEBUG", DebugLevel)

	data, err = ioutil.ReadFile(warnfile.Name())
	check(err)
	assert.Equal(t, "20 WARN\n", string(data), "Loggers are not Equal!")

	data, err = ioutil.ReadFile(debugfile.Name())
	check(err)
	assert.Equal(t, "20 WARN\n50 DEBUG\n20 WARN\n50 DEBUG\n", string(data), "Loggers are not Equal!")
}

func TestSimpleLogNonString(t *testing.T) {
	resetLogger()
	slc := GetLoggerConfig()
	sl := GetLogger()
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)
	slc.RemoveLogFile("STDOUT")
	slc.ForceFlush(true)
	slc.SetLevel(WarnLevel)

	tmpfile, err := ioutil.TempFile("", "TESTLOG")
	defer os.Remove(tmpfile.Name())
	check(err)
	slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
	sl.Warn(errors.New("Test "), WarnLevel)

	data, err := ioutil.ReadFile(tmpfile.Name())
	check(err)
	assert.Equal(t, "Test 20\n", string(data), "Loggers are not Equal!")
}

func TestSimpleLogSyncWait(t *testing.T) {
	resetLogger()
	slc := GetLoggerConfig()
	sl := GetLogger()
	sl2 := GetLoggerWithPrefix("")
	assert.Equal(t, sl, sl2, "Loggers are not Equal!")
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)
	slc.RemoveLogFile("STDOUT")
	slc.ForceFlush(false)
	slc.SetLevel(WarnLevel)

	tmpfile, err := ioutil.TempFile("", "TESTLOG")
	defer os.Remove(tmpfile.Name())
	check(err)
	slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
	sl.Warn("Warn {} {}", "Message", WarnLevel)
	time.Sleep(time.Millisecond * 200)
	data, err := ioutil.ReadFile(tmpfile.Name())
	check(err)
	assert.Equal(t, "Warn Message 20\n", string(data), "Loggers are not Equal!")
}

func TestSimpleLogPrefix(t *testing.T) {
	resetLogger()
	slc := GetLoggerConfig()
	sl := GetLoggerWithPrefix("PREFIX")
	sl2 := GetLoggerWithPrefix("PREFIX")
	assert.Equal(t, sl, sl2, "Loggers are not Equal!")
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)
	slc.RemoveLogFile("STDOUT")
	for _, LL := range LOGLEVELS {
		tmpfile, err := ioutil.TempFile("", "TESTLOG")
		defer os.Remove(tmpfile.Name())
		check(err)
		slc.SetLevel(LL)
		slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
		sl.Warn("Warn {} {}", "Message", LL)
		sl2.Info("Info {} {}", "Message", LL)
		sl.Verbose("Verbose {} {}", "Message", LL)
		sl2.Debug("Debug {} {}", "Message", LL)
		sl.Trace("Trace {} {}", "Message", LL)
		slc.ForceFlush(true)
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
			// fmt.Println(string(data))
			assert.Equal(t, "PREFIX:Warn Message 60\nPREFIX:Info Message 60\nPREFIX:Verbose Message 60\nPREFIX:Debug Message 60\nPREFIX:Trace Message 60\nFlush\n", string(data), "Loggers are not Equal!")
		}
	}
}

func BenchmarkNoSyncNoFormat(b *testing.B) {
	b.ReportAllocs()
	resetLogger()
	sl := GetLogger()
	slc := GetLoggerConfig()
	slc.RemoveLogFile("STDOUT")
	// tmpfile, err := ioutil.TempFile("", "TESTLOG")
	// defer os.Remove(tmpfile.Name())
	// check(err)
	// slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
	slc.ForceFlush(false)
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func() {
				sl.Info("TEST")
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkNoSyncNoFormatDateLevel(b *testing.B) {
	b.ReportAllocs()
	resetLogger()
	sl := GetLogger()
	slc := GetLoggerConfig()
	slc.RemoveLogFile("STDOUT")
	// tmpfile, err := ioutil.TempFile("", "TESTLOG")
	// defer os.Remove(tmpfile.Name())
	// check(err)
	// slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
	slc.ForceFlush(false)
	slc.EnableTimeLogging(true)
	slc.EnableLevelLogging(true)

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func() {
				sl.Info("TEST")
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkNoSyncSimpleFormat(b *testing.B) {
	b.ReportAllocs()
	resetLogger()
	sl := GetLogger()
	slc := GetLoggerConfig()
	slc.RemoveLogFile("STDOUT")
	// tmpfile, err := ioutil.TempFile("", "TESTLOG")
	// defer os.Remove(tmpfile.Name())
	// check(err)
	// slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
	slc.ForceFlush(false)
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func() {
				sl.Info("TEST {} {} {}", 1, 2, 3)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSyncNoFormat(b *testing.B) {
	b.ReportAllocs()
	resetLogger()
	sl := GetLogger()
	slc := GetLoggerConfig()
	slc.RemoveLogFile("STDOUT")
	// tmpfile, err := ioutil.TempFile("", "TESTLOG")
	// defer os.Remove(tmpfile.Name())
	// check(err)
	// slc.AddLogFile(tmpfile.Name(), sl.GetLogLevel())
	slc.ForceFlush(true)
	slc.EnableTimeLogging(false)
	slc.EnableLevelLogging(false)
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func() {
				sl.Info("TEST")
				wg.Done()
			}()
		}
		wg.Wait()
	}

}

//TODO: add date/time test
//TODO: add file log error test
