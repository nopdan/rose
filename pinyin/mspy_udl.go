package pinyin

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/cxcn/dtool/utils"
)

func ParseMspyUDL(rd io.Reader) []string {
	ret := make([]string, 0, 0xff)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)
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
		ret = append(ret, word)
	}
	return ret
}
