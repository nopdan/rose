package pinyin

import (
	"bytes"
	"io"
	"io/ioutil"

	"golang.org/x/text/encoding/unicode"
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

func ParseBaiduBdict(rd io.Reader) []Pinyin {
	ret := make([]Pinyin, 0, 1e5)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)

	// utf-16le 转换器
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()

	r.Seek(0x350, 0)
	for r.Len() > 4 {
		// 拼音长，也是词长
		codeLen, _ := r.ReadByte()
		r.Seek(3, 1) // 丢掉跟着的3个0

		// 判断下两个字节
		tmp := make([]byte, 2)
		r.Read(tmp)
		if tmp[0] == 0 && tmp[1] == 0 {
			r.Read(tmp)
			wordLen := bytesToInt(tmp)
			codeSli := make([]byte, codeLen*2)
			r.Read(codeSli)
			wordSli := make([]byte, wordLen*2)
			r.Read(wordSli)
			codeSli, _ = decoder.Bytes(codeSli)
			wordSli, _ = decoder.Bytes(wordSli)
			ret = append(ret, Pinyin{string(wordSli), []string{string(codeSli)}, 1})
			continue
		}

		// 全英文的词，编码和词是一样的
		if int(tmp[0]) >= len(bdictSm) && tmp[0] != 0xff {
			r.Seek(-2, 1)
			eng := make([]byte, codeLen)
			r.Read(eng)
			ret = append(ret, Pinyin{string(eng), []string{string(eng)}, 1})
			continue
		}

		r.Seek(-2, 1)
		// 一般格式
		code := make([]string, 0, codeLen)
		for i := 0; i < int(codeLen); i++ {
			smIdx, _ := r.ReadByte()
			ymIdx, _ := r.ReadByte()
			// 带英文的词组
			if smIdx == 0xff {
				code = append(code, string(ymIdx))
				continue
			}
			code = append(code, bdictSm[smIdx]+bdictYm[ymIdx])
		}
		// 读词
		wordSli := make([]byte, 2*codeLen)
		r.Read(wordSli)
		wordSli, _ = decoder.Bytes(wordSli)
		ret = append(ret, Pinyin{string(wordSli), code, 1})
	}
	return ret
}
