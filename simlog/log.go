package simlog

import (
	"log"
)

//--------------------
// LOG LEVEL
//--------------------

// Log levels to control the logging output.

type LogInterface interface {
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	SetLevel(Level)
	GetLevel() Level
	Close()
}

var logLevelNames = []string{
	"[TRACE]", "[DEBUG]", "[INFO]", "[WARN]", "[ERROR]", "[CRITICAL]",
}

func SetLevelName(names []string) {
	levelName = names
}

type Level int

func (this Level) String() string {
	return levelName[this]
}

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

var (
	LogFlag              = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
	Logger  LogInterface = &LogHandler{
		level:       LevelTrace,
		MaxDataSize: 4096,
	}
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

func Trace(v ...interface{}) {
	Logger.Trace(v...)
}

func Debug(v ...interface{}) {
	Logger.Debug(v...)
}

func Info(v ...interface{}) {
	Logger.Info(v...)
}

func Warn(v ...interface{}) {
	Logger.Warn(v...)
}

func Error(v ...interface{}) {
	Logger.Error(v...)
}

func Critical(v ...interface{}) {
	Logger.Critical(v...)
}

// LogLevel returns the global log level and can be used in
// own implementations of the logger interface.
func GetLevel() Level {
	return Logger.GetLevel()
}

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l Level) {
	Logger.SetLevel(l)
}

func SetSplit(maxsize, maxindex int) {
	if logger, ok := Logger.(*LogHandler); ok {
		logger.SetLogSplit(maxsize, maxindex)
	}
}

func Init(path string, maxday int, loglevel Level) {
	if logger, ok := Logger.(*LogHandler); ok {
		logger.Init(path, maxday, loglevel)
	}
}

func Close() {
	Logger.Close()
}

var (
	levelName = []string{"[TRACE]", "[DEBUG]", "[INFO]", "[WARN]", "[ERROR]", "[CRITICAL]"}
)
