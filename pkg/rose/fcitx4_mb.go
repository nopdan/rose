package rose

import (
	"bytes"
)

type Fcitx4Mb struct{ Dict }

func NewFcitx4Mb() *Fcitx4Mb {
	d := new(Fcitx4Mb)
	d.IsPinyin = false
	d.IsBinary = true
	d.Name = "fcitx4.mb"
	d.Suffix = "mb"
	return d
}

func (d *Fcitx4Mb) Parse() {
	r := bytes.NewReader(d.data)

	r.Seek(0x55, 0)
	// 词条数
	count := ReadUint32(r)
	table := make(Table, 0, count)

	for i := _u32; i < count; i++ {
		var tmp []byte
		tmp = make([]byte, 5)
		r.Read(tmp)
		code := trimSufZero(tmp)

		wordLen := ReadUint32(r)
		tmp = make([]byte, wordLen-1)
		r.Read(tmp)
		word := string(tmp)

		table = append(table, &TableEntry{word, code, 1})
		r.Seek(10, 1)
	}
	d.table = table
}

func trimSufZero(b []byte) string {
	for i := len(b); i > 0; i-- {
		if b[i-1] != 0 {
			return string(b[:i])
		}
	}
	return ""
}
