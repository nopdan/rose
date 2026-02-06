package jidian_mb

import (
	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var utf16 = util.NewEncoding("UTF-16LE")

type JidianMb struct {
	model.BaseFormat
}

func New() *JidianMb {
	return &JidianMb{
		BaseFormat: model.BaseFormat{
			ID:          "jidian_mb",
			Name:        "极点码表",
			Type:        model.FormatTypeWubi,
			Extension:   ".mb",
			Description: "极点五笔码表格式",
		},
	}
}

func (f *JidianMb) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	entries := make([]*model.Entry, 0, r.Size()>>8)
	r.Seek(0x17, 0)
	//!TODO
	// util.Info(r, 0x11F-0x17, "")

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
		codeBytes := r.ReadN(int(codeLen))
		code := string(codeBytes)
		// 读词
		wordBytes := r.ReadN(int(wordSize))
		word := utf16.Decode(wordBytes)

		entries = append(entries, model.NewEntry(word).WithSimpleCode(code))
		f.Debugf("%s\t%s\n", word, code)
	}
	return entries, nil
}
