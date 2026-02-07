package baidu_bdict

import (
	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var utf16 = util.NewEncoding("UTF-16LE")

// 声母列表
var smList = []string{
	"c", "d", "b", "f", "g", "h", "ch", "j", "k", "l", "m", "n",
	"", "p", "q", "r", "s", "t", "sh", "zh", "w", "x", "y", "z",
}

// 韵母列表
var ymList = []string{
	"uang", "iang", "iong", "ang", "eng", "ian", "iao", "ing", "ong",
	"uai", "uan", "ai", "an", "ao", "ei", "en", "er", "ua", "ie", "in", "iu",
	"ou", "ia", "ue", "ui", "un", "uo", "a", "e", "i", "o", "u", "v",
}

// BaiduBdict 百度分类词库.bdict格式
type BaiduBdict struct {
	model.BaseFormat
}

// New 创建百度词库格式处理器
func New() *BaiduBdict {
	return &BaiduBdict{
		BaseFormat: model.BaseFormat{
			ID:          "bdict",
			Name:        "百度分类词库",
			Type:        model.FormatTypePinyin,
			Extension:   ".bdict",
			Description: "百度分类词库二进制格式",
		},
	}
}

// NewBcd 创建百度分类词库（手机）格式处理器
func NewBcd() *BaiduBdict {
	format := New()
	format.ID = "bcd"
	format.Name = "百度分类词库（手机）"
	format.Extension = ".bcd"
	return format
}

func (f *BaiduBdict) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	r.Seek(0x40, 0)
	regularOffset := r.ReadUint32() // 常规词条偏移
	regularSize := r.ReadUint32()   // 常规词条占用空间

	r.Seek(0x50, 0)
	engOffset := r.ReadUint32() // 英文词条偏移
	engSize := r.ReadUint32()   // 英文词条占用空间

	r.Seek(0x60, 0)
	mixedOffset := r.ReadUint32() // 混合词条偏移
	mixedSize := r.ReadUint32()   // 混合词条占用空间

	r.Seek(0x70, 0)
	regularCount := r.ReadUint32() // 常规词条数
	engCount := r.ReadUint32()     // 英文词条数
	mixedCount := r.ReadUint32()   // 混合词条数

	f.Infof("常规词条偏移: 0x%X, 占用空间: 0x%X, 数量: %d\n", regularOffset, regularSize, regularCount)
	f.Infof("英文词条偏移: 0x%X, 占用空间: 0x%X, 数量: %d\n", engOffset, engSize, engCount)
	f.Infof("混合词条偏移: 0x%X, 占用空间: 0x%X, 数量: %d\n", mixedOffset, mixedSize, mixedCount)

	entries := make([]*model.Entry, 0, regularCount+engCount+mixedCount)

	r.Seek(0x90, 0)
	f.Infof("词库名: %s\n", r.ReadStringEnc(0xD0-0x90, utf16))
	f.Infof("词库作者: %s\n", r.ReadStringEnc(0x110-0xD0, utf16))
	f.Infof("示例词: %s\n", r.ReadStringEnc(0x150-0x110, utf16))
	f.Infof("词库描述: %s\n", r.ReadStringEnc(0x350-0x150, utf16))

	// 读取常规词条
	r.Seek(int64(regularOffset), 0)
	for range regularCount {
		pyLen := r.ReadIntN(2) // 拼音长度
		freq := r.ReadIntN(2)  // 词频
		pinyin := make([]string, pyLen)
		for i := range pyLen {
			smIdx, _ := r.ReadByte()
			ymIdx, _ := r.ReadByte()
			// 带英文的词组
			if smIdx == 0xff {
				pinyin[i] = string(ymIdx)
				continue
			}
			pinyin[i] = smList[smIdx] + ymList[ymIdx]
		}
		// 读词
		word := r.ReadStringEnc(pyLen*2, utf16)
		entry := model.NewEntry(word).WithMultiCode(pinyin...).WithFrequency(freq)
		f.Debugf("%s\t%v\t%d\n", word, pinyin, freq)
		entries = append(entries, entry)
	}

	// 读取英文词条
	r.Seek(int64(engOffset), 0)
	for range engCount {
		size := r.ReadIntN(2) // 长度
		freq := r.ReadIntN(2) // 词频
		word := r.ReadString(size)
		entry := model.NewEntry(word).WithSimpleCode(word).WithFrequency(freq)
		f.Debugf("[英文]: %s\t%d\n", word, freq)
		entries = append(entries, entry)
	}

	// 读取混合词条
	r.Seek(int64(mixedOffset), 0)
	for range mixedCount {
		pyLen := r.ReadIntN(2) // 拼音长度
		freq := r.ReadIntN(2)  // 词频
		r.Seek(2, 1)           // 跳过两个字节
		wordLen := r.ReadIntN(2)
		// 读编码
		code := r.ReadStringEnc(pyLen*2, utf16)
		// 读词
		word := r.ReadStringEnc(wordLen*2, utf16)
		entry := model.NewEntry(word).WithSimpleCode(code).WithFrequency(freq)
		f.Debugf("[混合]: %s\t%s\t%d\n", word, code, freq)
		entries = append(entries, entry)
	}

	f.Infof("成功导入词条数: %d\n", len(entries))
	return entries, nil
}
