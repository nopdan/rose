package rime_table

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestRimeTablePinyin(t *testing.T) {
	f := NewRimeTablePinyin()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/pinyin_simp.table.bin")
	_, err := f.Import(src)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRimeTableWubi(t *testing.T) {
	f := NewRimeTableWubi()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/wb2023.table.bin")
	_, err := f.Import(src)
	if err != nil {
		t.Fatal(err)
	}
}
