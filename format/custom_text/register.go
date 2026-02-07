package custom_text

import (
	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var (
	utf16 = util.NewEncoding("UTF-16LE")
	utf8  = util.NewEncoding("UTF-8")
	gb    = util.NewEncoding("GB18030")
)

// ------------- 纯词组模板 -------------

// NewWords 纯词组格式: 每行一个词组
func NewWords() *CustomText {
	return NewCustom(
		"words",
		"纯词组",
		model.FormatTypeWords,
		utf8,
		[]FieldConfig{
			{Type: FieldTypeWord},
		},
	).WithCommentPrefix("#")
}

// ------------- 拼音模板 -------------

// NewSogouPinyin 搜狗拼音格式: 词组
func NewSogouPinyin() *CustomText {
	return NewCustom(
		"sogou",
		"搜狗拼音",
		model.FormatTypeWords,
		gb,
		[]FieldConfig{
			{Type: FieldTypeWord},
		},
	)
}

// NewBaiduPinyin 百度拼音格式: 词组\tpin'yin'\t频率
func NewBaiduPinyin() *CustomText {
	return NewCustom(
		"baidu",
		"百度拼音",
		model.FormatTypePinyin,
		utf16,
		[]FieldConfig{
			{Type: FieldTypeWord},
			{Type: FieldTypeTab},
			{Type: FieldTypePinyin, PinyinSeparator: "'", PinyinSuffix: "'"},
			{Type: FieldTypeTab},
			{Type: FieldTypeFrequency},
		},
	)
}

// NewQQPinyin QQ拼音格式: pin'yin 词组 freq
func NewQQPinyin() *CustomText {
	return NewCustom(
		"qq",
		"QQ拼音",
		model.FormatTypePinyin,
		utf16,
		[]FieldConfig{
			{Type: FieldTypePinyin, PinyinSeparator: "'"},
			{Type: FieldTypeSpace},
			{Type: FieldTypeWord},
			{Type: FieldTypeSpace},
			{Type: FieldTypeFrequency},
		},
	)
}

// NewRimePinyin Rime拼音格式: 词组\tpin yin\t频率
func NewRimePinyin() *CustomText {
	return NewCustom(
		"rime_pinyin",
		"Rime拼音",
		model.FormatTypePinyin,
		utf8,
		[]FieldConfig{
			{Type: FieldTypeWord},
			{Type: FieldTypeTab},
			{Type: FieldTypePinyin, PinyinSeparator: " "},
			{Type: FieldTypeTab},
			{Type: FieldTypeFrequency},
		},
	).WithCommentPrefix("#").WithExtension(".dict.yaml").
		WithStartMarker("...")
}

// ------------- 五笔模板 -------------

func NewDuoduoWubi() *CustomText {
	return NewCustom(
		"duoduo",
		"多多生成器",
		model.FormatTypeWubi,
		utf16,
		[]FieldConfig{
			{Type: FieldTypeWord},
			{Type: FieldTypeTab},
			{Type: FieldTypeCode},
		},
	).WithSortByCode(true).WithCommentPrefix("$ddcmd")
}

func NewRimeWubi() *CustomText {
	return NewCustom(
		"rime_wubi",
		"Rime五笔",
		model.FormatTypeWubi,
		utf8,
		[]FieldConfig{
			{Type: FieldTypeWord},
			{Type: FieldTypeTab},
			{Type: FieldTypeCode},
		},
	).WithSortByCode(true).WithCommentPrefix("#").
		WithExtension(".dict.yaml").WithStartMarker("...")
}

func NewBaiduShouji() *CustomText {
	return NewCustom(
		"baidu_shouji",
		"百度手机码表",
		model.FormatTypeWubi,
		utf16,
		[]FieldConfig{
			{Type: FieldTypeCode},
			{Type: FieldTypeTab},
			{Type: FieldTypeWord},
		},
	).WithSortByCode(true)
}

func NewBingling() *CustomText {
	return NewCustom(
		"bingling",
		"冰凌",
		model.FormatTypeWubi,
		utf16,
		[]FieldConfig{
			{Type: FieldTypeCode},
			{Type: FieldTypeTab},
			{Type: FieldTypeWord},
		},
	).WithSortByCode(true).WithStartMarker("[CODETABLE]")
}

func NewUserPhrase() *CustomText {
	return NewCustom(
		"user_phrase",
		"用户自定义短语",
		model.FormatTypeWubi,
		utf8,
		[]FieldConfig{
			{Type: FieldTypeWord},
			{Type: FieldTypeLiteral, Literal: ","},
			{Type: FieldTypeCode},
			{Type: FieldTypeLiteral, Literal: ","},
			{Type: FieldTypeRank},
		},
	)
}
