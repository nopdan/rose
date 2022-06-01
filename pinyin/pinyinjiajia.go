package pinyin

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"

	encoder "github.com/cxcn/dtool/encoders"
	. "github.com/cxcn/dtool/utils"
)

// 模式一，只返回频率最高的拼音
// TODO: 模式二，作笛卡尔积
func ParseJiaJia(filename string) WpfDict {
	f, _ := os.Open(filename)
	defer f.Close()
	rd, err := DecodeIO(f)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make(WpfDict, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		tmp := strings.TrimSpace(scan.Text())
		entry := []rune(tmp)
		// 注释
		if entry[0] == rune(';') {
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
			codes := encoder.CharYinjieMap[char]
			if len(codes) == 0 {
				codes = []string{""}
			}
			pinyin = append(pinyin, codes[0])
			i++
		}
		ret = append(ret, WordPyFreq{string(word), pinyin, 1})
	}
	return ret
}

func GenJiaJia(dict WpfDict) []byte {
	var buf bytes.Buffer
	for _, v := range dict {
		words := []rune(v.Word)
		if len(words) != len(v.Pinyin) {
			continue
		}
		for i := 0; i < len(words); i++ {
			buf.WriteString(string(words[i]))
			buf.WriteString(v.Pinyin[i])
		}
		buf.WriteString(LineBreak)
	}
	return buf.Bytes()
}
