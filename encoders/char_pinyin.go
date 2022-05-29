package encoder

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"
)

//go:embed assets/char_pinyin.txt
var char_pinyin []byte

var CharPyMap = getCharPyMap()

func getCharPyMap() map[rune][]string {
	ret := make(map[rune][]string)
	rd := bytes.NewReader(char_pinyin)
	scan := bufio.NewScanner(rd)

	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		ret[[]rune(entry[0])[0]] = entry[1:]
	}
	return ret
}
