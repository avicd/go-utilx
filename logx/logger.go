package logx

import (
	"github.com/avicd/go-utilx/conv"
	"time"
)

type Level = int

const (
	ALL Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

type CKind = uint

const (
	CLabel CKind = iota
	CMsg
	CLabelMsg
	CAll
)

var Labels = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

type Logger interface {
	SetLevel(level Level)
	GetLevel() Level
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

type Appender interface {
	Write(level Level, args ...any)
}

var logger Logger

func init() {
	logger = LoggerOf(&ConsoleAppender{})
	logger.SetLevel(ALL)
}

func Default() Logger {
	return logger
}

func SetDefault(nlog Logger) {
	if nlog != nil {
		logger = nlog
	}
}

func LoggerOf(appender Appender) Logger {
	return &Proxy{
		Appender: appender,
	}
}

func timeNow(layout string) string {
	if layout == "" {
		layout = conv.DateTimeMirco
	}
	return time.Now().Format(layout) + " "
}

func Debug(args ...any) {
	Default().Debug(args...)
}

func Debugf(format string, args ...any) {
	Default().Debugf(format, args...)
}

func Info(args ...any) {
	Default().Info(args...)
}
func Infof(format string, args ...any) {
	Default().Infof(format, args...)
}

func Warn(args ...any) {
	Default().Warn(args...)
}

func Warnf(format string, args ...any) {
	Default().Warnf(format, args...)
}

func Error(args ...any) {
	Default().Error(args...)
}

func Errorf(format string, args ...any) {
	Default().Errorf(format, args...)
}

func Fatal(args ...any) {
	Default().Fatal(args...)
}

func Fatalf(format string, args ...any) {
	Default().Fatalf(format, args...)
}
