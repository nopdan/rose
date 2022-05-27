package zici

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/cxcn/dtool/utils"
)

func ParseJidianMb(rd io.Reader) []ZcEntry {
	ret := make([]ZcEntry, 0, 1e5) // 初始化
	data, _ := ioutil.ReadAll(rd)  // 全部读到内存
	r := bytes.NewReader(data)
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
		word := string(DecUtf16le(tmp))

		ret = append(ret, ZcEntry{word, code})
	}
	return ret
}
