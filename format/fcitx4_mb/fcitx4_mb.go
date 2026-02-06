package fcitx4_mb

import (
	"github.com/nopdan/rose/model"
)

type Fcitx4Mb struct {
	model.BaseFormat
}

func New() *Fcitx4Mb {
	return &Fcitx4Mb{
		BaseFormat: model.BaseFormat{
			ID:          "fcitx4_mb",
			Name:        "Fcitx4码表",
			Type:        model.FormatTypeWubi,
			Extension:   ".mb",
			Description: "Fcitx4输入法码表格式",
		},
	}
}

func (f *Fcitx4Mb) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	r.Seek(0x55, 0)
	// 词条数
	count := r.ReadIntN(4)
	entries := make([]*model.Entry, 0, count)

	for range count {
		codeBytes := r.ReadN(5)
		code := trimSufZero(codeBytes)

		wordSize := r.ReadIntN(4) - 1
		wordBytes := r.ReadN(wordSize)
		word := string(wordBytes)

		entries = append(entries, model.NewEntry(word).WithSimpleCode(code))
		f.Debugf("%s\t%s\n", word, code)
		r.Seek(10, 1)
	}
	return entries, nil
}

// 去掉末尾的 0
func trimSufZero(b []byte) string {
	for i := len(b); i > 0; i-- {
		if b[i-1] != 0 {
			return string(b[:i])
		}
	}
	return ""
}
