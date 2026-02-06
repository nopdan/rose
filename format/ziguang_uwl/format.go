package ziguang_uwl

import (
	"fmt"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

type ZiguangUwl struct {
	model.BaseFormat
	encoding string
	smList   []string
	ymList   []string
}

func New() *ZiguangUwl {
	f := &ZiguangUwl{
		BaseFormat: model.BaseFormat{
			ID:          "ziguang_uwl",
			Name:        "紫光华宇拼音词库",
			Type:        model.FormatTypePinyin,
			Extension:   ".uwl",
			Description: "紫光华宇拼音词库格式",
		},
		smList: []string{
			"", "b", "c", "ch", "d", "f", "g", "h", "j", "k", "l", "m", "n",
			"p", "q", "r", "s", "sh", "t", "w", "x", "y", "z", "zh",
		},
		ymList: []string{
			"ang", "a", "ai", "an", "ang", "ao", "e", "ei", "en", "eng", "er",
			"i", "ia", "ian", "iang", "iao", "ie", "in", "ing", "iong", "iu",
			"o", "ong", "ou", "u",
			"ua", "uai", "uan", "uang", "ue", "ui", "un", "uo", "v",
		},
	}
	return f
}

func (f *ZiguangUwl) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	entries := make([]*model.Entry, 0, r.Size()>>8)

	r.Seek(2, 0)
	// 编码格式，08 为 GBK，09 为 UTF-16LE
	encoding, _ := r.ReadByte()
	switch encoding {
	case 0x08:
		f.encoding = "GBK"
	case 0x09:
		f.encoding = "UTF-16LE"
	default:
		return nil, fmt.Errorf("ZiguangUwl: Unsupported encoding: %x", encoding)
	}
	enc := util.NewEncoding(f.encoding)

	// 分段
	r.Seek(0x48, 0)
	partLen := r.ReadUint32()
	for i := range partLen {
		r.Seek(0xC00+int64(i)<<10, 0)
		r.Seek(12, 1)
		part, err := f.parse(r, enc)
		if err != nil {
			return nil, err
		}
		entries = append(entries, part...)
	}
	return entries, nil
}

func (f *ZiguangUwl) parse(r *model.Reader, enc *util.Encoding) ([]*model.Entry, error) {
	// 词条占用字节数
	max := r.ReadUint32()
	entries := make([]*model.Entry, 0, max>>8)
	// 当前字节
	for curr := uint32(0); curr < max; {
		head := r.ReadN(4)
		// 词长 * 2
		wordLen := head[0]%0x80 - 1
		// 拼音长
		pyLen := (head[1]&0xF)*2 + head[0]/0x80
		// 频率
		freq := util.Bytes2Int(head[2:])
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
		wordBytes := r.ReadN(int(wordLen))
		word := enc.Decode(wordBytes)

		entry := model.NewEntry(word).
			WithMultiCode(pinyin...).
			WithFrequency(freq)
		if pyLen != wordLen {
			entry.CodeType = model.CodeTypeIncompletePinyin
		}
		entries = append(entries, entry)
		f.Debugf("%s\t%s\t%d\n", word, entry.Code, freq)
	}
	return entries, nil
}
