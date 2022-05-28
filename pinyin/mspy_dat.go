package pinyin

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/cxcn/dtool/utils"
)

func ParseMspyDat(rd io.Reader) []PyEntry {

	ret := make([]PyEntry, 0, 0xff)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)
	var tmp []byte

	// 词库偏移量
	r.Seek(0x14, 0)
	phrase_start := ReadInt(r, 4)
	// fmt.Printf("%x", phrase_start)
	// phrase_end
	ReadInt(r, 4)
	// 词条数
	phrase_count := ReadInt(r, 4)

	r.Seek(int64(phrase_start), 0)
	for i := 0; i < phrase_count; i++ {
		// 这 16 个字节有信息，但没什么用
		r.Seek(16, 1)
		tmp = make([]byte, 2)

		// 编码，居然是连续的，分隔符都没有
		codeSli := make([]byte, 0, 2)
		for {
			r.Read(tmp)
			if tmp[0] == 0 && tmp[1] == 0 {
				break
			}
			codeSli = append(codeSli, tmp...)
		}

		wordSli := make([]byte, 0, 2)
		for {
			r.Read(tmp)
			if tmp[0] == 0 && tmp[1] == 0 {
				break
			}
			wordSli = append(wordSli, tmp...)
		}

		ret = append(ret, PyEntry{string(DecUtf16le(wordSli)), []string{string(DecUtf16le(codeSli))}, 1})
	}
	return ret
}
