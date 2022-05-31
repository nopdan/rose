package encoder

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"
)

//go:embed assets/word_pinyin.txt
var word_pinyin []byte

var WordPyMap = genWordPyMap(word_pinyin)

func genWordPyMap(data []byte) map[string][]string {
	ret := make(map[string][]string)
	rd := bytes.NewReader(data)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		ret[entry[0]] = strings.Split(entry[1], "'")
	}
	return ret
}
