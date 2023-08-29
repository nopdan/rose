package wubi

import (
	"bytes"
)

type DuoDB struct{ Template }

func init() {
	FormatList = append(FormatList, NewDuoDB())
}
func NewDuoDB() *DuoDB {
	f := new(DuoDB)
	f.Name = "多多v4.duodb"
	f.ID = "duoduo_duodb,duodb"
	return f
}

func (DuoDB) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)

	r.Seek(0x4086C, 0)
	offsetList := make([]uint32, 0, 12)
	for {
		offset := ReadUint32(r)
		if offset == 0 {
			break
		}
		offsetList = append(offsetList, offset)
	}
	for _, offset := range offsetList {
		r.Seek(int64(offset), 0)
		r.Seek(4, 1)
		codeLen := ReadIntN(r, 1)
		code := string(ReadN(r, codeLen))
		wordSize := ReadIntN(r, 2)
		word := string(ReadN(r, wordSize))

		di = append(di, &Entry{
			Word: word,
			Code: code,
		})
	}
	return di
}
