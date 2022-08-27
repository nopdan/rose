package pinyin

import (
	"bytes"
	"io/ioutil"

	. "github.com/cxcn/dtool/pkg/util"
)

func ParseMspyDat(filename string) WpfDict {
	data, _ := ioutil.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(WpfDict, 0, r.Len()>>8)
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
		code, _ := Decode(codeSli, "utf16")

		wordSli := make([]byte, 0, 2)
		for {
			r.Read(tmp)
			if tmp[0] == 0 && tmp[1] == 0 {
				break
			}
			wordSli = append(wordSli, tmp...)
		}
		word, _ := Decode(wordSli, "utf16")

		ret = append(ret, WordPyFreq{word, []string{code}, 1})
	}
	return ret
}
