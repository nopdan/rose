package mswb_lex

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nopdan/rose/model"
)

func TestMswbLex(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	inputPath := "../../testdata/ChsWubi.lex"
	entries, err := f.Import(model.NewFileSource(inputPath))
	if err != nil {
		t.Fatal(err)
	}

	outputPath := "test_export" + filepath.Ext(inputPath)
	out, err := os.Create(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	if err := f.Export(entries, out); err != nil {
		t.Fatal(err)
	}
}
