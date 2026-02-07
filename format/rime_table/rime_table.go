package rime_table

import (
	"github.com/nopdan/rose/format/custom_text"
	"github.com/nopdan/rose/model"
)

// RimeTableWubi
type RimeTableWubi struct {
	model.BaseFormat
}

func NewRimeTableWubi() *RimeTableWubi {
	return &RimeTableWubi{
		BaseFormat: model.BaseFormat{
			ID:          "rime_table_wubi",
			Name:        "Rime Table 五笔",
			Type:        model.FormatTypeWubi,
			Extension:   ".table.bin",
			Description: "Rime::Table/4.0",
		},
	}
}

func (f *RimeTableWubi) Import(src model.Source) ([]*model.Entry, error) {
	b, err := src.Bytes()
	if err != nil {
		return nil, err
	}
	b, err = Decompile(b)
	if err != nil {
		return nil, err
	}
	src = model.NewBytesSource(src.Name(), b)
	f.Debugf("%s\n", string(b))
	rime := custom_text.NewRimeWubi()
	return rime.Import(src)
}

// RimeTablePinyin
type RimeTablePinyin struct {
	model.BaseFormat
}

func NewRimeTablePinyin() *RimeTablePinyin {
	return &RimeTablePinyin{
		BaseFormat: model.BaseFormat{
			ID:          "rime_table_pinyin",
			Name:        "Rime Table 拼音",
			Type:        model.FormatTypePinyin,
			Extension:   ".table.bin",
			Description: "Rime::Table/4.0",
		},
	}
}

func (f *RimeTablePinyin) Import(src model.Source) ([]*model.Entry, error) {
	b, err := src.Bytes()
	if err != nil {
		return nil, err
	}
	b, err = Decompile(b)
	if err != nil {
		return nil, err
	}
	src = model.NewBytesSource(src.Name(), b)
	f.Debugf("%s\n", string(b))
	rime := custom_text.NewRimePinyin()
	return rime.Import(src)
}
