package conv

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	Date              = "2006-01-02"
	DateTime          = "2006-01-02 15:04:05"
	DateTimeMirco     = "2006-01-02 15:04:05.000000"
	DateTimeMircoZone = "2006-01-02 15:04:05.000000Z07:00"
	DateTimeMilli     = "2006-01-02 15:04:05.000"
	DateTimeMilliZone = "2006-01-02 15:04:05.000Z07:00"
)

func CamelCase(text string) string {
	return regexp.MustCompile("[:\\-_]+[a-zA-z]").
		ReplaceAllStringFunc(text, func(s string) string {
			return strings.ToUpper(s[len(s)-1:])
		})
}

func BigCamelCase(text string) string {
	return FirstUpper(CamelCase(text))
}

func FirstUpper(text string) string {
	if len(text) < 1 {
		return text
	}
	return strings.ToUpper(text[0:1]) + text[1:]
}

func UnderLineCase(text string) string {
	return strings.ToLower(regexp.MustCompile("[a-z][A-Z]").
		ReplaceAllStringFunc(text, func(s string) string {
			return s[:1] + "_" + s[1:]
		}))
}

func Append(dest string, src any, delimiter ...string) string {
	if dest == "" {
		return fmt.Sprint(src)
	}
	var sep string
	if len(delimiter) > 0 {
		sep = delimiter[0]
	}
	return fmt.Sprint(dest, sep, src)
}

func StrToArr(text string, delimiter string) []string {
	if text == "" {
		return nil
	}
	return strings.Split(text, delimiter)
}

func IsDigit(text string) bool {
	return strings.IndexFunc(text, func(r rune) bool {
		return !unicode.IsDigit(r)
	}) < 0
}
