package zici

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/cxcn/dtool/utils"
)

func ParseFcitx4Mb(rd io.Reader) []ZcEntry {
	ret := make([]ZcEntry, 1e5)   // 初始化
	data, _ := ioutil.ReadAll(rd) // 全部读到内存
	r := bytes.NewReader(data)
	var tmp []byte

	r.Seek(0x55, 0)
	// 词条数
	dictLen := ReadUint32(r)

	for i := 0; i < dictLen; i++ {
		tmp = make([]byte, 5)
		r.Read(tmp)
		code := TrimSufZero(tmp)

		wordLen := ReadUint32(r)
		tmp = make([]byte, wordLen-1)
		r.Read(tmp)
		word := string(tmp)

		ret = append(ret, ZcEntry{word, code})
		r.Seek(10, 1)
	}
	return ret
}

func TrimSufZero(b []byte) string {
	for i := len(b); i > 0; i-- {
		if b[i] != 0 {
			return string(b[:i])
		}
	}
	return ""
}
