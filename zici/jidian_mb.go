package zici

import (
	"bytes"
	"io"
	"io/ioutil"

	"golang.org/x/text/encoding/unicode"
)

func ParseJidianMb(rd io.Reader) Dict {

	ret := make(Dict, 1e5)       // 初始化
	tmp, _ := ioutil.ReadAll(rd) // 全部读到内存
	r := bytes.NewReader(tmp)
	r.Seek(0x1B620, 0) // 从 0x1B620 开始读
	// utf-16le 转换
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	for r.Len() > 3 {
		codeLen, _ := r.ReadByte()
		if codeLen == 0xff {
			r.Seek(1, 1)
			continue
		}
		wordLen, _ := r.ReadByte()

		r.Seek(1, 1) // 丢掉一个字节
		codeSli := make([]byte, codeLen)
		r.Read(codeSli)
		wordSli := make([]byte, wordLen)
		r.Read(wordSli)
		wordSli, _ = decoder.Bytes(wordSli)
		ret.insert(string(codeSli), string(wordSli))
	}

	return ret
}
