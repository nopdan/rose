package custom_text

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

func TestCustomText_BaiduToQQ(t *testing.T) {
	testdataPath := "../../testdata/baidupinyin.txt"
	src := model.NewFileSource(testdataPath)

	importer := NewBaiduPinyin()
	importer.LogLevel = model.LogDebug
	entries, err := importer.Import(src)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no entries imported")
	}

	exporter := NewQQPinyin()
	var buf bytes.Buffer
	if err := exporter.Export(entries, &buf); err != nil {
		t.Fatalf("export failed: %v", err)
	}
	os.WriteFile("test_export.txt", buf.Bytes(), 0o644)

	decoded := util.NewEncoding("UTF-16LE").Decode(buf.Bytes())
	decoded = strings.TrimPrefix(decoded, "\uFEFF")

	expected := "a'bao'zhi'gong' 阿保之功 1"
	if !strings.Contains(decoded, expected) {
		t.Fatalf("export content missing expected line: %s", expected)
	}
}

func TestCustomText_DuoduoToQQ(t *testing.T) {
	testdataPath := "../../testdata/duoduo.txt"
	src := model.NewFileSource(testdataPath)

	importer := NewDuoduoWubi()
	importer.LogLevel = model.LogDebug
	entries, err := importer.Import(src)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no entries imported")
	}

	exporter := NewBaiduWubi()
	var buf bytes.Buffer
	if err := exporter.Export(entries, &buf); err != nil {
		t.Fatalf("export failed: %v", err)
	}
	os.WriteFile("test_export.txt", buf.Bytes(), 0o644)
}
