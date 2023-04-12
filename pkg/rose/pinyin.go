package rose

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	util "github.com/flowerime/goutil"
)

// 通用规则
type CommonPyTable struct {
	Dict
	Sep   byte // 分隔符
	PySep byte // 拼音分隔符

	// w 词，c 无前缀拼音，p 有前缀拼音，f 词频
	Rule string
}

func NewCommonPyTable(format string) *CommonPyTable {
	d := new(CommonPyTable)
	d.IsPinyin = true
	d.IsBinary = false
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

func (d *CommonPyTable) GetDict() *Dict {
	return &d.Dict
}

// 拼音通用格式解析
func (d *CommonPyTable) Parse() {
	pyt := make(PyTable, 0, 0xff)
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
		pyt = append(pyt, &PinyinEntry{word, pinyin, freq})
	}
	d.pyt = pyt
}

// 拼音通用格式生成
func (d *CommonPyTable) GenFrom(src *Dict) []byte {
	var buf bytes.Buffer
	for _, v := range src.pyt {
		for i := 0; i < len(d.Rule); i++ {
			if i != 0 {
				buf.WriteByte(d.Sep)
			}
			switch d.Rule[i] {
			case 'w':
				buf.WriteString(v.Word)
			case 'f':
				buf.WriteString(strconv.Itoa(v.Freq))
			case 'c', 'p':
				if d.Rule[i] == 'p' {
					buf.WriteByte(d.PySep)
				}
				pinyin := strings.Join(v.Pinyin, string(d.PySep))
				buf.WriteString(pinyin)
			}
		}
		buf.WriteString(util.LineBreak)
	}
	return buf.Bytes()
}
