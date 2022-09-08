package table

import (
	"bytes"
	"os"

	"github.com/cxcn/dtool/pkg/util"
)

type JidianMb struct{}

func (JidianMb) Parse(filename string) Table {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(Table, 0, r.Len()>>8)
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
		word, _ := util.Decode(tmp, "utf16")

		ret = append(ret, Entry{word, code})
	}
	return ret
}
