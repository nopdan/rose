package pinyinjiajia

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nopdan/rose/model"
)

func TestPinyinJiaJia(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	inputPath := filepath.Join("..", "..", "testdata", "jiajia.txt")
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
