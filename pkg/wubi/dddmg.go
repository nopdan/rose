package wubi

import (
	"bytes"
)

type DDdmg struct{ Template }

func init() {
	FormatList = append(FormatList, NewDDdmg())
}
func NewDDdmg() *DDdmg {
	f := new(DDdmg)
	f.Name = "多多v3.dmg"
	f.ID = "duoduo_dmg,dmg"
	return f
}

func (DDdmg) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)
	r.Seek(0x4089C, 0)
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
		rank := ReadIntN(r, 4)
		_ = rank
		codeLen := ReadIntN(r, 1)
		code := string(ReadN(r, codeLen))
		wordSize := ReadIntN(r, 1)
		word := string(ReadN(r, wordSize))

		di = append(di, &Entry{
			Word: word,
			Code: code,
		})
	}
	return di
}
