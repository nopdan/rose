package zhuyin

import (
	"bufio"
	"bytes"
	"strings"

	_ "embed"

	"github.com/flowerime/gosmq/pkg/matcher"
)

//go:embed assets/char_pinyin.txt
var char_pinyin []byte

// 词库来源：
// 现代汉语常用词表（2008年李行健课题组）.txt
// 深蓝 WordPinyin.txt
//
//go:embed assets/word_pinyin.txt
var word_pinyin []byte

var m matcher.Matcher

func init() {
	m = matcher.NewStableTrie()
	// 词组
	rd := bytes.NewReader(word_pinyin)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) >= 2 {
			m.Insert(entry[0], entry[1], 1)
		}
	}
	// 单字
	rd.Reset(char_pinyin)
	scan = bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) >= 2 {
			m.Insert(entry[0], entry[1], 1)
		}
	}
	// ascii
	var a byte = 32
	for ; a < 127; a++ {
		m.Insert(string(a), string(a), 1)
	}
}
