package simlog

// from beego

import (
	"fmt"
	"log"
)

//--------------------
// LOG LEVEL
//--------------------

// Log levels to control the logging output.
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

var level = LevelTrace //日志等级，传参设定

// LogLevel returns the global log level and can be used in
// own implementations of the logger interface.
func LogLevel() int {
	return level
}

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLogLevel(l int) {
	level = l
}

// logger references the used application logger.
//var NettaoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var (
	SimLogger *log.Logger
	logSplitFunc func() = func(){}
)

// SetLogger sets a new logger.
func SetLogger(l *log.Logger) {
	SimLogger = l
}

func SetSplit(maxsize, maxindex int) {
	setLogSplit(maxsize, maxindex)
}

func Init(path string, maxday, loglevel int) {
	logInit(path, maxday, loglevel)
}

func Close() {
	logClose()
}

func Trace(v ...interface{}) {
	if level < LevelTrace {
		logSplitFunc()
		SimLogger.Output(3, fmt.Sprint(v...))
	}
}

func Debug(v ...interface{}) {
	if level < LevelDebug {
		logSplitFunc()
		SimLogger.Output(3, fmt.Sprint(v...))
	}
}

func Info(v ...interface{}) {
	if level < LevelInfo {
		logSplitFunc()
		SimLogger.Output(3, fmt.Sprint(v...))
	}
}

func Warn(v ...interface{}) {
	if level < LevelWarning {
		logSplitFunc()
		SimLogger.Output(3, fmt.Sprint(v...))
	}
}

func Error(v ...interface{}) {
	if level < LevelError {
		logSplitFunc()
		SimLogger.Output(3, fmt.Sprint(v...))
	}
}

func Critical(v ...interface{}) {
	if level < LevelCritical {
		logSplitFunc()
		SimLogger.Output(3, fmt.Sprint(v...))
	}
}

// Trace logs a message at trace level.
func LogTrace(format string) {
	if level <= LevelTrace {
		SimLogger.Printf("[T]" + format)
	}
}

// Debug logs a message at debug level.
func LogDebug(format string, skips ...int) {
	if level <= LevelDebug {
		logPrintf("[D]"+format, skips...)
	}
}

// Info logs a message at info level.
func LogInfo(format string, skips ...int) {
	if level <= LevelInfo {
		logPrintf("[I]"+format, skips...)
	}
}

// Warning logs a message at warning level.
func LogWarn(format string, skips ...int) {
	if level <= LevelWarning {
		logPrintf("[W]"+format, skips...)
	}
}

// Error logs a message at error level.
func LogError(format string, skips ...int) {
	if level <= LevelError {
		logPrintf("[E]"+format, skips...)
	}
}

// Critical logs a message at critical level.
func LogCritical(format string, skips ...int) {
	if level <= LevelCritical {
		logPrintf("[C]"+format, skips...)
	}
}

func logPrintf(format string, v ...int) {
	logSplitFunc()
	skip := 3
	if len(v) > 0 {skip = v[0]}
	SimLogger.Output(skip, format)
}

func logSplit() {
	logfilelock.Lock()
	defer logfilelock.Unlock()
	logcnt++
	if logcnt > MAXLOGCNT {
		changelogindex(1)
		logcnt = 0
	}
}