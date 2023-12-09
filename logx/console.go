package logx

import (
	"fmt"
	"log"
)

var Colors = map[Level]string{
	DEBUG: "34",
	INFO:  "32",
	WARN:  "33",
	ERROR: "31",
	FATAL: "35",
}

type ConsoleAppender struct {
	logger     *log.Logger
	Prefix     string
	TimeLayout string
	ColorOff   bool
	ColorKind  CKind
	Colors     map[Level]string
}

func (it *ConsoleAppender) getLogger() *log.Logger {
	if it.logger == nil {
		it.logger = log.New(log.Writer(), it.Prefix, 0)
	}
	return it.logger
}

func (it *ConsoleAppender) Write(level Level, args ...any) {
	label := Labels[level] + " "
	color := it.Colors[level]
	var dest []any
	if !it.ColorOff {
		if color == "" {
			color = Colors[level]
		}
		switch it.ColorKind {
		case CLabel:
			dest = append(dest, timeNow(it.TimeLayout))
			label = fmt.Sprintf("\033[%sm%s\033[0m", color, label)
			dest = append(dest, label)
			dest = append(dest, args...)
		case CMsg:
			dest = append(dest, timeNow(it.TimeLayout))
			dest = append(dest, label)
			dest = append(dest, fmt.Sprintf("\033[%sm", color))
			dest = append(dest, args...)
			dest = append(dest, "\033[0m")
		case CLabelMsg:
			dest = append(dest, timeNow(it.TimeLayout))
			dest = append(dest, fmt.Sprintf("\033[%sm", color))
			dest = append(dest, label)
			dest = append(dest, args...)
			dest = append(dest, "\033[0m")
		case CAll:
			dest = append(dest, fmt.Sprintf("\033[%sm", color))
			dest = append(dest, timeNow(it.TimeLayout))
			dest = append(dest, label)
			dest = append(dest, args...)
			dest = append(dest, "\033[0m")
		}
	}
	it.getLogger().Print(dest...)
}
