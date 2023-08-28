package wubi

import (
	"bytes"
	"fmt"
	"slices"
	"time"

	"github.com/nopdan/rose/util"
)

type MsUDP struct{ Template }

func init() {
	FormatList = append(FormatList, NewMsUDP())
}
func NewMsUDP() *MsUDP {
	f := new(MsUDP)
	f.Name = "微软用户自定义短语.dat"
	f.ID = "udp"
	f.CanMarshal = true
	f.HasRank = true
	return f
}

func (f *MsUDP) Unmarshal(r *bytes.Reader) []*Entry {
	// 词库偏移量
	r.Seek(0x10, 0)
	offset_start := ReadIntN(r, 4) // 偏移表开始
	entry_start := ReadUint32(r)   // 词条开始
	entry_end := ReadUint32(r)     // 词条结束
	count := ReadIntN(r, 4)        // 词条数
	di := make([]*Entry, 0, count)

	export_stamp := ReadUint32(r) // 导出的时间戳
	export_time := time.Unix(int64(export_stamp), 0)
	fmt.Printf("时间: %v\n", export_time)

	_ = entry_end
	for i := 0; i < count; i++ {
		r.Seek(int64(offset_start+4*i), 0)
		offset := ReadUint32(r)
		r.Seek(int64(entry_start+offset), 0)
		r.Seek(6, 1)
		rank := ReadIntN(r, 1)        // 顺序 1-9
		r.ReadByte()                  // 0x06 不明
		r.Seek(4, 1)                  // 4 个空字节
		insert_stamp := ReadUint32(r) // 时间戳
		_ = insert_stamp
		// insert_time := util.MsToTime(insert_stamp)
		// fmt.Println(insert_time)

		codeBytes := make([]byte, 0, 1)
		wordBytes := make([]byte, 0, 1)
		for {
			tmp := ReadN(r, 2)
			if bytes.Equal(tmp, []byte{0, 0}) {
				break
			}
			codeBytes = append(codeBytes, tmp...)
		}
		for {
			tmp := ReadN(r, 2)
			if bytes.Equal(tmp, []byte{0, 0}) {
				break
			}
			wordBytes = append(wordBytes, tmp...)
		}
		code := util.DecodeMust(codeBytes, "UTF-16LE")
		word := util.DecodeMust(wordBytes, "UTF-16LE")
		// util.PrintHex(codeBytes)
		// util.PrintHex(wordBytes)
		di = append(di, &Entry{
			Word: word,
			Code: code,
			Rank: rank,
		})
	}
	return di
}

func (MsUDP) Marshal(di []*Entry, hasRank bool) []byte {
	now := time.Now()
	export_stamp := To4Bytes(now.Unix())
	insert_stamp := util.MsTimeTo(now)

	slices.DeleteFunc(di, func(e *Entry) bool {
		return e.Code == ""
	})
	if !hasRank {
		di = GenRank(di)
	}

	b := make([]byte, 0, len(di))
	b = append(b, 0x6D, 0x73, 0x63, 0x68, 0x78, 0x75, 0x64, 0x70,
		0x02, 0x00, 0x60, 0x00, 0x01, 0x00, 0x00, 0x00)
	b = append(b, To4Bytes(0x40)...)
	b = append(b, To4Bytes(0x40+4*len(di))...)
	b = append(b, 0, 0, 0, 0) // 待定 文件总长
	b = append(b, To4Bytes(len(di))...)
	b = append(b, export_stamp...)
	b = append(b, make([]byte, 28)...)

	// 偏移表，每个词占 4 字节
	b = append(b, make([]byte, 4*len(di))...)
	offset := 0
	for i, v := range di {
		copy(b[0x40+4*i:0x40+4*i+4], To4Bytes(offset))

		b = append(b, 0x10, 0, 0x10, 0)
		wordBytes := util.EncodeMust(v.Word, "UTF-16LE")
		codeBytes := util.EncodeMust(v.Code, "UTF-16LE")
		b = append(b, To2Bytes(len(codeBytes)+18)...)
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
	copy(b[0x18:0x1c], To4Bytes(len(b)))
	return b
}
