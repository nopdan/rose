package pinyin

import (
	"bytes"
	"os"

	"github.com/cxcn/dtool/pkg/util"
)

type MspyDat struct{}

func (MspyDat) Parse(filename string) Dict {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(Dict, 0, r.Len()>>8)
	var tmp []byte

	// 词库偏移量
	r.Seek(0x14, 0)
	phrase_start := ReadUint32(r)
	// fmt.Printf("%x", phrase_start)
	// phrase_end
	ReadUint32(r)
	// 词条数
	phrase_count := ReadUint32(r)

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
		code, _ := util.Decode(codeSli, "UTF-16LE")

		wordSli := make([]byte, 0, 2)
		for {
			r.Read(tmp)
			if tmp[0] == 0 && tmp[1] == 0 {
				break
			}
			wordSli = append(wordSli, tmp...)
		}
		word, _ := util.Decode(wordSli, "UTF-16LE")

		ret = append(ret, Entry{word, []string{code}, 1})
	}
	return ret
}
