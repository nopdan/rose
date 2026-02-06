package fcitx4_mb

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestFcitx4MbImport(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/98wb_ci.mb")
	if _, err := f.Import(src); err != nil {
		t.Fatal(err)
	}
}
