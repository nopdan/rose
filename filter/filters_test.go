package filter

import (
	"testing"

	"github.com/nopdan/rose/model"
)

func TestLengthFilter(t *testing.T) {
	entries := []*model.Entry{
		{Word: "一", Code: model.NewSimpleCode("")},
		{Word: "你好", Code: model.NewSimpleCode("")},
		{Word: "测试词汇", Code: model.NewSimpleCode("")},
		{Word: "这是一个很长的词条", Code: model.NewSimpleCode("")},
	}

	filtered := applyFilter(NewLengthFilter(2, 0), entries)
	if len(filtered) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(filtered))
	}

	filtered = applyFilter(NewLengthFilter(0, 4), entries)
	if len(filtered) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(filtered))
	}

	filtered = applyFilter(NewLengthFilter(2, 4), entries)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(filtered))
	}
}

func TestFrequencyFilter(t *testing.T) {
	entries := []*model.Entry{
		{Word: "词1", Code: model.NewSimpleCode(""), Frequency: 10},
		{Word: "词2", Code: model.NewSimpleCode(""), Frequency: 50},
		{Word: "词3", Code: model.NewSimpleCode(""), Frequency: 100},
		{Word: "词4", Code: model.NewSimpleCode(""), Frequency: 0},
	}

	filtered := applyFilter(NewFrequencyFilter(20, 0), entries)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(filtered))
	}

	// 验证过滤结果
	for _, entry := range filtered {
		if entry.Frequency < 20 {
			t.Errorf("Entry with frequency %d should be filtered out", entry.Frequency)
		}
	}

	filtered = applyFilter(NewFrequencyFilter(0, 50), entries)
	if len(filtered) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(filtered))
	}
	for _, entry := range filtered {
		if entry.Frequency > 50 {
			t.Errorf("Entry with frequency %d should be filtered out", entry.Frequency)
		}
	}

	filtered = applyFilter(NewFrequencyFilter(20, 80), entries)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(filtered))
	}
	if len(filtered) == 1 && filtered[0].Frequency != 50 {
		t.Errorf("Expected frequency 50, got %d", filtered[0].Frequency)
	}
}

func TestCharacterFilter(t *testing.T) {
	entries := []*model.Entry{
		{Word: "纯中文", Code: model.NewSimpleCode("")},
		{Word: "English", Code: model.NewSimpleCode("")},
		{Word: "数字123", Code: model.NewSimpleCode("")},
		{Word: "混合abc123", Code: model.NewSimpleCode("")},
	}

	filtered := applyFilter(NewCharacterFilter(true, false), entries)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(filtered))
	}

	filtered = applyFilter(NewCharacterFilter(false, true), entries)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(filtered))
	}

	filtered = applyFilter(NewCharacterFilter(true, true), entries)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(filtered))
	}

	if filtered[0].Word != "纯中文" {
		t.Errorf("Expected '纯中文', got '%s'", filtered[0].Word)
	}
}

func TestRegexFilter(t *testing.T) {
	entries := []*model.Entry{
		{Word: "测试", Code: model.NewSimpleCode("")},
		{Word: "测试词汇", Code: model.NewSimpleCode("")},
		{Word: "其他词条", Code: model.NewSimpleCode("")},
		{Word: "特殊@符号", Code: model.NewSimpleCode("")},
	}

	filtered := applyFilter(NewRegexFilter([]string{"测试"}), entries)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(filtered))
	}

	filtered = applyFilter(NewRegexFilter([]string{"测试", "@"}), entries)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(filtered))
	}

	if filtered[0].Word != "其他词条" {
		t.Errorf("Expected '其他词条', got '%s'", filtered[0].Word)
	}
}

func applyFilter(filter Filter, entries []*model.Entry) []*model.Entry {
	filtered := make([]*model.Entry, 0, len(entries))
	for _, entry := range entries {
		if !filter.Filter(entry) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
