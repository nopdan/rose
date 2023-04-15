package rose

import (
	"bytes"
)

type JidianMb struct{ Dict }

func NewJidianMb() *JidianMb {
	d := new(JidianMb)
	d.Name = "极点码表.mb"
	d.Suffix = "mb"
	return d
}

func (d *JidianMb) Parse() {
	wl := make([]Entry, 0, d.size>>8)

	r := bytes.NewReader(d.data)
	r.Seek(0x17, 0)
	PrintInfo(r, 0x11F-0x17, "")
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

		wl = append(wl, &WubiEntry{word, code, 1})
	}
	d.WordLibrary = wl
}
