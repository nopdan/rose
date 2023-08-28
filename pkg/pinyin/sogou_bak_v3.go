package pinyin

import (
	"bytes"
	"fmt"

	"github.com/nopdan/rose/util"
)

type SogouBak struct {
	Template
	pyList []string
}

func init() {
	FormatList = append(FormatList, NewSogouBak())
}
func NewSogouBak() *SogouBak {
	f := new(SogouBak)
	f.Name = "搜狗拼音备份.bin"
	f.ID = "sgbak"
	pyList := NewMspyUDL().pyList
	pyList = append(pyList, []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p",
		"q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5",
		"6", "7", "8", "9", "#"}...)
	f.pyList = pyList
	return f
}

func (f *SogouBak) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)
	header := ReadN(r, 4) // SGPU
	if !bytes.Equal(header, []byte{0x53, 0x47, 0x50, 0x55}) {
		r.Seek(0, 0)
		return f.unmarshalV2(r)
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
		r.Seek(int64(idx+dictBegin), 0)
		freq := ReadIntN(r, 2)
		ReadUint16(r) // unknown
		r.Seek(5, 1)  // unknown 5 bytes, same in every entry

		pyLen := ReadIntN(r, 2) / 2
		py := make([]string, 0, pyLen)
		for j := 0; j < pyLen; j++ {
			p := ReadUint16(r)
			py = append(py, f.pyList[p])
		}
		ReadUint16(r) // word size + code size（include idx）
		wordSize := ReadUint16(r)
		wordBytes := ReadN(r, wordSize)
		word := util.DecodeMust(wordBytes, "UTF-16LE")

		di = append(di, &Entry{word, py, freq})
		// repeat code
	}
	return di
}
