package encoder

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"testing"

	. "github.com/cxcn/dtool/utils"
)

// 处理原始单字拼音表
func TestGenCharPyText(t *testing.T) {
	f, err := os.Open("own/src_char_pinyin.txt")
	if err != nil {
		log.Panic(err)
	}
	rd, _ := DecodeIO(f)

	type opys struct {
		order int
		codes []string
	}
	charMap := make(map[rune]*opys)
	var buf bytes.Buffer
	scan := bufio.NewScanner(rd)
	order := 0
	for ; scan.Scan(); order++ {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		char := []rune(entry[0])[0]
		if _, ok := charMap[char]; !ok {
			charMap[char] = &opys{order, entry[1:]}
			continue
		}
		charMap[char].codes = append(charMap[char].codes, entry[1:]...)
	}

	type owpys struct {
		word  rune
		order int
		codes []string
	}
	cmSli := make([]owpys, 0, len(charMap))
	for k, v := range charMap {
		cmSli = append(cmSli, owpys{k, v.order, RmRepeat(v.codes)})
	}
	sort.Slice(cmSli, func(i, j int) bool {
		return cmSli[i].order < cmSli[j].order
	})
	for _, v := range cmSli {
		buf.WriteRune(v.word)
		for _, vv := range v.codes {
			buf.WriteByte('\t')
			// buf.WriteByte('\'')
			buf.WriteString(vv)
		}
		buf.WriteString(LineBreak)
	}

	ioutil.WriteFile("assets/char_pinyin.txt", buf.Bytes(), 0777)
}
