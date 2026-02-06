package kafan_pinyin_bak

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestKafanPinyinBakImport(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	input := "../../testdata/拼音词库备份2023-08-19.dict"
	src := model.NewFileSource(input)
	if _, err := f.Import(src); err != nil {
		t.Fatal(err)
	}
}
