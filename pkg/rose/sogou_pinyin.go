package rose

import (
	"bytes"
	"fmt"

	util "github.com/flowerime/goutil"
)

var sg_pinyin = append(mspy, []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p",
	"q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5",
	"6", "7", "8", "9", "#"}...)

type SogouBin struct{ Dict }

func NewSogouBin() *SogouBin {
	d := new(SogouBin)
	d.Name = "搜狗拼音备份.bin"
	d.Suffix = "bin"
	return d
}

func (d *SogouBin) Parse() {
	wl := make([]Entry, 0, d.size>>8)

	r := bytes.NewReader(d.data)
	header := make([]byte, 4) // SGPU
	r.Read(header)
	if !bytes.Equal(header, []byte{0x53, 0x47, 0x50, 0x55}) {
		r.Seek(0, 0)
		d.ParseOld()
		return
	}

	r.Seek(12, 1)
	fileSize := ReadUint32(r) // file total size
	r.Seek(36, 1)
	idxBegin := ReadUint32(r) // index begin
	idxSize := ReadUint32(r)  // index size
	wordCount := ReadUint32(r)
	dictBegin := ReadUint32(r)     // = idxBegin + idxSize
	dictTotalSize := ReadUint32(r) // file total size - dictBegin
	dictSize := ReadUint32(r)      // effective dict size
	fmt.Printf("fileSize: 0x%x\n", fileSize)
	fmt.Printf("idxBegin: 0x%x\n", idxBegin)
	fmt.Printf("idxSize: 0x%x\n", idxSize)
	fmt.Printf("dictBegin: 0x%x\n", dictBegin)
	fmt.Printf("dictTotalSize: 0x%x\n", dictTotalSize)
	fmt.Printf("dictSize: 0x%x\n", dictSize)

	for i := _u32; i < wordCount; i++ {
		r.Seek(int64(idxBegin+4*i), 0)
		idx := ReadUint32(r)
		if idx == 0 && i != 0 {
			break
		}
		r.Seek(int64(idx+dictBegin), 0)
		freq := ReadUint32(r)
		r.Seek(5, 1) // 00 00 FE 07 02
		pyLen := ReadUint16(r) / 2
		py := make([]string, 0, pyLen)
		for j := _u16; j < pyLen; j++ {
			p := ReadUint16(r)
			py = append(py, sg_pinyin[p])
		}
		ReadUint16(r) // word size + code size（include idx）
		wordSize := ReadUint16(r)
		tmp := make([]byte, wordSize)
		r.Read(tmp)
		word, _ := util.Decode(tmp, "UTF-16LE")

		wl = append(wl, &PinyinEntry{word, py, int(freq)})
		// repeat code
	}
	d.WordLibrary = wl
}
