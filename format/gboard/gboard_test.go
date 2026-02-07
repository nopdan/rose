package gboard

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nopdan/rose/model"
)

func TestGboard(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	inputPath := "../../testdata/Gboard词库.zip"
	src := model.NewFileSource(inputPath)
	entries, err := f.Import(src)
	if err != nil {
		t.Fatal(err)
	}

	outputPath := "test_export" + filepath.Ext(inputPath)
	out, err := os.Create(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	if err := New().Export(entries, out); err != nil {
		t.Fatal(err)
	}
}
