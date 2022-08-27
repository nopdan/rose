package pinyin

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"

	. "github.com/cxcn/dtool/pkg/util"
)

// 通用规则
type PinyinRule struct {
	Sep   byte // 分隔符
	PySep byte // 拼音分隔符

	// w 词，c 无前缀拼音，p 有前缀拼音，f 词频
	Rule string
}

var (
	R_sogou     = PinyinRule{' ', '\'', "pw"}
	R_qq        = PinyinRule{' ', '\'', "cwf"}
	R_baidu     = PinyinRule{'\t', '\'', "wcf"}
	R_google    = PinyinRule{'\t', ' ', "wfc"}
	R_word_only = PinyinRule{' ', ' ', "w"}
)

// 拼音通用格式解析
func ParsePinyin(filename string, rule PinyinRule) WpfDict {
	f, _ := os.Open(filename)
	defer f.Close()
	rd, err := DecodeIO(f)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make(WpfDict, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		e := strings.Split(scan.Text(), string(rule.Sep))
		// TODO: 纯词生成拼音
		if len(e) < 2 {
			continue
		}
		var word string
		pinyin := make([]string, 0, 1)
		freq := 1
		for i := 0; i < len(rule.Rule); i++ {
			switch rule.Rule[i] {
			case 'w':
				word = e[i]
			case 'f':
				freq, _ = strconv.Atoi(e[i])
			case 'c', 'p':
				tmp := strings.TrimLeft(e[i], string(rule.PySep))
				pinyin = strings.Split(tmp, string(rule.PySep))
			}
		}
		ret = append(ret, WordPyFreq{word, pinyin, freq})
	}
	return ret
}

// 拼音通用格式生成
func GenPinyin(dict WpfDict, rule PinyinRule) []byte {
	var buf bytes.Buffer
	for _, v := range dict {
		for i := 0; i < len(rule.Rule); i++ {
			if i != 0 {
				buf.WriteByte(rule.Sep)
			}
			switch rule.Rule[i] {
			case 'w':
				buf.WriteString(v.Word)
			case 'f':
				buf.WriteString(strconv.Itoa(v.Freq))
			case 'c', 'p':
				if rule.Rule[i] == 'p' {
					buf.WriteByte(rule.PySep)
				}
				pinyin := strings.Join(v.Pinyin, string(rule.PySep))
				buf.WriteString(pinyin)
			}
		}
		buf.WriteString(LineBreak)
	}
	return buf.Bytes()
}
