package pinyin

import (
	"bytes"
	"os"

	"github.com/cxcn/dtool/pkg/util"
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

type BaiduBdict struct{}

func (BaiduBdict) Parse(filename string) Dict {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(Dict, 0, r.Len()>>8)
	var tmp []byte

	r.Seek(0x350, 0)
	for r.Len() > 4 {
		// 拼音长
		pyLen := ReadUint16(r)
		// 词频
		freq := ReadUint16(r)

		// 判断下两个字节
		tmp = make([]byte, 2)
		r.Read(tmp)

		// 编码和词不等长，全按 utf-16le
		if tmp[0] == 0 && tmp[1] == 0 {
			wordLen := ReadUint16(r)
			// 读编码
			tmp = make([]byte, pyLen*2)
			r.Read(tmp)
			code, _ := util.Decode(tmp, "UTF-16LE")
			// 读词
			tmp = make([]byte, wordLen*2)
			r.Read(tmp)
			word, _ := util.Decode(tmp, "UTF-16LE")

			ret = append(ret, Entry{
				Word:   word,
				Pinyin: []string{code},
				Freq:   freq,
			})
			continue
		}

		// 全英文的词，编码和词是一样的
		if int(tmp[0]) >= len(bdictSm) && tmp[0] != 0xff {
			r.Seek(-2, 1)
			eng := make([]byte, pyLen)
			r.Read(eng)
			ret = append(ret, Entry{
				Word:   string(eng),
				Pinyin: []string{string(eng)},
				Freq:   freq,
			})
			continue
		}

		// 一般格式
		r.Seek(-2, 1)
		pinyin := make([]string, 0, pyLen)
		for i := 0; i < pyLen; i++ {
			smIdx, _ := r.ReadByte()
			ymIdx, _ := r.ReadByte()
			// 带英文的词组
			if smIdx == 0xff {
				pinyin = append(pinyin, string(ymIdx))
				continue
			}
			pinyin = append(pinyin, bdictSm[smIdx]+bdictYm[ymIdx])
		}
		// 读词
		tmp = make([]byte, pyLen*2)
		r.Read(tmp)
		word, _ := util.Decode(tmp, "UTF-16LE")
		ret = append(ret, Entry{
			Word:   word,
			Pinyin: pinyin,
			Freq:   freq,
		})
	}
	return ret
}
