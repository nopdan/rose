package sogou_scel

import (
	"errors"
	"time"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var utf16 = util.NewEncoding("UTF-16LE")

// SogouScel 搜狗细胞词库格式
type SogouScel struct {
	model.BaseFormat
}

// New 创建搜狗细胞词库格式处理器
func New() *SogouScel {
	return &SogouScel{
		BaseFormat: model.BaseFormat{
			ID:          "sogou_scel",
			Name:        "搜狗细胞词库.scel",
			Type:        model.FormatTypePinyin,
			Extension:   ".scel",
			Description: "搜狗细胞词库二进制格式",
		},
	}
}

// NewQcel 创建QQ拼音细胞词库格式处理器
func NewQcel() *SogouScel {
	format := New()
	format.ID = "qq_qcel"
	format.Name = "QQ拼音v6以上.qcel"
	format.Extension = ".qcel"
	return format
}

const (
	PHRASE_OFFSET    = 0x5C  // 短语数量+偏移
	DELTBL_OFFSET    = 0x74  // 黑名单数量+偏移
	TIMESTAMP_OFFSET = 0x11C // 时间戳
	ENTRY_OFFSET     = 0x120 // 词条数量
	NAME_OFFSET      = 0x130 // 词库名称
	LOCATION_OFFSET  = 0x338 // 地点
	DESC_OFFSET      = 0x540 // 备注
	EXAMPLE_OFFSET   = 0xD40 // 示例词
)

func (f *SogouScel) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	entryOffset := r.ReadUint32()
	magic := r.ReadUint32()
	f.Infof("魔数: %X\n", magic)

	r.Seek(0xC, 0)
	chksum := r.ReadN(16)
	f.Infof("校验和: %X\n", chksum)

	id := r.ReadStringEnc(16, utf16)
	f.Infof("ID: %s\n", id)

	r.Seek(PHRASE_OFFSET, 0)
	phraseCount := r.ReadUint32()
	phraseOffset := r.ReadUint32()
	phraseSize := r.ReadUint32()
	r.Seek(4, 1)
	phraseEndOffset := r.ReadUint32()
	if phraseEndOffset != phraseOffset+phraseSize {
		return nil, errors.New("phrase offset and size mismatch")
	}
	f.Infof("短语数量: %d, 偏移: 0x%X, 数据大小: %d\n", phraseCount, phraseOffset, phraseSize)

	r.Seek(DELTBL_OFFSET, 0)
	deltblCount := r.ReadUint32()
	deltblOffset := r.ReadUint32()
	deltblSize := r.ReadUint32()
	f.Infof("黑名单数量: %d, 偏移: 0x%X, 数据大小: %d\n", deltblCount, deltblOffset, deltblSize)

	r.Seek(TIMESTAMP_OFFSET, 0)
	timestamp := r.ReadUint32()
	t := time.Unix(int64(timestamp), 0)
	f.Infof("时间: %s\n", t.Format(time.RFC3339))

	r.Seek(ENTRY_OFFSET, 0)
	codeCount := r.ReadUint32()
	wordCount := r.ReadUint32()
	codeSize := r.ReadUint32()
	wordSize := r.ReadUint32()
	f.Infof("编码数量: %d, 数据大小: %d, 词语数量: %d, 数据大小: %d\n", codeCount, codeSize, wordCount, wordSize)

	r.Seek(NAME_OFFSET, 0)
	name := r.ReadStringEnc(LOCATION_OFFSET-NAME_OFFSET, utf16)
	f.Infof("词库名称: %s\n", name)

	r.Seek(LOCATION_OFFSET, 0)
	location := r.ReadStringEnc(DESC_OFFSET-LOCATION_OFFSET, utf16)
	f.Infof("地点: %s\n", location)

	r.Seek(DESC_OFFSET, 0)
	desc := r.ReadStringEnc(EXAMPLE_OFFSET-DESC_OFFSET, utf16)
	f.Infof("备注: %s\n", desc)

	r.Seek(EXAMPLE_OFFSET, 0)
	example := r.ReadStringEnc(int(entryOffset)-EXAMPLE_OFFSET, utf16)
	f.Infof("示例词: %s\n", example)

	r.Seek(int64(entryOffset), 0)
	pyListLen := r.ReadUint32()
	f.Infof("拼音列表长度: %d\n", pyListLen)
	pys := make([]string, pyListLen)
	// 读取拼音列表
	for range pyListLen {
		index := r.ReadUint16()
		length := r.ReadUint16()
		pys[index] = r.ReadStringEnc(int(length), utf16)
	}
	if pyListLen == 0 {
		f.Infof("拼音列表为空，使用默认拼音\n")
		pys = pyList
	}

	entries := make([]*model.Entry, 0, wordCount)
	// 读取词条
	for range codeCount {
		// 重码数
		repeat := r.ReadUint16()
		// 拼音索引长度
		pyLen := r.ReadUint16() / 2
		pinyin := make([]string, 0, pyLen)
		for range pyLen {
			idx := r.ReadUint16()
			pinyin = append(pinyin, indexPinyin(int(idx), pys))
		}
		for range repeat {
			wordSize := r.ReadUint16()
			word := r.ReadStringEnc(int(wordSize), utf16)

			extSize := r.ReadUint16()
			freq := r.ReadUint32()
			r.Seek(int64(extSize-4), 1) // 跳过扩展信息

			entry := model.NewEntry(word).WithMultiCode(pinyin...).WithFrequency(int(freq))
			entries = append(entries, entry)
			f.Debugf("%s\t%s\t%d\n", word, pinyin, freq)
		}
	}

	// 读取短语
	r.Seek(int64(phraseOffset), 0)
	for range phraseCount {
		info := r.ReadN(17)

		codeSize := r.ReadUint16()
		if info[2] == 0x01 {
			// 拼音
			pyLen := codeSize / 2
			pinyin := make([]string, pyLen)
			for i := range pyLen {
				index := r.ReadUint16()
				pinyin[i] = indexPinyin(int(index), pys)
			}
			wordSize := r.ReadUint16()
			word := r.ReadStringEnc(int(wordSize), utf16)
			entry := model.NewEntry(word).WithMultiCode(pinyin...)
			entries = append(entries, entry)
			f.Debugf("[拼音短语]: %s\t%v\n", word, pinyin)
		} else {
			// 自定义
			code := r.ReadStringEnc(int(codeSize), utf16)
			wordSize := r.ReadUint16()
			word := r.ReadStringEnc(int(wordSize), utf16)
			entry := model.NewEntry(word).WithSimpleCode(code)
			entries = append(entries, entry)
			f.Debugf("[自定义短语]: %s\t%s\n", word, code)
		}
	}

	// 读取黑名单
	r.Seek(int64(deltblOffset), 0)
	for range deltblCount {
		size := r.ReadUint16() * 2
		word := r.ReadStringEnc(int(size), utf16)
		f.Debugf("[黑名单]: %s\n", word)
	}

	return entries, nil
}
