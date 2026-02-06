package encoder

import (
	"strings"
	"unicode/utf8"

	marker "github.com/nopdan/pinyin-marker"
	"github.com/nopdan/rose/model"
)

// PinyinEncoder 拼音编码器
type PinyinEncoder struct {
	*BaseEncoder
	m *marker.Marker
}

// NewPinyinEncoder 创建拼音编码器
func NewPinyinEncoder() *PinyinEncoder {
	return &PinyinEncoder{
		BaseEncoder: &BaseEncoder{},
	}
}

// Encode 将Entry转换为拼音编码（原地修改）
func (e *PinyinEncoder) Encode(entry *model.Entry) {
	if e.m == nil {
		e.m = marker.New()
	}
	switch entry.CodeType {
	case model.CodeTypePinyin:
		return // 已经是拼音编码，无需转换
	case model.CodeTypeIncompletePinyin:
		word := []rune(entry.Word)
		newCodes := make([]string, 0, len(word))
		codes := entry.Code.Strings()
		for i, c := range word {
			if i >= len(codes) || (utf8.RuneCountInString(codes[i]) == 1 && string(c) == codes[i]) {
				// 重新生成拼音
				newCodes = append(newCodes, e.m.Mark(string(c))...)
			} else {
				newCodes = append(newCodes, codes[i])
			}
		}
		entry.Code = model.NewMultiCode(newCodes...)
	case model.CodeTypePinyinString:
		// 匹配所有可能的拼音
		all := e.m.Mark2(entry.Word)
		for _, pinyin := range all {
			if strings.Join(pinyin, "") == entry.Code.String() {
				entry.Code = model.NewMultiCode(pinyin...)
				goto LABEL
			}
		}
		fallthrough
	default:
		// 其他编码类型，直接生成拼音
		entry.Code = model.NewMultiCode(e.m.Mark(entry.Word)...)
		// fmt.Printf("%s: %s\n", entry.Word, entry.Code.String())
	}
LABEL:
	entry.CodeType = model.CodeTypePinyin
}

// EncodeBatch 批量编码（原地修改）
func (e *PinyinEncoder) EncodeBatch(entries []*model.Entry) {
	for _, entry := range entries {
		e.Encode(entry)
	}
}
