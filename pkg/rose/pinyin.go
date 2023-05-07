package rose

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"github.com/nopdan/ku"
)

// 通用规则
type Pinyin struct {
	Dict
	Sep   byte // 分隔符
	PySep byte // 拼音分隔符

	// w 词，c 无前缀拼音，p 有前缀拼音，f 词频
	Rule string
}

func NewPinyin(format string) *Pinyin {
	d := new(Pinyin)
	switch format {
	case "sg":
		d.Sep = ' '
		d.PySep = '\''
		d.Rule = "pw"
		d.Name = "搜狗拼音.txt"
	case "qq":
		d.Sep = ' '
		d.PySep = '\''
		d.Rule = "cwf"
		d.Name = "QQ 拼音.txt"
	case "bd":
		d.Sep = '\t'
		d.PySep = '\''
		d.Rule = "wcf"
		d.Name = "百度拼音.txt"
	case "gg":
		d.Sep = '\t'
		d.PySep = ' '
		d.Rule = "wfc"
		d.Name = "谷歌拼音.txt"
	case "rime":
		d.Sep = '\t'
		d.PySep = ' '
		d.Rule = "wcf"
		d.Name = "rime 拼音.txt"
	}
	return d
}

// 拼音通用格式解析
func (d *Pinyin) Parse() {
	wl := make([]Entry, 0, d.size>>8)

	scan := bufio.NewScanner(d.rd)
	for scan.Scan() {
		e := strings.Split(scan.Text(), string(d.Sep))
		// TODO: 纯词生成拼音
		if len(e) < 2 {
			continue
		}
		var word string
		pinyin := make([]string, 0, 1)
		freq := 1
		for i := 0; i < len(d.Rule); i++ {
			if i >= len(e) {
				continue
			}
			switch d.Rule[i] {
			case 'w':
				word = e[i]
			case 'f':
				freq, _ = strconv.Atoi(e[i])
			case 'c', 'p':
				tmp := strings.TrimLeft(e[i], string(d.PySep))
				pinyin = strings.Split(tmp, string(d.PySep))
			}
		}
		wl = append(wl, &PinyinEntry{word, pinyin, freq})
	}
	d.WordLibrary = wl
}

// 拼音通用格式生成
func (d *Pinyin) GenFrom(wl WordLibrary) []byte {
	var buf bytes.Buffer
	for _, v := range wl {
		for i := 0; i < len(d.Rule); i++ {
			if i != 0 {
				buf.WriteByte(d.Sep)
			}
			switch d.Rule[i] {
			case 'w':
				buf.WriteString(v.GetWord())
			case 'f':
				buf.WriteString(strconv.Itoa(v.GetFreq()))
			case 'c', 'p':
				if d.Rule[i] == 'p' {
					buf.WriteByte(d.PySep)
				}
				pinyin := strings.Join(v.GetPinyin(), string(d.PySep))
				buf.WriteString(pinyin)
			}
		}
		buf.WriteString(ku.LineBreak)
	}
	return buf.Bytes()
}
