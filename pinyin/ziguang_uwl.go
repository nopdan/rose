package pinyin

import (
	"bytes"

	"github.com/nopdan/rose/util"
)

type ZiguangUwl struct {
	Template
	encoding string
	smList   []string
	ymList   []string
}

func init() {
	FormatList = append(FormatList, NewZiguangUwl())
}
func NewZiguangUwl() *ZiguangUwl {
	f := new(ZiguangUwl)
	f.Name = "紫光华宇拼音.uwl"
	f.ID = "uwl"
	f.smList = []string{
		"", "b", "c", "ch", "d", "f", "g", "h", "j", "k", "l", "m", "n",
		"p", "q", "r", "s", "sh", "t", "w", "x", "y", "z", "zh",
	}
	f.ymList = []string{
		"ang", "a", "ai", "an", "ang", "ao", "e", "ei", "en", "eng", "er",
		"i", "ia", "ian", "iang", "iao", "ie", "in", "ing", "iong", "iu",
		"o", "ong", "ou", "u",
		"ua", "uai", "uan", "uang", "ue", "ui", "un", "uo", "v",
	}
	return f
}

func (f *ZiguangUwl) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)

	r.Seek(2, 0)
	// 编码格式，08 为 GBK，09 为 UTF-16LE
	encoding, _ := r.ReadByte()
	switch encoding {
	case 0x08:
		f.encoding = "GBK"
	case 0x09:
		f.encoding = "UTF-16LE"
	}

	// 分段
	r.Seek(0x48, 0)
	partLen := ReadUint32(r)
	for i := _u32; i < partLen; i++ {
		r.Seek(0xC00+int64(i)<<10, 0)
		r.Seek(12, 1)
		part := f.parse(r)
		di = append(di, part...)
	}
	return di
}

func (f *ZiguangUwl) parse(r *bytes.Reader) []*Entry {
	// 词条占用字节数
	max := ReadUint32(r)
	di := make([]*Entry, 0, max>>8)
	// 当前字节
	for curr := _u32; curr < max; {
		head := ReadN(r, 4)
		// 词长 * 2
		wordLen := head[0]%0x80 - 1
		// 拼音长
		pyLen := (head[1]&0xF)*2 + head[0]/0x80
		// 频率
		freq := BytesToInt(head[2:])
		curr += uint32(4 + wordLen + pyLen*2)
		// 拼音
		pinyin := make([]string, 0, pyLen)
		for i := 0; i < int(pyLen); i++ {
			bsm, _ := r.ReadByte()
			bym, _ := r.ReadByte()
			smIdx := bsm & 0x1F
			ymIdx := (bsm >> 5) + (bym << 3)
			if bym >= 0x10 || smIdx >= 24 || ymIdx >= 34 {
				break
			}
			pinyin = append(pinyin, f.smList[smIdx]+f.ymList[ymIdx])
		}
		// 词
		wordBytes := ReadN(r, wordLen)
		word := util.DecodeMust(wordBytes, f.encoding)
		di = append(di, &Entry{word, pinyin, freq})
	}
	return di
}
