package conv

import (
	"math"
	"strings"
	"unicode"
)

type numInfo struct {
	negative bool
	base     int
	integer  string
	decimal  string
}

func isNumChar(r rune, base int) bool {
	if base <= 10 {
		return r >= '0' && r < '0'+rune(base)
	} else {
		return unicode.IsNumber(r) || r >= 'a' && r < 'a'+rune(base-10)
	}
}

func parseNum(num string) *numInfo {
	buf := strings.ToLower(num)
	info := &numInfo{}
	var cut int
	if cut = strings.Index(buf, "+"); cut > -1 {
		buf = buf[cut+1:]
	}
	if cut = strings.Index(buf, "-"); cut > -1 {
		buf = buf[cut+1:]
		info.negative = true
	}
	base := 10
	tokens := map[string]int{"0b": 2, "0o": 8, "0x": 16}
	for tk, bs := range tokens {
		if cut = strings.Index(buf, tk); cut > -1 {
			base = bs
			buf = buf[cut+len(tk):]
			break
		}
	}
	info.base = base
	valid := func(r rune) bool {
		return isNumChar(r, base)
	}
	invalid := func(r rune) bool {
		return !isNumChar(r, base)
	}
	if cut = strings.IndexFunc(buf, valid); cut > -1 {
		buf = buf[cut:]
	}
	if cut = strings.IndexFunc(buf, invalid); cut > -1 {
		info.integer = buf[:cut]
		buf = buf[cut:]
		if buf[0] == '.' {
			buf = buf[1:]
			if cut = strings.IndexFunc(buf, valid); cut > -1 {
				if cut = strings.IndexFunc(buf, invalid); cut > -1 {
					info.decimal = buf[:cut]
				} else {
					info.decimal = buf
				}
			}
		}
	} else {
		info.integer = buf
	}
	return info
}

func parseDigit(r rune) int32 {
	if r >= 'a' {
		return r - 'a' + 10
	} else if r <= '9' && r >= '0' {
		return r - '0'
	}
	return 0
}

func ParseInt(input string) int64 {
	info := parseNum(input)
	high := len(info.integer)
	var sum int64 = 0
	var level int64 = 1
	for _, ch := range info.integer {
		high--
		level = int64(math.Pow(float64(info.base), float64(high)))
		sum += int64(parseDigit(ch)) * level
	}
	if info.negative {
		sum = -sum
	}
	return sum
}

func ParseFloat(input string) float64 {
	info := parseNum(input)
	var sum = 0.0
	high := len(info.integer)
	plain := info.integer + info.decimal
	for _, ch := range plain {
		high--
		level := math.Pow(float64(info.base), float64(high))
		sum += float64(parseDigit(ch)) * level
	}
	if info.negative {
		sum = -sum
	}
	return sum
}
