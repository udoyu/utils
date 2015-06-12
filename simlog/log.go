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
var SimLogger *log.Logger

// SetLogger sets a new logger.
func SetLogger(l *log.Logger) {
	SimLogger = l
}

//用来替换NettaoLogger.Printf

func logPrintf(format string, v ...int) {
	logfilelock.Lock()
	defer logfilelock.Unlock()
	logcnt++
	if logcnt > MAXLOGCNT {
		changelogindex(len(format) + 1)
		logcnt = 0
	}
	skip := 3
	if len(v) > 0 {
		skip = v[0]
	}
	SimLogger.Output(skip, fmt.Sprintf(format))
}

// Trace logs a message at trace level.
func Trace(format string) {
	if level <= LevelTrace {
		SimLogger.Printf("[T]" + format)
	}
}

// Debug logs a message at debug level.
func Debug(format string, skips ...int) {
	if level <= LevelDebug {
		logPrintf("[D]"+format, skips...)
	}
}

// Info logs a message at info level.
func Info(format string, skips ...int) {
	if level <= LevelInfo {
		logPrintf("[I]"+format, skips...)
	}
}

// Warning logs a message at warning level.
func Warn(format string, skips ...int) {
	if level <= LevelWarning {
		logPrintf("[W]"+format, skips...)
	}
}

// Error logs a message at error level.
func Error(format string, skips ...int) {
	if level <= LevelError {
		logPrintf("[E]"+format, skips...)
	}
}

// Critical logs a message at critical level.
func Critical(format string, skips ...int) {
	if level <= LevelCritical {
		logPrintf("[C]"+format, skips...)
	}
}
