package rose

import (
	"bufio"
	"bytes"
	"strings"
)

type Jidian struct{ Dict }

func NewJidian() *Jidian {
	d := new(Jidian)
	d.IsPinyin = false
	d.IsBinary = false
	d.Name = "极点码表.txt"
	return d
}

func (d *Jidian) Parse() {
	table := make(CodeTable, 0, d.size>>8)

	scan := bufio.NewScanner(d.rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), " ")
		if len(entry) < 2 {
			continue
		}
		table = append(table, &CodeEntry{entry[0], entry[1:]})
	}
	d.codet = table
}

func (Jidian) GenFrom(d *Dict) []byte {
	d.ToCodeTable()
	var buf bytes.Buffer
	for _, v := range d.codet {
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Words, " "))
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}
