package rose

import (
	"bytes"
)

type JidianMb struct{ Dict }

func NewJidianMb() *JidianMb {
	d := new(JidianMb)
	d.IsPinyin = false
	d.IsBinary = true
	d.Name = "极点码表.mb"
	d.Suffix = "mb"
	return d
}

func (d *JidianMb) Parse() {
	table := make(Table, 0, d.size>>8)

	r := bytes.NewReader(d.data)
	r.Seek(0x1B620, 0) // 从 0x1B620 开始读
	for r.Len() > 3 {
		var tmp []byte
		codeLen, _ := r.ReadByte()
		if codeLen == 0xff {
			r.Seek(1, 1)
			continue
		}
		wordLen, _ := r.ReadByte()
		r.Seek(1, 1)

		// 读编码
		tmp = make([]byte, codeLen)
		r.Read(tmp)
		code := string(tmp)

		// 读词
		tmp = make([]byte, wordLen)
		r.Read(tmp)
		word, _ := Decode(tmp, "UTF-16LE")

		table = append(table, &TableEntry{word, code, 1})
	}
	d.table = table
}
