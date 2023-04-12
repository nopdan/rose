package rose

import (
	"bufio"
	"fmt"
	"strings"
)

type CommonTable struct {
	Dict
	WordFirst bool
	Encoding  string
}

func NewCommonTable(format string) *CommonTable {
	d := new(CommonTable)
	d.IsPinyin = false
	d.IsBinary = false
	if format == "dd" {
		d.WordFirst = true
		d.Encoding = "UTF-8"
		d.Name = "多多.txt"
	} else if format == "bl" {
		d.WordFirst = false
		d.Encoding = "UTF-16LE"
		d.Name = "冰凌.txt"
	}
	return d
}

func (d *CommonTable) GetDict() *Dict {
	return &d.Dict
}

func (d *CommonTable) Parse() {
	table := make(Table, 0, d.size>>8)

	scan := bufio.NewScanner(d.rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		word, code := entry[0], entry[1]
		if !d.WordFirst {
			word, code = code, word
		}
		if strings.HasPrefix(word, "$ddcmd") {
			fmt.Println("多多的命令" + word)
			continue
		}
		table = append(table, &TableEntry{word, code, 1})
	}
	d.table = table
}

func (d *CommonTable) GenFrom(src *Dict) []byte {
	var sb strings.Builder
	sb.Grow(len(src.table))
	for _, v := range src.table {
		if d.WordFirst {
			sb.WriteString(v.Word)
			sb.WriteByte('\t')
			sb.WriteString(v.Code)
		} else {
			sb.WriteString(v.Code)
			sb.WriteByte('\t')
			sb.WriteString(v.Word)
		}
		sb.WriteByte('\r')
		sb.WriteByte('\n')
	}
	ret, _ := Encode(sb.String(), d.Encoding)
	if d.Encoding == "UTF-16LE" {
		ret = append([]byte{0xff, 0xfe}, ret...)
	}
	return ret
}
