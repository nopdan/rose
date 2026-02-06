package qq_qpyd

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestQpydImport(t *testing.T) {
	f := New()
	f.LogLevel = model.LogDebug
	src := model.NewFileSource("../../testdata/网络用语.qpyd")
	if _, err := f.Import(src); err != nil {
		t.Fatal(err)
	}
}
