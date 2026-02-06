package converter

import (
	"testing"

	"github.com/nopdan/rose/encoder"
	"github.com/nopdan/rose/filter"
	"github.com/nopdan/rose/model"
)

func TestConverter(t *testing.T) {
	conv := NewConverter()

	// 测试转换器初始化
	if conv.registry == nil {
		t.Fatal("Registry should not be nil")
	}

	// 测试添加过滤器
	conv.AddFilter(filter.NewLengthFilter(1, 10))
	conv.AddFilter(filter.NewFrequencyFilter(1))

	if len(conv.filters) != 2 {
		t.Errorf("Expected 2 filters, got %d", len(conv.filters))
	}

	// 测试注册编码器
	pinyinEncoder := encoder.NewPinyinEncoder()
	conv.RegisterEncoder(pinyinEncoder)

	if conv.encoder == nil {
		t.Errorf("Expected encoder to be set")
	}

	// 测试获取支持的格式
	formats := conv.GetSupportedFormats()
	if len(formats) == 0 {
		t.Error("Should have at least one supported format")
	}

	t.Logf("Supported formats: %d", len(formats))
	for _, format := range formats {
		t.Logf("- %s: %s", format.ID, format.Name)
	}
}

func TestConverterApplyFilters(t *testing.T) {
	conv := NewConverter()

	// 添加过滤器
	conv.AddFilter(filter.NewLengthFilter(2, 4))
	conv.AddFilter(filter.NewCharacterFilter(true, false))

	// 创建测试数据
	entries := []*model.Entry{
		{Word: "一", Code: model.NewSimpleCode("")},
		{Word: "你好", Code: model.NewSimpleCode("")},
		{Word: "English", Code: model.NewSimpleCode("")},
		{Word: "测试词汇", Code: model.NewSimpleCode("")},
	}

	// 应用过滤器
	filtered := conv.applyFilters(entries)

	// 应该过滤掉"一"（长度不够）和"English"（包含英文）
	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries after filtering, got %d", len(filtered))
	}

	// 验证剩余的词条
	expectedWords := map[string]bool{"你好": true, "测试词汇": true}
	for _, entry := range filtered {
		if !expectedWords[entry.Word] {
			t.Errorf("Unexpected word after filtering: %s", entry.Word)
		}
	}
}

func TestConverterIntegration(t *testing.T) {
	conv := NewConverter()

	// 添加过滤器
	conv.AddFilter(filter.NewLengthFilter(1, 10))

	// 添加编码器
	pinyinEncoder := encoder.NewPinyinEncoder()
	conv.RegisterEncoder(pinyinEncoder)

	// 覆盖编码器
	wubiEncoder := encoder.NewWubiEncoder("wubi86", nil, false)
	conv.RegisterEncoder(wubiEncoder)

	// 测试编码器注册
	if conv.encoder != wubiEncoder {
		t.Error("Expected encoder to be replaced by the latest registration")
	}

	t.Log("Integration test passed - converter can handle filters and encoder")
}
