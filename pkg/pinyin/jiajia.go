package pinyin

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/nopdan/rose/pkg/encoder"
)

type JiaJia struct{ Template }

func init() {
	FormatList = append(FormatList, NewJiaJia())
}
func NewJiaJia() *JiaJia {
	f := new(JiaJia)
	f.Name = "拼音加加"
	f.ID = "pyjj,jj"
	f.CanMarshal = true
	return f
}

func (f *JiaJia) Unmarshal(r *bytes.Reader) []*Entry {
	d := make([]*Entry, 0, r.Size()>>8)

	enc := encoder.NewPinyin()
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		entry := []rune(strings.TrimSpace(scan.Text()))
		// 注释
		if entry[0] == ';' {
			continue
		}
		wordRunes := make([]rune, 0, 1)
		pinyin := make([]string, 0, 1)
		for i := 0; i < len(entry); {
			char := entry[i]
			wordRunes = append(wordRunes, char)
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
			pys := enc.Encode(string(wordRunes))
			if len(pys) != 0 {
				pinyin = append(pinyin, pys[0])
			}
			i++
		}
		d = append(d, &Entry{string(wordRunes), pinyin, 1})
	}
	return d
}

func (JiaJia) Marshal(di []*Entry) []byte {
	var buf bytes.Buffer
	for _, v := range di {
		words := []rune(v.Word)
		if len(words) != len(v.Pinyin) {
			continue
		}
		for i := 0; i < len(words); i++ {
			buf.WriteString(string(words[i]))
			buf.WriteString(v.Pinyin[i])
		}
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}
