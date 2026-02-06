package mspy_udl

import (
	"bytes"
	"io"
	"slices"
	"time"
	"unicode"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var utf16 = util.NewEncoding("UTF-16LE")

type MspyUDL struct {
	model.BaseFormat
	pyList []string
	pyMap  map[string]uint16
}

func New() *MspyUDL {
	f := &MspyUDL{
		BaseFormat: model.BaseFormat{
			ID:          "mspy_udl",
			Name:        "微软拼音自学习词汇",
			Type:        model.FormatTypePinyin,
			Extension:   ".dat",
			Description: "微软拼音自学习词汇格式",
		},
	}

	f.pyList = pyList
	f.pyMap = make(map[string]uint16)
	for idx, py := range f.pyList {
		f.pyMap[py] = uint16(idx)
	}
	return f
}

// 自学习词库，纯汉字
func (f *MspyUDL) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	r.Seek(0xC, 0)
	count := r.ReadIntN(4)
	entries := make([]*model.Entry, 0, count)

	r.Seek(4, 1)
	export_stamp := r.ReadUint32()
	export_time := MsToTime(export_stamp)
	f.Infof("时间: %v\n", export_time)

	for i := range count {
		r.Seek(0x2400+60*int64(i), 0)
		data := r.ReadN(60)
		// insert_stamp := util.BytesToInt(data[:4])
		// insert_time := MspyTime(uint32(insert_stamp))
		// jianpin := data[4:7]
		wordLen := int(data[10])
		p := 12 + wordLen*2
		if wordLen > 12 {
			f.Infof("词长超限: %d\n", wordLen)
			continue
		}
		wordBytes := data[12:p]
		word := utf16.Decode(wordBytes)

		py := make([]string, 0, wordLen)
		for j := range wordLen {
			idx := util.Bytes2Int(data[p+2*j : p+2*(j+1)])
			if idx < len(f.pyList) {
				py = append(py, f.pyList[idx])
			} else {
				f.Infof("拼音索引超出范围: %d > %d\n", idx, len(f.pyList))
			}
		}

		entry := model.NewEntry(word).
			WithMultiCode(py...).
			WithFrequency(1)
		entries = append(entries, entry)
		f.Debugf("%s\t%s\t%d\n", word, entry.Code, 1)
	}
	return entries, nil
}

func (f *MspyUDL) Export(entries []*model.Entry, w io.Writer) error {
	// 过滤只保留纯汉字词条
	validEntries := make([]*model.Entry, 0)
	for _, v := range entries {
		valid := true
		for _, r := range v.Word {
			if !unicode.Is(unicode.Han, r) {
				valid = false
				break
			}
		}
		if valid && len([]rune(v.Word)) <= 12 {
			validEntries = append(validEntries, v)
		}
	}

	// 转换为内部结构
	type udlEntry struct {
		Word    string
		Pinyin  []string
		Jianpin []byte
	}

	dict := make([]*udlEntry, 0, len(validEntries))
	now := time.Now()
	timeBytes := MsTimeTo(now)

	for _, entry := range validEntries {
		codes := entry.Code.Strings()
		if len(codes) != len([]rune(entry.Word)) {
			continue
		}

		jianpin := f.jianpin(codes)
		dict = append(dict, &udlEntry{
			Word:    entry.Word,
			Pinyin:  codes,
			Jianpin: jianpin,
		})
	}

	// 按照简拼排序
	slices.SortStableFunc(dict, func(a, b *udlEntry) int {
		return bytes.Compare(a.Jianpin, b.Jianpin)
	})

	count := len(dict)
	var buf bytes.Buffer
	buf.Grow(0x2400 + 60*count)

	// 写入头部
	buf.Write([]byte{0x55, 0xAA, 0x88, 0x81, 0x02, 0x00, 0x60, 0x00, 0x55, 0xAA, 0x55, 0xAA})
	buf.Write(util.To4Bytes(count))
	buf.Write(make([]byte, 4))
	buf.Write(timeBytes)
	buf.Write(make([]byte, 0x2400-0x18))

	// 写入词条
	for _, v := range dict {
		b := make([]byte, 60)
		copy(b[:4], timeBytes)
		copy(b[4:7], v.Jianpin) // 3 bytes jianpin
		copy(b[7:10], []byte{0, 0, 4})

		wordBytes := utf16.Encode(v.Word)
		b[10] = byte(len(wordBytes) / 2)
		b[11] = 0x5A
		copy(b[12:], wordBytes)

		if len(wordBytes)/2 > 12 {
			f.Infof("词组过长，跳过: %s\n", v.Word)
			continue
		}

		// 写入拼音索引
		indexBytes := f.GetIndex(v.Pinyin)
		copy(b[12+len(wordBytes):], indexBytes)

		buf.Write(b)
	}

	// 补齐到0x400的倍数
	size := 0x400 - buf.Len()%0x400
	buf.Write(make([]byte, size))

	b := buf.Bytes()
	b[0xE6C] = 1
	b[0xE98] = 2

	_, err := w.Write(b)
	return err
}

func (f *MspyUDL) GetIndex(py []string) []byte {
	b := make([]byte, 0, len(py)*2)
	for _, v := range py {
		if idx, ok := f.pyMap[v]; ok {
			b = append(b, util.To2Bytes(idx)...)
		}
	}
	return b
}

// 三位简拼
func (f *MspyUDL) jianpin(py []string) []byte {
	ret := make([]byte, 3)
	for i := 0; i < 3; i++ {
		if i >= len(py) {
			break
		}
		if len(py[i]) > 0 {
			ret[i] = py[i][0]
		}
	}
	return ret
}

func MsToTime(stamp uint32) time.Time {
	return time.Unix(int64(stamp), 0).Add(946684800 * time.Second)
}

func MsTimeTo(t time.Time) []byte {
	return util.To4Bytes(uint32(t.Add(-946684800 * time.Second).Unix()))
}
