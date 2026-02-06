package sogou_bak

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestSogouBakImport(t *testing.T) {
	f := NewSogouBak()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/sogou_bak_v3.bin")
	if _, err := f.Import(src); err != nil {
		t.Fatal(err)
	}
}

func TestSogouBakImportV2(t *testing.T) {
	f := NewSogouBak()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/sogou_bak_v2.bin")
	if _, err := f.Import(src); err != nil {
		t.Fatal(err)
	}
}
