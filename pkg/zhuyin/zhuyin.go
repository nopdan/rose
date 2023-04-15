package zhuyin

import (
	"strings"
)

func Get(word string) []string {
	chars := []rune(word)
	pinyin := make([]string, 0, 1)
	for i := 0; i < len(chars); {
		len, code, _ := m.Match(chars[i:])
		if len == 0 {
			pinyin = append(pinyin, string([]rune{'#', chars[i]}))
			i++
			continue
		}
		// fmt.Println(len, code)
		py := strings.Split(code, "'")
		pinyin = append(pinyin, py...)
		i += len
	}
	return pinyin
}

func GetOne(char rune) string {
	_, code, _ := m.Match([]rune{char})
	if code == "" {
		code = "#"
	}
	return code
}
