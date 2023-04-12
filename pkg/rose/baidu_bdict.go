package rose

import (
	"bytes"
)

type BaiduBdict struct{ Dict }

func NewBaiduBdict() *BaiduBdict {
	d := new(BaiduBdict)
	d.IsPinyin = true
	d.IsBinary = true
	d.Name = "百度分类词库.bdict(.bcd)"
	d.Suffix = "bdict"
	return d
}

func (d *BaiduBdict) GetDict() *Dict {
	return &d.Dict
}

var bdictSm = []string{
	"c", "d", "b", "f", "g", "h", "ch", "j", "k", "l", "m", "n",
	"", "p", "q", "r", "s", "t", "sh", "zh", "w", "x", "y", "z",
}

var bdictYm = []string{
	"uang", "iang", "iong", "ang", "eng", "ian", "iao", "ing", "ong",
	"uai", "uan", "ai", "an", "ao", "ei", "en", "er", "ua", "ie", "in", "iu",
	"ou", "ia", "ue", "ui", "un", "uo", "a", "e", "i", "o", "u", "v",
}

func (d *BaiduBdict) Parse() {
	pyt := make(PyTable, 0, d.size>>8)

	r := bytes.NewReader(d.data)
	r.Seek(0x350, 0)
	for r.Len() > 4 {
		var tmp []byte
		// 拼音长
		pyLen := ReadUint16(r)
		// 词频
		freq := int(ReadUint16(r))

		// 判断下两个字节
		tmp = make([]byte, 2)
		r.Read(tmp)

		// 编码和词不等长，全按 utf-16le
		if tmp[0] == 0 && tmp[1] == 0 {
			wordLen := ReadUint16(r)
			// 读编码
			tmp = make([]byte, pyLen*2)
			r.Read(tmp)
			code, _ := Decode(tmp, "UTF-16LE")
			// 读词
			tmp = make([]byte, wordLen*2)
			r.Read(tmp)
			word, _ := Decode(tmp, "UTF-16LE")

			pyt = append(pyt, &PinyinEntry{word, []string{code}, freq})
			continue
		}

		// 全英文的词，编码和词是一样的
		if int(tmp[0]) >= len(bdictSm) && tmp[0] != 0xff {
			r.Seek(-2, 1)
			eng := make([]byte, pyLen)
			r.Read(eng)
			pyt = append(pyt, &PinyinEntry{string(eng), []string{string(eng)}, freq})
			continue
		}

		// 一般格式
		r.Seek(-2, 1)
		pinyin := make([]string, 0, pyLen)
		for i := _u16; i < pyLen; i++ {
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
		word, _ := Decode(tmp, "UTF-16LE")
		pyt = append(pyt, &PinyinEntry{word, pinyin, freq})
	}
	d.pyt = pyt
}

func (d *BaiduBdict) GenFrom(src *Dict) []byte {
	return genErr(d.Name)
}
