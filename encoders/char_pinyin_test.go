package encoder

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	. "github.com/cxcn/dtool/utils"
)

// 处理原始单字拼音表
func TestGenCharPyText(t *testing.T) {
	f, err := os.Open("own/char_pinyin.txt")
	if err != nil {
		log.Panic(err)
	}
	rd, _ := DecodeIO(f)

	charMap := make(map[rune][]string)
	var buf bytes.Buffer
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		char := []rune(entry[0])[0]
		if _, ok := charMap[char]; !ok {
			charMap[char] = []string{}
		}
		charMap[char] = append(charMap[char], entry[1:]...)
		charMap[char] = RmRepeat(charMap[char])

		buf.WriteRune(char)
		for _, v := range charMap[char] {
			buf.WriteByte('\t')
			// buf.WriteByte('\'')
			buf.WriteString(v)
		}
		buf.WriteString(LineBreak)
	}
	ioutil.WriteFile("assets/char_pinyin.txt", buf.Bytes(), 0777)
}
