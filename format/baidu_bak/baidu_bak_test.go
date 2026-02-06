package baidu_bak

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestBaiduBakImport(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/百度输入法词库导出_2023_12_28.bin")
	_, err := f.Import(src)
	if err != nil {
		t.Fatal(err)
	}
}
