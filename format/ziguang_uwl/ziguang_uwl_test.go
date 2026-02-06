package ziguang_uwl

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestZiguangUwlImport(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/music.uwl")
	if _, err := f.Import(src); err != nil {
		t.Fatal(err)
	}
}
