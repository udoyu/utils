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
type Level int
const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

var (
	levelName = []string{"[T]", "[D]", "[I]", "[W]", "[E]", "[C]"}
	level = LevelTrace //日志等级，传参设定
	SimLogger *log.Logger
	logSplitFunc func()
	Trace func(v ...interface{}) = func(v ...interface{}){PrintFunc(LevelTrace, v...)}
	Debug func(v ...interface{}) = func(v ...interface{}){PrintFunc(LevelDebug, v...)}
	Info func(v ...interface{}) = func(v ...interface{}){PrintFunc(LevelInfo, v...)}
	Warn func(v ...interface{}) = func(v ...interface{}){PrintFunc(LevelWarning, v...)}
	Error func(v ...interface{}) = func(v ...interface{}){PrintFunc(LevelError, v...)}
	Critical func(v ...interface{}) = func(v ...interface{}){PrintFunc(LevelCritical, v...)}
)

func (this Level) String() string {
	return levelName[this]
}

// LogLevel returns the global log level and can be used in
// own implementations of the logger interface.
func LogLevel() Level {
	return level
}

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLogLevel(l Level) {
	level = l
}

// logger references the used application logger.
//var NettaoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)


// SetLogger sets a new logger.
func SetLogger(l *log.Logger) {
	SimLogger = l
}

func SetSplit(maxsize, maxindex int) {
	setLogSplit(maxsize, maxindex)
	logSplitFunc = logSplit
}

func Init(path string, maxday int, loglevel Level) {
	logInit(path, maxday, loglevel)
	Trace = TraceFunc
	Debug = DebugFunc
	Info = InfoFunc
	Warn = WarnFunc
	Error = ErrorFunc
	Critical = CriticalFunc
}

func Close() {
	logClose()
}

func PrintFunc (loglevel Level, v ...interface{}) {
	if level <= loglevel {
		log.Println(loglevel.String() + fmt.Sprint(v...))
	}
}

func TraceFunc(v ...interface{}) {
	if level <= LevelTrace {
		logSplitFunc()
		SimLogger.Output(2, LevelTrace.String() + fmt.Sprint(v...))
	}
}

func DebugFunc(v ...interface{}) {
	if level <= LevelDebug {
		logSplitFunc()
		SimLogger.Output(2,LevelDebug.String() + fmt.Sprint(v...))
	}
}

func InfoFunc(v ...interface{}) {
	if level <= LevelInfo {
		logSplitFunc()
		SimLogger.Output(2, LevelInfo.String() + fmt.Sprint(v...))
	}
}

func WarnFunc(v ...interface{}) {
	if level <= LevelWarning {
		logSplitFunc()
		SimLogger.Output(2, LevelWarning.String() + fmt.Sprint(v...))
	}
}

func ErrorFunc(v ...interface{}) {
	if level <= LevelError {
		logSplitFunc()
		SimLogger.Output(2, LevelError.String() + fmt.Sprint(v...))
	}
}

func CriticalFunc(v ...interface{}) {
	if level <= LevelCritical {
		logSplitFunc()
		SimLogger.Output(2, LevelCritical.String() + fmt.Sprint(v...))
	}
}

// Trace logs a message at trace level.
func LogTrace(format string) {
	if level <= LevelTrace {
		SimLogger.Printf(LevelTrace.String() + format)
	}
}

// Debug logs a message at debug level.
func LogDebug(format string, skips ...int) {
	if level <= LevelDebug {
		logPrintf(LevelDebug.String() + format, skips...)
	}
}

// Info logs a message at info level.
func LogInfo(format string, skips ...int) {
	if level <= LevelInfo {
		logPrintf(LevelInfo.String() + format, skips...)
	}
}

// Warning logs a message at warning level.
func LogWarn(format string, skips ...int) {
	if level <= LevelWarning {
		logPrintf(LevelWarning.String() + format, skips...)
	}
}

// Error logs a message at error level.
func LogError(format string, skips ...int) {
	if level <= LevelError {
		logPrintf(LevelError.String() + format, skips...)
	}
}

// Critical logs a message at critical level.
func LogCritical(format string, skips ...int) {
	if level <= LevelCritical {
		logPrintf(LevelCritical.String() + format, skips...)
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