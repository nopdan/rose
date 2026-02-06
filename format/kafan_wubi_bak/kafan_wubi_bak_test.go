package kafan_wubi_bak

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestKafanWubiBakImport(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/五笔词库备份2023-08-29.dict")
	if _, err := f.Import(src); err != nil {
		t.Fatal(err)
	}
}
