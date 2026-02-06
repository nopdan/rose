package encoder

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func newTestWubiEncoder(isAABC bool) *WubiEncoder {
	data := []byte("你\twqiy\n好\tvbdt\n世\tavn\n界\tjgah\n测\tiyy\n试\tyaxw\n")
	return NewWubiEncoder("custom", data, isAABC)
}

func TestWubiEncode_SingleAndPair(t *testing.T) {
	enc := newTestWubiEncoder(false)

	e1 := &model.Entry{Word: "你", CodeType: model.CodeTypeNone}
	enc.Encode(e1)
	if e1.CodeType != model.CodeTypeWubi {
		t.Fatalf("expected CodeTypeWubi, got %v", e1.CodeType)
	}
	if got, want := e1.Code.String(), "wqiy"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}

	e2 := &model.Entry{Word: "你好", CodeType: model.CodeTypeNone}
	enc.Encode(e2)
	if got, want := e2.Code.String(), "wqvb"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestWubiEncode_ThreeChar_AABC(t *testing.T) {
	enc := newTestWubiEncoder(true)
	entry := &model.Entry{Word: "你好世", CodeType: model.CodeTypeNone}

	enc.Encode(entry)

	if got, want := entry.Code.String(), "wqva"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestWubiEncode_ThreeChar_NonAABC(t *testing.T) {
	enc := newTestWubiEncoder(false)
	entry := &model.Entry{Word: "你好世", CodeType: model.CodeTypeNone}

	enc.Encode(entry)

	if got, want := entry.Code.String(), "wvav"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestWubiEncode_FourPlus(t *testing.T) {
	enc := newTestWubiEncoder(false)
	entry := &model.Entry{Word: "你好世界", CodeType: model.CodeTypeNone}

	enc.Encode(entry)

	if got, want := entry.Code.String(), "wvaj"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestWubiEncode_UnknownWord(t *testing.T) {
	enc := newTestWubiEncoder(false)
	entry := &model.Entry{Word: "未知", CodeType: model.CodeTypeNone}

	enc.Encode(entry)

	if entry.CodeType != model.CodeTypeNone {
		t.Fatalf("expected CodeTypeNone, got %v", entry.CodeType)
	}
	if entry.Code != nil {
		t.Fatalf("expected Code to remain nil")
	}
}

func TestWubiEncodeBatch(t *testing.T) {
	enc := newTestWubiEncoder(false)
	entries := []*model.Entry{
		{Word: "你好", CodeType: model.CodeTypeNone},
		{Word: "世界", CodeType: model.CodeTypeNone},
	}

	enc.EncodeBatch(entries)

	for i, entry := range entries {
		if entry.CodeType != model.CodeTypeWubi {
			t.Fatalf("entry %d expected CodeTypeWubi, got %v", i, entry.CodeType)
		}
		if entry.Code == nil || entry.Code.IsEmpty() {
			t.Fatalf("entry %d expected Code to be set", i)
		}
	}
}

func TestWubiEncode_Schema86(t *testing.T) {
	enc := NewWubiEncoder("wubi86", nil, false)

	e1 := &model.Entry{Word: "一", CodeType: model.CodeTypeNone}
	enc.Encode(e1)
	if got, want := e1.Code.String(), "ggll"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}

	e2 := &model.Entry{Word: "一丁", CodeType: model.CodeTypeNone}
	enc.Encode(e2)
	if got, want := e2.Code.String(), "ggsg"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
