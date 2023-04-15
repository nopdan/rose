package rose

import (
	"bytes"
)

type Fcitx4Mb struct{ Dict }

func NewFcitx4Mb() *Fcitx4Mb {
	d := new(Fcitx4Mb)
	d.Name = "fcitx4.mb"
	d.Suffix = "mb"
	return d
}

func (d *Fcitx4Mb) Parse() {
	r := bytes.NewReader(d.data)

	r.Seek(0x55, 0)
	// 词条数
	count := ReadUint32(r)
	wl := make([]Entry, 0, count)

	for i := _u32; i < count; i++ {
		var tmp []byte
		tmp = make([]byte, 5)
		r.Read(tmp)
		code := trimSufZero(tmp)

		wordLen := ReadUint32(r)
		tmp = make([]byte, wordLen-1)
		r.Read(tmp)
		word := string(tmp)

		wl = append(wl, &WubiEntry{word, code, 1})
		r.Seek(10, 1)
	}
	d.WordLibrary = wl
}

func trimSufZero(b []byte) string {
	for i := len(b); i > 0; i-- {
		if b[i-1] != 0 {
			return string(b[:i])
		}
	}
	return ""
}
