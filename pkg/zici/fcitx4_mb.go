package zici

import (
	"bytes"
	"io/ioutil"

	. "github.com/cxcn/dtool/pkg/util"
)

func ParseFcitx4Mb(filename string) WcTable {
	data, _ := ioutil.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(WcTable, 0, r.Len()>>8)
	var tmp []byte

	r.Seek(0x55, 0)
	// 词条数
	dictLen := ReadUint32(r)

	for i := 0; i < dictLen; i++ {
		tmp = make([]byte, 5)
		r.Read(tmp)
		code := trimSufZero(tmp)

		wordLen := ReadUint32(r)
		tmp = make([]byte, wordLen-1)
		r.Read(tmp)
		word := string(tmp)

		ret = append(ret, WordCode{word, code})
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
