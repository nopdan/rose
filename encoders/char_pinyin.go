package encoder

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"
)

//go:embed assets/char_pinyin.txt
var char_pinyin []byte

var CharPyMap = genCharPyMap()

func genCharPyMap() map[rune][]string {
	ret := make(map[rune][]string)
	rd := bytes.NewReader(char_pinyin)
	scan := bufio.NewScanner(rd)

	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		ret[[]rune(entry[0])[0]] = entry[1:]
	}
	// ascii
	var a byte = 32
	for ; a < 127; a++ {
		ret[rune(a)] = []string{string(a)}
	}
	return ret
}
