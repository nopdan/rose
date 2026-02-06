package sogou_scel

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestSogouScelImport(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/test.scel")
	if _, err := f.Import(src); err != nil {
		t.Fatal(err)
	}
}
