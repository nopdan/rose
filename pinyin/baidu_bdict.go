package pinyin

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/cxcn/dtool/utils"
)

var bdictSm = []string{
	"c", "d", "b", "f", "g", "h", "ch", "j", "k", "l", "m", "n",
	"", "p", "q", "r", "s", "t", "sh", "zh", "w", "x", "y", "z",
}

var bdictYm = []string{
	"uang", "iang", "iong", "ang", "eng", "ian", "iao", "ing", "ong",
	"uai", "uan", "ai", "an", "ao", "ei", "en", "er", "ua", "ie", "in", "iu",
	"ou", "ia", "ue", "ui", "un", "uo", "a", "e", "i", "o", "u", "v",
}

func ParseBaiduBdict(rd io.Reader) []PyEntry {
	ret := make([]PyEntry, 0, 0xff)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)
	var tmp []byte

	r.Seek(0x350, 0)
	for r.Len() > 4 {
		// 拼音长
		codeLen := ReadInt(r, 2)
		// 词频
		freq := ReadInt(r, 2)

		// 判断下两个字节
		tmp = make([]byte, 2)
		r.Read(tmp)

		// 编码和词不等长，全按 utf-16le
		if tmp[0] == 0 && tmp[1] == 0 {
			wordLen := ReadInt(r, 2)
			// 读编码
			tmp = make([]byte, codeLen*2)
			r.Read(tmp)
			code := string(DecUtf16le(tmp))
			// 读词
			tmp = make([]byte, wordLen*2)
			r.Read(tmp)
			word := string(DecUtf16le(tmp))

			ret = append(ret, PyEntry{word, []string{code}, freq})
			continue
		}

		// 全英文的词，编码和词是一样的
		if int(tmp[0]) >= len(bdictSm) && tmp[0] != 0xff {
			r.Seek(-2, 1)
			eng := make([]byte, codeLen)
			r.Read(eng)
			ret = append(ret, PyEntry{string(eng), []string{string(eng)}, freq})
			continue
		}

		// 一般格式
		r.Seek(-2, 1)
		codes := make([]string, 0, codeLen)
		for i := 0; i < codeLen; i++ {
			smIdx, _ := r.ReadByte()
			ymIdx, _ := r.ReadByte()
			// 带英文的词组
			if smIdx == 0xff {
				codes = append(codes, string(ymIdx))
				continue
			}
			codes = append(codes, bdictSm[smIdx]+bdictYm[ymIdx])
		}
		// 读词
		tmp = make([]byte, 2*codeLen)
		r.Read(tmp)
		word := string(DecUtf16le(tmp))
		ret = append(ret, PyEntry{word, codes, freq})
	}
	return ret
}
