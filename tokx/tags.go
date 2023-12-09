package tokx

import (
	"io"
	"strings"
)

func ParseTags(text string) map[string][]string {
	rd := strings.NewReader(strings.TrimSpace(text))
	info := map[string][]string{}
	buffer := &strings.Builder{}
	var key string
	var list []string
	for {
		ch, err := rd.ReadByte()
		if err == io.EOF {
			break
		}
		switch ch {
		case ':':
			if len(list) > 0 {
				info[key] = list
				list = nil
			}
			key = buffer.String()
			buffer.Reset()
			continue
		case ',', ';':
			list = append(list, buffer.String())
			buffer.Reset()
			continue
		}
		buffer.WriteByte(ch)
	}
	if buffer.Len() > 0 {
		list = append(list, buffer.String())
	}
	if _, exist := info[key]; !exist {
		info[key] = list
	}
	return info
}
