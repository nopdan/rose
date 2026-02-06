package msudp

import (
	"bytes"
	"io"
	"slices"
	"time"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var utf16 = util.NewEncoding("UTF-16LE")

type MsUDP struct {
	model.BaseFormat
}

func New() *MsUDP {
	return &MsUDP{
		BaseFormat: model.BaseFormat{
			ID:          "msudp",
			Name:        "微软用户自定义短语",
			Type:        model.FormatTypeWubi,
			Extension:   ".dat",
			Description: "微软输入法用户自定义短语格式",
		},
	}
}

func (f *MsUDP) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	// 词库偏移量
	r.Seek(0x10, 0)
	offset_start := r.ReadIntN(4) // 偏移表开始
	entry_start := r.ReadUint32() // 词条开始
	entry_end := r.ReadUint32()   // 词条结束
	count := r.ReadIntN(4)        // 词条数
	entries := make([]*model.Entry, 0, count)

	export_stamp := r.ReadUint32() // 导出的时间戳
	export_time := time.Unix(int64(export_stamp), 0)
	f.Infof("时间: %v\n", export_time)

	_ = entry_end
	for i := range count {
		r.Seek(int64(offset_start+4*i), 0)
		offset := r.ReadUint32()
		r.Seek(int64(entry_start+offset), 0)
		r.Seek(6, 1)
		rank := r.ReadIntN(1)          // 顺序 1-9
		r.ReadByte()                   // 0x06 不明
		r.Seek(4, 1)                   // 4 个空字节
		insert_stamp := r.ReadUint32() // 时间戳
		_ = insert_stamp
		// insert_time := util.MsToTime(insert_stamp)
		// fmt.Println(insert_time)

		codeBytes := make([]byte, 0, 1)
		wordBytes := make([]byte, 0, 1)
		for {
			tmp := r.ReadN(2)
			if bytes.Equal(tmp, []byte{0, 0}) {
				break
			}
			codeBytes = append(codeBytes, tmp...)
		}
		for {
			tmp := r.ReadN(2)
			if bytes.Equal(tmp, []byte{0, 0}) {
				break
			}
			wordBytes = append(wordBytes, tmp...)
		}
		code := utf16.Decode(codeBytes)
		word := utf16.Decode(wordBytes)

		entries = append(entries, model.NewEntry(word).WithSimpleCode(code).WithRank(rank))
		f.Debugf("%s\t%s\t%d\n", word, code, rank)
	}
	return entries, nil
}

func (MsUDP) Export(di []*model.Entry, w io.Writer) error {
	now := time.Now()
	export_stamp := util.To4Bytes(now.Unix())
	insert_stamp := MsTimeTo(now)

	di = slices.DeleteFunc(di, func(e *model.Entry) bool {
		return e.Code.String() == ""
	})
	// if !hasRank {
	// 	di = GenRank(di)
	// }

	b := make([]byte, 0, len(di))
	b = append(b, 0x6D, 0x73, 0x63, 0x68, 0x78, 0x75, 0x64, 0x70,
		0x02, 0x00, 0x60, 0x00, 0x01, 0x00, 0x00, 0x00)
	b = append(b, util.To4Bytes(0x40)...)
	b = append(b, util.To4Bytes(0x40+4*len(di))...)
	b = append(b, 0, 0, 0, 0) // 待定 文件总长
	b = append(b, util.To4Bytes(len(di))...)
	b = append(b, export_stamp...)
	b = append(b, make([]byte, 28)...)

	// 偏移表，每个词占 4 字节
	b = append(b, make([]byte, 4*len(di))...)
	offset := 0
	for i, v := range di {
		copy(b[0x40+4*i:0x40+4*i+4], util.To4Bytes(offset))

		b = append(b, 0x10, 0, 0x10, 0)
		wordBytes := utf16.Encode(v.Word)
		codeBytes := utf16.Encode(v.Code.String())
		b = append(b, util.To2Bytes(len(codeBytes)+18)...)
		rank := v.Rank
		if rank < 1 {
			rank = 1
		}
		b = append(b, byte(rank))
		b = append(b, 0x06, 0, 0, 0, 0)
		b = append(b, insert_stamp...)
		b = append(b, codeBytes...)
		b = append(b, 0, 0)
		b = append(b, wordBytes...)
		b = append(b, 0, 0)

		offset += len(wordBytes) + len(codeBytes) + 20
	}
	copy(b[0x18:0x1c], util.To4Bytes(len(b)))
	_, err := w.Write(b)
	return err
}

func MsToTime(stamp uint32) time.Time {
	return time.Unix(int64(stamp), 0).Add(946684800 * time.Second)
}

func MsTimeTo(t time.Time) []byte {
	return util.To4Bytes(t.Add(-946684800 * time.Second).Unix())
}
