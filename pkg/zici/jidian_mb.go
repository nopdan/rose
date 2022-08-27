package zici

import (
	"bytes"
	"io/ioutil"

	. "github.com/cxcn/dtool/pkg/util"
)

func ParseJidianMb(filename string) WcTable {
	data, _ := ioutil.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(WcTable, 0, r.Len()>>8)
	var tmp []byte

	r.Seek(0x1B620, 0) // 从 0x1B620 开始读
	for r.Len() > 3 {
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
		word, _ := Decode(tmp, "utf16")

		ret = append(ret, WordCode{word, code})
	}
	return ret
}
