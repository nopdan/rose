package pinyin

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"

	. "github.com/cxcn/dtool/utils"
)

// 通用规则
type GenRule struct {
	Sep   byte // 分隔符
	PySep byte // 拼音分隔符

	// w 词，c 无前缀拼音，p 有前缀拼音，f 词频
	Rule string
}

func GenGeneral(pe []PyEntry, rule GenRule) []byte {
	var buf bytes.Buffer

	for _, v := range pe {
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
				codes := strings.Join(v.Codes, string(rule.PySep))
				buf.WriteString(codes)
			}
		}
		buf.WriteString(LineBreak)
	}
	return buf.Bytes()
}

func ParseGeneral(rd io.Reader, rule GenRule) []PyEntry {
	rd, err := Decode(rd)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make([]PyEntry, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		e := strings.Split(scan.Text(), string(rule.Sep))

		// TODO: 纯词生成拼音
		if len(e) < 2 {
			continue
		}
		var word string
		var codes []string
		freq := 1
		for i := 0; i < len(rule.Rule); i++ {
			switch rule.Rule[i] {
			case 'w':
				word = e[i]
			case 'f':
				freq, _ = strconv.Atoi(e[i])
			case 'c', 'p':
				tmp := strings.TrimLeft(e[i], string(rule.PySep))
				codes = strings.Split(tmp, string(rule.PySep))
			}
		}
		ret = append(ret, PyEntry{word, codes, freq})
	}
	return ret
}
