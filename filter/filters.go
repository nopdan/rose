package filter

import (
	"regexp"
	"unicode"

	"github.com/nopdan/rose/model"
)

// 过滤器用于对词条进行筛选
type Filter interface {
	// 过滤词条，返回 true 则过滤掉
	Filter(entry *model.Entry) bool
}

// LengthFilter 词长过滤器
type LengthFilter struct {
	MinLength int
	MaxLength int
}

// NewLengthFilter 创建词长过滤器
func NewLengthFilter(minLength, maxLength int) *LengthFilter {
	return &LengthFilter{MinLength: minLength, MaxLength: maxLength}
}

// Filter 实现Filter接口
func (f *LengthFilter) Filter(entry *model.Entry) bool {
	if f.MinLength <= 0 && f.MaxLength <= 0 {
		return false
	}
	length := len([]rune(entry.Word))
	if f.MinLength > 0 && length < f.MinLength {
		return true
	}
	if f.MaxLength > 0 && length > f.MaxLength {
		return true
	}
	return false
}

// FrequencyFilter 词频过滤器
type FrequencyFilter struct {
	MinFrequency int
	MaxFrequency int
}

// NewFrequencyFilter 创建词频过滤器
func NewFrequencyFilter(minFrequency, maxFrequency int) *FrequencyFilter {
	return &FrequencyFilter{MinFrequency: minFrequency, MaxFrequency: maxFrequency}
}

// Filter 实现Filter接口
func (f *FrequencyFilter) Filter(entry *model.Entry) bool {
	if f.MinFrequency <= 0 && f.MaxFrequency <= 0 {
		return false
	}
	if f.MinFrequency > 0 && entry.Frequency < f.MinFrequency {
		return true
	}
	if f.MaxFrequency > 0 && entry.Frequency > f.MaxFrequency {
		return true
	}
	return false
}

// CharacterFilter 字符类型过滤器
type CharacterFilter struct {
	FilterEnglish bool
	FilterNumber  bool
}

// NewCharacterFilter 创建字符过滤器
func NewCharacterFilter(filterEnglish, filterNumber bool) *CharacterFilter {
	return &CharacterFilter{FilterEnglish: filterEnglish, FilterNumber: filterNumber}
}

// Filter 实现Filter接口
func (f *CharacterFilter) Filter(entry *model.Entry) bool {
	if !f.FilterEnglish && !f.FilterNumber {
		return false
	}

	if f.FilterEnglish && f.containsEnglish(entry.Word) {
		return true
	}

	if f.FilterNumber && f.containsNumber(entry.Word) {
		return true
	}

	return false
}

// containsEnglish 检查是否包含英文字符
func (f *CharacterFilter) containsEnglish(text string) bool {
	for _, r := range text {
		if unicode.IsLetter(r) && r < 128 {
			return true
		}
	}
	return false
}

// containsNumber 检查是否包含数字
func (f *CharacterFilter) containsNumber(text string) bool {
	for _, r := range text {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// RegexFilter 正则表达式过滤器
type RegexFilter struct {
	patterns []*regexp.Regexp
}

// NewRegexFilter 创建正则过滤器（忽略无效规则）
func NewRegexFilter(rules []string) *RegexFilter {
	var patterns []*regexp.Regexp
	for _, rule := range rules {
		pattern, err := regexp.Compile(rule)
		if err != nil {
			continue
		}
		patterns = append(patterns, pattern)
	}

	return &RegexFilter{patterns: patterns}
}

// Filter 实现Filter接口
func (f *RegexFilter) Filter(entry *model.Entry) bool {
	if len(f.patterns) == 0 {
		return false
	}
	for _, pattern := range f.patterns {
		if pattern.MatchString(entry.Word) {
			return true
		}
	}
	return false
}
