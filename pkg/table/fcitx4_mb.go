package table

import (
	"bytes"
	"os"
)

type Fcitx4Mb struct{}

func (Fcitx4Mb) Parse(filename string) Table {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	var tmp []byte

	r.Seek(0x55, 0)
	// 词条数
	count := ReadUint32(r)
	ret := make(Table, 0, count)

	for i := _u32; i < count; i++ {
		tmp = make([]byte, 5)
		r.Read(tmp)
		code := trimSufZero(tmp)

		wordLen := ReadUint32(r)
		tmp = make([]byte, wordLen-1)
		r.Read(tmp)
		word := string(tmp)

		ret = append(ret, Entry{word, code, 1})
		r.Seek(10, 1)
	}
	return ret
}

func trimSufZero(b []byte) string {
	for i := len(b); i > 0; i-- {
		if b[i-1] != 0 {
			return string(b[:i])
		}
	}
	return ""
}
