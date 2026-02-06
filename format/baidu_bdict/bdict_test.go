package baidu_bdict

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestBaiduBdictImport(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/计算机硬件词汇.bdict")
	_, err := f.Import(src)
	if err != nil {
		t.Fatal(err)
	}
}
