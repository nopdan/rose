package rime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nopdan/rose/model"
)

func TestRime(t *testing.T) {
	inputPath := filepath.Join("..", "..", "testdata", "rime_sample.dict.yaml")
	content := "# Rime dictionary\n---\nname: test\n...\n\n你好\tni hao\n"
	if err := os.WriteFile(inputPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	src := model.NewFileSource(inputPath)
	entries, err := (&Format{}).Import(src)
	if err != nil {
		t.Fatal(err)
	}

	base := filepath.Base(inputPath)
	idx := strings.Index(base, ".")
	suffix := ""
	if idx >= 0 {
		suffix = base[idx:]
	}
	outputPath := "test_export" + suffix
	out, err := os.Create(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	if err := (&Format{}).Export(entries, out); err != nil {
		t.Fatal(err)
	}
}
