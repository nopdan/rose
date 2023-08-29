package wubi

import (
	"bytes"
	"io"

	"github.com/nopdan/rose/pkg/util"
)

type Kafan struct {
	Template
}

func init() {
	FormatList = append(FormatList, NewKafan())
}

func NewKafan() *Kafan {
	f := new(Kafan)
	f.Name = "卡饭五笔备份.dict"
	f.ID = "kfwbbak"
	return f
}

func (f *Kafan) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, 0xff)
	for r.Len() > 8 {
		check := make([]byte, 8)
		r.Read(check)
		if string(check) == "ProtoDic" {
			r.Seek(8, io.SeekCurrent)
			break
		}
	}
	for r.Len() > 0x28 {
		tmp := ReadN(r, 4)
		// wubi86
		if bytes.Equal(tmp, []byte{0x77, 0x75, 0x62, 0x69}) {
			r.Seek(4, io.SeekCurrent)
			continue
		}
		if util.BytesToInt(tmp) == 0 {
			continue
		}
		r.Seek(-4, io.SeekCurrent)
		codeBytes := make([]byte, 0, 2)
		var word string
		for {
			tmp := ReadN[int](r, 0x28) // 40
			if bytes.Equal(tmp[:8], []byte{4, 0, 0, 0, 3, 0, 1, 0}) {
				r.Seek(8, io.SeekCurrent) // 未知
				wordBytes := make([]byte, 0, 4)
				for {
					b := ReadN[int](r, 4)
					wordBytes = append(wordBytes, b...)
					if b[3] == 0 {
						break
					}
				}
				// 去除末尾的 0
				for i := len(wordBytes) - 1; i >= 0 && wordBytes[i] == 0; i-- {
					wordBytes = wordBytes[:i]
				}
				word = string(wordBytes)
				break
			}
			codeBytes = append(codeBytes, tmp[0])
		}
		di = append(di, &Entry{
			Word: word,
			Code: string(codeBytes),
		})
	}
	return di
}
