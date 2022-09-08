package table

import (
	"bytes"
	"os"

	"github.com/cxcn/dtool/pkg/util"
)

type Fcitx4Mb struct{}

func (Fcitx4Mb) Parse(filename string) Table {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(Table, 0, r.Len()>>8)
	var tmp []byte

	r.Seek(0x55, 0)
	// 词条数
	dictLen := util.ReadUint32(r)

	for i := 0; i < dictLen; i++ {
		tmp = make([]byte, 5)
		r.Read(tmp)
		code := trimSufZero(tmp)

		wordLen := util.ReadUint32(r)
		tmp = make([]byte, wordLen-1)
		r.Read(tmp)
		word := string(tmp)

		ret = append(ret, Entry{word, code})
		r.Seek(10, 1)
	}
	return ret
}

func trimSufZero(b []byte) string {
	for i := len(b); i > 0; i-- {
		if b[i] != 0 {
			return string(b[:i])
		}
	}
	return ""
}
