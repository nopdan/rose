package model

import "strings"

// CodeType 编码类型枚举
type CodeType int

const (
	CodeTypeNone             CodeType = iota // 无编码
	CodeTypeEnglish                          // 英文
	CodeTypeWubi                             // 五笔
	CodeTypePinyin                           // 拼音
	CodeTypePinyinString                     // 拼音字符串
	CodeTypeIncompletePinyin                 // 不完整的拼音
)

// Entry 代表词库中的一个词条
type Entry struct {
	// 词组文本，如："中华人民共和国"
	Word string

	// 编码信息 - 优化为接口类型以支持不同格式
	Code Encoding

	// 词频，表示该词的使用频率，0表示未设置
	Frequency int

	// 候选顺序，在重码时的排序位置，0表示未设置
	Rank int

	// 编码类型
	CodeType CodeType
}

// NewEntry 创建一个新的词条
func NewEntry(word string) *Entry {
	return &Entry{
		Word:      word,
		Code:      nil,
		Frequency: 0,
		Rank:      0,
		CodeType:  CodeTypeNone,
	}
}

// WithSimpleCode 设置简单编码（用于五笔等）
func (e *Entry) WithSimpleCode(code string) *Entry {
	e.Code = NewSimpleCode(code)
	return e
}

// WithMultiCode 设置多段编码（用于拼音等）
func (e *Entry) WithMultiCode(codes ...string) *Entry {
	e.Code = NewMultiCode(codes...)
	return e
}

// WithFrequency 设置词频
func (e *Entry) WithFrequency(freq int) *Entry {
	e.Frequency = freq
	return e
}

// WithRank 设置候选顺序
func (e *Entry) WithRank(rank int) *Entry {
	e.Rank = rank
	return e
}

// WithCodeType 设置编码类型
func (e *Entry) WithCodeType(codeType CodeType) *Entry {
	e.CodeType = codeType
	return e
}

// Encoding 编码接口，支持不同类型的编码格式
type Encoding interface {
	// String 返回编码的字符串表示
	String() string

	// Strings 返回编码的字符串数组表示（用于拼音等多段编码）
	Strings() []string

	// IsEmpty 检查编码是否为空
	IsEmpty() bool
}

// SimpleCode 简单编码（用于五笔等单一编码）
type SimpleCode struct {
	code string
}

// NewSimpleCode 创建简单编码
func NewSimpleCode(code string) *SimpleCode {
	return &SimpleCode{code: code}
}

func (s *SimpleCode) String() string {
	return s.code
}

func (s *SimpleCode) Strings() []string {
	if s.code == "" {
		return []string{}
	}
	return []string{s.code}
}

func (s *SimpleCode) IsEmpty() bool {
	return s.code == ""
}

// MultiCode 多段编码（用于拼音等多段编码）
type MultiCode struct {
	codes []string
}

// NewMultiCode 创建多段编码
func NewMultiCode(codes ...string) *MultiCode {
	return &MultiCode{codes: codes}
}

func (m *MultiCode) String() string {
	if len(m.codes) == 0 {
		return ""
	}
	if len(m.codes) == 1 {
		return m.codes[0]
	}
	// 使用常见的分隔符连接
	var result strings.Builder
	for i, code := range m.codes {
		if i > 0 {
			result.WriteByte('\'')
		}
		result.WriteString(code)
	}
	return result.String()
}

func (m *MultiCode) Strings() []string {
	return m.codes
}

func (m *MultiCode) IsEmpty() bool {
	return len(m.codes) == 0
}
