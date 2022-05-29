package pinyin

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	encoder "github.com/cxcn/dtool/encoders"
	. "github.com/cxcn/dtool/utils"
)

func GenPyJiaJia(pe []PyEntry) []byte {
	var buf bytes.Buffer
	for _, v := range pe {
		words := []rune(v.Word)
		if len(words) != len(v.Codes) {
			continue
		}
		for i := 0; i < len(words); i++ {
			buf.WriteString(string(words[i]))
			buf.WriteString(v.Codes[i])
		}
		buf.WriteString(LineBreak)
	}

	return buf.Bytes()
}

// 模式一，只返回频率最高的拼音
// 模式二，作笛卡尔积
func ParsePyJiaJia(rd io.Reader) []PyEntry {
	ret := make([]PyEntry, 0, 0xff)

	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		tmp := strings.TrimSpace(scan.Text())
		entry := []rune(tmp)
		// 注释
		if entry[0] == rune(';') {
			continue
		}
		var word []rune
		var codes []string
		for i := 0; i < len(entry); {
			char := entry[i]
			word = append(word, char)

			// 下一个是英文（拼音）
			if i+1 != len(entry) && entry[i+1] < 128 {
				j := 1 // 已匹配的字母数
				for i+j+1 < len(entry) && entry[i+j+1] < 128 {
					j++
				}
				codes = append(codes, string(entry[i+1:i+j+1]))
				i += j + 1
				continue
			}
			// 读到汉字
			py := encoder.CharPyMap[char]
			if len(py) == 0 {
				py = []string{""}
			}
			codes = append(codes, py[0])
			i++
		}
		ret = append(ret, PyEntry{string(word), codes, 1})
	}
	return ret
}
