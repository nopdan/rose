package pinyin

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

// 通用规则
type GenRule struct {
	Sep       byte // 分隔符
	PySep     byte // 拼音分隔符
	WordFirst bool // 词在前
	HasPy     bool // 有无拼音
	PrePySep  bool // 拼音开头有无分隔符
	HasFreq   bool // 有无词频
}

func NewGenRule(sep, pySep byte, rule int) GenRule {
	var ret GenRule
	ret.Sep = sep
	ret.PySep = pySep
	if rule&0b1000 == 0b1000 {
		ret.WordFirst = true
	}
	if rule&0b100 == 0b100 {
		ret.HasPy = true
	}
	if rule&0b10 == 0b10 {
		ret.PrePySep = true
	}
	if rule&0b1 == 0b1 {
		ret.HasFreq = true
	}
	return ret
}

func GenGenersal(pe []PyEntry, rl GenRule) []byte {
	var buf bytes.Buffer
	for _, v := range pe {
		if rl.WordFirst {
			buf.WriteString(v.Word)
			if rl.HasPy {
				buf.WriteByte(rl.Sep)
				if rl.PrePySep {
					buf.WriteByte(rl.PySep)
				}
				buf.WriteString(strings.Join(v.Codes, string(rl.PySep)))
			}
		} else {
			if rl.HasPy {
				if rl.PrePySep {
					buf.WriteByte(rl.PySep)
				}
				buf.WriteString(strings.Join(v.Codes, string(rl.PySep)))
				buf.WriteByte(rl.Sep)
			}
			buf.WriteString(v.Word)
		}
		if rl.HasFreq {
			buf.WriteByte(rl.Sep)
			buf.WriteString(strconv.Itoa(v.Freq))
		}
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}

func ParseGeneral(rd io.Reader, rl GenRule) []PyEntry {
	ret := make([]PyEntry, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		e := strings.Split(scan.Text(), string(rl.Sep))

		// TODO: 纯词生成拼音
		if len(e) < 2 {
			continue
		}
		var word, codes string
		var freq int
		if rl.WordFirst {
			word = e[0]
			codes = e[1]
			if len(e) == 3 && rl.HasFreq {
				tmp := e[2]
				freq, _ = strconv.Atoi(tmp)
			}
		} else {
			word = e[1]
			codes = e[0]
			if len(e) == 3 && rl.HasFreq {
				tmp := e[2]
				freq, _ = strconv.Atoi(tmp)
			}
		}
		codes = strings.TrimLeft(codes, string(rl.PySep))
		pys := strings.Split(codes, string(rl.PySep))
		ret = append(ret, PyEntry{word, pys, freq})
	}
	return ret
}
