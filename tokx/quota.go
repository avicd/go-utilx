package tokx

import "strings"

func DoubleQuota(text string) string {
	buf := text
	str := &strings.Builder{}
	var left string
	for {
		var si int
		if left == "" {
			si = strings.IndexFunc(buf, func(r rune) bool {
				return r == '"' || r == '\''
			})
		} else {
			si = strings.Index(buf, left)
		}
		if si > -1 {
			if left == "" {
				left = buf[si : si+1]
			} else {
				left = ""
			}
			str.WriteString(buf[:si])
			str.WriteString("\"")
			buf = buf[si+1:]
		} else {
			str.WriteString(buf)
			break
		}
	}
	return str.String()
}
