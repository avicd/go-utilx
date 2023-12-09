package logx

import (
	"fmt"
	"os"
	"runtime/debug"
)

type Proxy struct {
	Level    Level
	Appender Appender
}

func (it *Proxy) GetLevel() Level {
	return it.Level
}

func (it *Proxy) SetLevel(level Level) {
	it.Level = level
}

func (it *Proxy) Debug(args ...any) {
	if DEBUG >= it.Level {
		it.Appender.Write(DEBUG, args...)
	}
}

func (it *Proxy) Debugf(format string, args ...any) {
	if DEBUG >= it.Level {
		it.Appender.Write(DEBUG, fmt.Sprintf(format, args...))
	}
}

func (it *Proxy) Info(args ...any) {
	if INFO >= it.Level {
		it.Appender.Write(INFO, args...)
	}
}

func (it *Proxy) Infof(format string, args ...any) {
	if INFO >= it.Level {
		it.Appender.Write(INFO, fmt.Sprintf(format, args...))
	}
}

func (it *Proxy) Warn(args ...any) {
	if WARN >= it.Level {
		it.Appender.Write(WARN, args...)
	}
}

func (it *Proxy) Warnf(format string, args ...any) {
	if WARN >= it.Level {
		it.Appender.Write(WARN, fmt.Sprintf(format, args...))
	}
}

func (it *Proxy) Error(args ...any) {
	if ERROR >= it.Level {
		args = append(args, "\n"+string(debug.Stack()))
		it.Appender.Write(ERROR, args...)
	}
}

func (it *Proxy) Errorf(format string, args ...any) {
	if ERROR >= it.Level {
		var dest []any
		dest = append(dest, fmt.Sprintf(format, args...))
		dest = append(dest, "\n"+string(debug.Stack()))
		it.Appender.Write(ERROR, dest...)
	}
}

func (it *Proxy) Fatal(args ...any) {
	if FATAL >= it.Level {
		args = append(args, "\n"+string(debug.Stack()))
		it.Appender.Write(FATAL, args...)
	}
	os.Exit(1)
}

func (it *Proxy) Fatalf(format string, args ...any) {
	if FATAL >= it.Level {
		var dest []any
		dest = append(dest, fmt.Sprintf(format, args...))
		dest = append(dest, "\n"+string(debug.Stack()))
		it.Appender.Write(FATAL, dest...)
	}
	os.Exit(1)
}
