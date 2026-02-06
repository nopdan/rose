package jidian_mb

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestJidianMbImport(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/jidian.mb")
	if _, err := f.Import(src); err != nil {
		t.Fatal(err)
	}
}
