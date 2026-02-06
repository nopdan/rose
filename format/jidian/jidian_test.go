package jidian

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nopdan/rose/model"
)

func TestJidian(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	inputPath := filepath.Join(".", "test_import.txt")
	content := "abcd 你好 世界\nefgh 测试\n"
	if err := os.WriteFile(inputPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
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

	if err := f.Export(entries, out); err != nil {
		t.Fatal(err)
	}
}
