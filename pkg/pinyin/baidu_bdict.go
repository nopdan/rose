package pinyin

import (
	"bytes"

	"github.com/nopdan/rose/pkg/util"
)

type BaiduBdict struct {
	Template
	smList []string
	ymList []string
}

func init() {
	FormatList = append(FormatList, NewBaiduBdict(), NewBaiduBcd())
}
func NewBaiduBdict() *BaiduBdict {
	f := new(BaiduBdict)
	f.Name = "百度分类词库.bdict"
	f.ID = "bdict"
	f.smList = []string{
		"c", "d", "b", "f", "g", "h", "ch", "j", "k", "l", "m", "n",
		"", "p", "q", "r", "s", "t", "sh", "zh", "w", "x", "y", "z",
	}
	f.ymList = []string{
		"uang", "iang", "iong", "ang", "eng", "ian", "iao", "ing", "ong",
		"uai", "uan", "ai", "an", "ao", "ei", "en", "er", "ua", "ie", "in", "iu",
		"ou", "ia", "ue", "ui", "un", "uo", "a", "e", "i", "o", "u", "v",
	}
	return f
}

func NewBaiduBcd() *BaiduBdict {
	f := NewBaiduBdict()
	f.Name = "百度手机分类词库.bcd"
	f.ID = "bcd"
	return f
}

func (f *BaiduBdict) Unmarshal(r *bytes.Reader) []*Entry {
	r.Seek(0x70, 0)
	count := ReadUint32(r) // 词条数
	di := make([]*Entry, 0, count)
	r.Seek(0x90, 0)
	util.Info(r, 0xD0-0x90, "词库名: ")
	util.Info(r, 0x110-0xD0, "词库作者: ")
	util.Info(r, 0x150-0x110, "示例词: ")
	util.Info(r, 0x350-0x150, "词库描述: ")

	for r.Len() > 4 {
		// 拼音长
		pyLen := ReadIntN(r, 2)
		// 词频
		freq := ReadIntN(r, 2)

		// 判断下两个字节
		tmp := ReadN(r, 2)

		// 编码和词不等长，全按 utf-16le
		if tmp[0] == 0 && tmp[1] == 0 {
			wordLen := ReadIntN(r, 2)
			// 读编码
			tmp = ReadN(r, pyLen*2)
			code := util.DecodeMust(tmp, "UTF-16LE")
			// 读词
			tmp = ReadN(r, wordLen*2)
			word := util.DecodeMust(tmp, "UTF-16LE")

			di = append(di, &Entry{word, []string{code}, freq})
			continue
		}

		// 全英文的词，编码和词是一样的
		if int(tmp[0]) >= len(f.smList) && tmp[0] != 0xff {
			r.Seek(-2, 1)
			eng := ReadN(r, pyLen)
			di = append(di, &Entry{string(eng), []string{string(eng)}, freq})
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
			pinyin = append(pinyin, f.smList[smIdx]+f.ymList[ymIdx])
		}
		// 读词
		tmp = ReadN(r, pyLen*2)
		word := util.DecodeMust(tmp, "UTF-16LE")
		di = append(di, &Entry{word, pinyin, freq})
	}
	return di
}
