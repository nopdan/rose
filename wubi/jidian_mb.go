package wubi

import (
	"bytes"

	"github.com/nopdan/rose/util"
)

type JidianMb struct{ Template }

func NewJidianMb() *JidianMb {
	f := new(JidianMb)
	f.Name = "极点码表.mb"
	f.ID = "jidian_mb"
	f.Rank = false
	return f
}

func (d *JidianMb) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)
	r.Seek(0x17, 0)
	util.Info(r, 0x11F-0x17, "")

	r.Seek(0x1B620, 0)
	for r.Len() > 3 {
		codeLen, _ := r.ReadByte()
		if codeLen == 0xff {
			r.Seek(1, 1)
			continue
		}
		wordSize, _ := r.ReadByte()
		r.Seek(1, 1)

		// 读编码
		codeBytes := ReadN(r, codeLen)
		code := string(codeBytes)
		// 读词
		wordBytes := ReadN(r, wordSize)
		word := util.DecodeMust(wordBytes, "UTF-16LE")

		di = append(di, &Entry{
			Word: word,
			Code: code,
		})
	}
	return di
}
