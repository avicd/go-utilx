package tokx

import (
	"strings"
)

type Pair struct {
	left  string
	right string
}

func NewPair(left string, right string) *Pair {
	return &Pair{left: left, right: right}
}

func (it *Pair) Match(text string) bool {
	if text == "" {
		return false
	}
	buf := text
	left := false
	expect := it.left
	for {
		si := strings.Index(buf, expect)
		if si < 0 {
			return false
		} else if si > 0 && buf[si-1] == '\\' {
			buf = buf[si+len(expect):]
			continue
		}
		if left {
			return true
		}
		left = !left
		expect = it.right
	}
}

func (it *Pair) Map(text string, rel func(s string) string) string {
	if text == "" {
		return text
	}
	buf := text
	left := false
	expect := it.left
	str := &strings.Builder{}
	key := &strings.Builder{}
	for {
		si := strings.Index(buf, expect)
		if si < 0 {
			str.WriteString(buf)
			break
		} else if si > 0 && buf[si-1] == '\\' {
			if left {
				key.WriteString(buf[:si+len(expect)])
			} else {
				str.WriteString(buf[:si+len(expect)])
			}
			buf = buf[si+len(expect):]
			continue
		}
		next := buf[si+len(expect):]
		if left {
			key.WriteString(buf[:si])
			str.WriteString(rel(key.String()))
			key.Reset()
			expect = it.left
		} else {
			str.WriteString(buf[:si])
			expect = it.right
		}
		buf = next
		left = !left
	}
	return str.String()
}
