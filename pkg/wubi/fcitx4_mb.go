package wubi

import (
	"bytes"
)

type Fcitx4Mb struct{ Template }

func NewFcitx4Mb() *Fcitx4Mb {
	f := new(Fcitx4Mb)
	f.Name = "fcitx4.mb"
	f.ID = "fcitx4"
	return f
}

func (Fcitx4Mb) Unmarshal(r *bytes.Reader) []*Entry {
	r.Seek(0x55, 0)
	// 词条数
	count := ReadIntN(r, 4)
	di := make([]*Entry, 0, count)

	for i := 0; i < count; i++ {
		codeBytes := ReadN(r, 5)
		code := trimSufZero(codeBytes)

		wordSize := ReadIntN(r, 4) - 1
		wordBytes := ReadN(r, wordSize)
		word := string(wordBytes)

		di = append(di, &Entry{
			Word: word,
			Code: code,
		})
		r.Seek(10, 1)
	}
	return di
}

// 去掉末尾的 0
func trimSufZero(b []byte) string {
	for i := len(b); i > 0; i-- {
		if b[i-1] != 0 {
			return string(b[:i])
		}
	}
	return ""
}
