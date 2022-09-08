package pinyin

import (
	"bytes"
	"os"

	"github.com/cxcn/dtool/pkg/encoder"
	"github.com/cxcn/dtool/pkg/util"
)

type MspyUDL struct{}

// 自学习词库，纯汉字
func (MspyUDL) Parse(filename string) Dict {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(Dict, 0, r.Len()>>8)
	r.Seek(0xC, 0)
	dictLen := ReadUint32(r)

	for i := 0; i < dictLen; i++ {
		r.Seek(0x2400+60*int64(i), 0)
		r.Seek(10, 1)
		wordLen, _ := r.ReadByte()
		r.ReadByte()
		wordSli := make([]byte, wordLen*2)
		r.Read(wordSli)
		word, _ := util.Decode(wordSli, "utf16")
		ret = append(ret, Entry{word, encoder.GetPinyin(word), 1})
	}
	return ret
}
