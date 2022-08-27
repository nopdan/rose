package pinyin

import (
	"bytes"
	"io/ioutil"

	"github.com/cxcn/dtool/pkg/encoder"
	. "github.com/cxcn/dtool/pkg/util"
)

// 自学习词库，纯汉字
func ParseMspyUDL(filename string) WpfDict {
	data, _ := ioutil.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(WpfDict, 0, r.Len()>>8)
	r.Seek(0xC, 0)
	dictLen := ReadUint32(r)

	for i := 0; i < dictLen; i++ {
		r.Seek(0x2400+60*int64(i), 0)
		r.Seek(10, 1)
		wordLen, _ := r.ReadByte()
		r.ReadByte()
		wordSli := make([]byte, wordLen*2)
		r.Read(wordSli)
		word, _ := Decode(wordSli, "utf16")
		ret = append(ret, WordPyFreq{word, encoder.GetPinyin(word), 1})
	}
	return ret
}
