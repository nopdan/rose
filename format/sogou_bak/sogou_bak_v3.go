package sogou_bak

import (
	"bytes"

	"github.com/nopdan/rose/model"
)

type SogouBak struct {
	model.BaseFormat
}

func NewSogouBak() *SogouBak {
	f := &SogouBak{
		BaseFormat: model.BaseFormat{
			ID:          "sogou_bak",
			Name:        "搜狗拼音备份",
			Type:        model.FormatTypePinyin,
			Extension:   ".bin",
			Description: "搜狗拼音备份词库格式",
		},
	}
	return f
}

func (f *SogouBak) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	entries := make([]*model.Entry, 0, r.Size()>>8)

	header := r.ReadN(4) // SGPU
	if !bytes.Equal(header, []byte{0x53, 0x47, 0x50, 0x55}) {
		return f.importV2(r)
	}

	f.Infof("检测到v3格式文件头: SGPU\n")

	r.Seek(12, 1)

	fileSize := r.ReadUint32() // file total size
	f.Infof("文件总大小: 0x%08X (%d 字节)\n", fileSize, fileSize)

	r.Seek(36, 1)

	idxBegin := r.ReadUint32() // index begin
	f.Infof("索引区起始位置: 0x%08X\n", idxBegin)

	idxSize := r.ReadUint32() // index size
	f.Infof("索引区大小: 0x%08X (%d 字节)\n", idxSize, idxSize)

	wordCount := r.ReadUint32()
	f.Infof("词条数量: %d\n", wordCount)

	dictBegin := r.ReadUint32() // = idxBegin + idxSize
	f.Infof("词典区起始位置: 0x%08X\n", dictBegin)

	dictTotalSize := r.ReadUint32() // file total size - dictBegin
	f.Infof("词典区总大小: 0x%08X (%d 字节)\n", dictTotalSize, dictTotalSize)

	dictSize := r.ReadUint32() // effective dict size
	f.Infof("词典区有效大小: 0x%08X (%d 字节)\n", dictSize, dictSize)

	for i := range wordCount {
		r.Seek(int64(idxBegin+4*i), 0)
		idx := r.ReadUint32()
		r.Seek(int64(idx+dictBegin), 0)
		freq := r.ReadIntN(2)
		r.ReadUint16() // unknown
		r.Seek(5, 1)   // unknown 5 bytes, same in every entry

		pyLen := r.ReadUint16() / 2
		py := make([]string, 0, pyLen)
		for range pyLen {
			p := r.ReadUint16()
			if int(p) < len(pyList) {
				py = append(py, pyList[p])
			}
		}
		r.ReadUint16() // word size + code size（include idx）
		wordSize := r.ReadUint16()
		word := r.ReadStringEnc(int(wordSize), utf16)

		entry := model.NewEntry(word).WithMultiCode(py...).WithFrequency(freq)
		entries = append(entries, entry)

		f.Debugf("%s\t%v\t%d\n", word, py, freq)
	}
	return entries, nil
}
