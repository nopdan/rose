package rose

import (
	"bufio"
	"bytes"
	"strings"

	util "github.com/flowerime/goutil"
	"github.com/flowerime/rose/pkg/zhuyin"
)

type JiaJia struct{ Dict }

func NewJiaJia() *JiaJia {
	d := new(JiaJia)
	d.IsPinyin = true
	d.IsBinary = false
	d.Name = "拼音加加.txt"
	return d
}

func (d *JiaJia) GetDict() *Dict {
	return &d.Dict
}

func (d *JiaJia) Parse() {
	pyt := make(PyTable, 0, 0xff)
	scan := bufio.NewScanner(d.rd)
	for scan.Scan() {
		entry := []rune(strings.TrimSpace(scan.Text()))
		// 注释
		if entry[0] == ';' {
			continue
		}
		word := make([]rune, 0, 1)
		pinyin := make([]string, 0, 1)

		for i := 0; i < len(entry); {
			char := entry[i]
			word = append(word, char)
			// 下一个是英文（拼音）
			if i+1 != len(entry) && entry[i+1] < 128 {
				j := 1 // 已匹配的字母数
				for i+j+1 < len(entry) && entry[i+j+1] < 128 {
					j++
				}
				pinyin = append(pinyin, string(entry[i+1:i+j+1]))
				i += j + 1
				continue
			}
			// 读到汉字
			code := zhuyin.GetOne(char)
			pinyin = append(pinyin, code)
			i++
		}
		pyt = append(pyt, &PinyinEntry{string(word), pinyin, 1})
	}
	d.pyt = pyt
}

func (d *JiaJia) GenFrom(src *Dict) []byte {
	var buf bytes.Buffer
	for _, v := range src.pyt {
		words := []rune(v.Word)
		if len(words) != len(v.Pinyin) {
			continue
		}
		for i := 0; i < len(words); i++ {
			buf.WriteString(string(words[i]))
			buf.WriteString(v.Pinyin[i])
		}
		buf.WriteString(util.LineBreak)
	}
	return buf.Bytes()
}
