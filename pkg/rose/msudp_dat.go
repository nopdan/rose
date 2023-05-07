package rose

import (
	"bytes"
	"fmt"
	"time"

	"github.com/nopdan/ku"
)

type MsUDP struct{ Dict }

func NewMsUDP() *MsUDP {
	d := new(MsUDP)
	d.Name = "微软用户自定义短语.dat"
	d.Suffix = "dat"
	return d
}

func (d *MsUDP) Parse() {
	r := bytes.NewReader(d.data)

	// 词库偏移量
	r.Seek(0x10, 0)
	offset_start := ReadUint32(r) // 偏移表开始
	entry_start := ReadUint32(r)  // 词条开始
	// entry_end := ReadUint32(r)    // 词条结束
	r.Seek(4, 1)
	count := ReadUint32(r) // 词条数
	wl := make([]Entry, 0, count)

	export_stamp := ReadUint32(r) // 导出的时间戳
	export_time := time.Unix(int64(export_stamp), 0)
	fmt.Printf("时间: %v\n", export_time)

	for i := _u32; i < count; i++ {
		r.Seek(int64(offset_start+4*i), 0)
		offset := ReadUint32(r)
		r.Seek(int64(entry_start+offset), 0)
		r.Seek(6, 1)
		pos, _ := r.ReadByte() // 顺序
		r.ReadByte()           // 0x06 不明
		r.Seek(4, 1)           // 4 个空字节
		r.Seek(4, 1)           // 时间戳
		// insert_stamp := ReadUint32(r)
		// insert_time := MspyTime(insert_stamp)
		// fmt.Println(insert_time)
		code := make([]byte, 0, 1)
		word := make([]byte, 0, 1)
		tmp := make([]byte, 2)
	CODE:
		r.Read(tmp)
		if !bytes.Equal(tmp, []byte{0, 0}) {
			code = append(code, tmp...)
			goto CODE
		}
	WORD:
		r.Read(tmp)
		if !bytes.Equal(tmp, []byte{0, 0}) {
			word = append(word, tmp...)
			goto WORD
		}
		c := DecodeMust(code, "UTF-16LE")
		w := DecodeMust(word, "UTF-16LE")
		// fmt.Println(c, w)
		wl = append(wl, &WubiEntry{w, c, int(pos)})
	}
	d.WordLibrary = wl
}

func (MsUDP) GenFrom(wl WordLibrary) []byte {
	var buf bytes.Buffer
	now := time.Now()
	export_stamp := ku.To4Bytes(now.Unix())
	insert_stamp := MspyTimeTo(now)
	buf.Write([]byte{0x6D, 0x73, 0x63, 0x68, 0x78, 0x75, 0x64, 0x70,
		0x02, 0x00, 0x60, 0x00, 0x01, 0x00, 0x00, 0x00})
	buf.Write(ku.To4Bytes(0x40))
	buf.Write(ku.To4Bytes(0x40 + 4*len(wl)))
	buf.Write(make([]byte, 4)) // 待定 文件总长
	buf.Write(ku.To4Bytes(len(wl)))
	buf.Write(export_stamp)
	buf.Write(make([]byte, 28))
	buf.Write(make([]byte, 4))

	words := make([][]byte, 0, len(wl))
	codes := make([][]byte, 0, len(wl))
	sum := 0
	for i := range wl {
		word := EncodeMust(wl[i].GetWord(), "UTF-16LE")
		code := EncodeMust(wl[i].GetCode(), "UTF-16LE")
		words = append(words, word)
		codes = append(codes, code)
		if i != len(wl)-1 {
			sum += len(word) + len(code) + 20
			buf.Write(ku.To4Bytes(sum))
		}
	}
	for i := range wl {
		buf.Write([]byte{0x10, 0x00, 0x10, 0x00})
		// fmt.Println(words[i], len(words[i]), codes[i], len(codes[i]))
		buf.Write(ku.To2Bytes(len(codes[i]) + 18))
		pos := wl[i].GetPos()
		if pos < 1 {
			pos = 1
		}
		buf.WriteByte(byte(pos))
		buf.WriteByte(0x06)
		buf.Write(make([]byte, 4))
		buf.Write(insert_stamp)
		buf.Write(codes[i])
		buf.Write([]byte{0, 0})
		buf.Write(words[i])
		buf.Write([]byte{0, 0})
	}
	b := buf.Bytes()
	copy(b[0x18:0x1c], ku.To4Bytes(len(b)))
	return b
}
