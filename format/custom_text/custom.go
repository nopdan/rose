package custom_text

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

// FieldType 表示字段类型
type FieldType int

const (
	FieldTypeWord FieldType = iota
	FieldTypePinyin
	FieldTypeCode
	FieldTypeFrequency
	FieldTypeRank
	FieldTypeTab     // 制表符
	FieldTypeSpace   // 空格
	FieldTypeLiteral // 自定义字面量
)

// FieldConfig 字段配置
type FieldConfig struct {
	Type FieldType

	// 拼音配置（仅对 FieldTypePinyin 有效）
	PinyinSeparator string // 拼音分隔符，如 "'"
	PinyinPrefix    string // 拼音前缀
	PinyinSuffix    string // 拼音后缀

	// 字面量配置（仅对 FieldTypeLiteral 有效）
	Literal string // 字面量内容
}

// CustomText 通用纯文本格式
//
// 支持灵活的字段组合：
// - 词组、拼音（可配置分隔符/前后缀）、编码、词频、候选顺序
// - Tab、空格、自定义字面量
type CustomText struct {
	model.BaseFormat

	// 导出的编码格式（导入时使用自动检测）
	Encoding *util.Encoding

	// 字段配置（按顺序组合）
	Fields []FieldConfig

	// 是否按编码排序
	SortByCode bool

	// 注释前缀（如 "#"），空字符串表示不支持注释
	CommentPrefix string
}

// NewCustom 创建通用纯文本格式
func NewCustom(
	id, name string,
	formatType model.FormatType,
	encoding *util.Encoding,
	fields []FieldConfig,
	sortByCode bool,
	commentPrefix string,
) *CustomText {
	if encoding == nil {
		encoding = util.NewEncoding("UTF-8")
	}

	return &CustomText{
		BaseFormat: model.BaseFormat{
			ID:          id,
			Name:        name,
			Type:        formatType,
			Extension:   ".txt",
			Description: "通用纯文本格式",
		},
		Encoding:      encoding,
		Fields:        fields,
		SortByCode:    sortByCode,
		CommentPrefix: commentPrefix,
	}
}

// Import 读取通用纯文本格式
func (f *CustomText) Import(src model.Source) ([]*model.Entry, error) {
	textReader, _, closeFn, err := model.OpenTextReader(src)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	entries := make([]*model.Entry, 0)
	scan := bufio.NewScanner(textReader)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			continue
		}
		// 跳过注释行
		if f.CommentPrefix != "" && strings.HasPrefix(line, f.CommentPrefix) {
			f.Infof("跳过注释行: %s\n", line)
			continue
		}

		entry := model.NewEntry("")
		var hasPinyin bool
		var hasCode bool
		// 当前解析位置
		pos := 0

		for _, fieldCfg := range f.Fields {
			if pos >= len(line) {
				break
			}

			switch fieldCfg.Type {
			case FieldTypeWord:
				value, nextPos := f.extractField(line, pos)
				entry.Word = value
				pos = nextPos

			case FieldTypePinyin:
				value, nextPos := f.extractField(line, pos)
				value = strings.TrimSpace(value)
				if value != "" {
					value = f.stripAffix(value, fieldCfg.PinyinPrefix, fieldCfg.PinyinSuffix)
					var codes []string
					if fieldCfg.PinyinSeparator != "" {
						codes = strings.Split(value, fieldCfg.PinyinSeparator)
					} else {
						codes = []string{value}
					}
					entry = entry.WithMultiCode(codes...)
					hasPinyin = true
				}
				pos = nextPos

			case FieldTypeCode:
				value, nextPos := f.extractField(line, pos)
				value = strings.TrimSpace(value)
				if value != "" {
					entry = entry.WithSimpleCode(value)
					hasCode = true
				}
				pos = nextPos

			case FieldTypeFrequency:
				value, nextPos := f.extractField(line, pos)
				value = strings.TrimSpace(value)
				if freq, err := strconv.Atoi(value); err == nil {
					entry.Frequency = freq
				}
				pos = nextPos

			case FieldTypeRank:
				value, nextPos := f.extractField(line, pos)
				value = strings.TrimSpace(value)
				if rank, err := strconv.Atoi(value); err == nil {
					entry.Rank = rank
				}
				pos = nextPos

			case FieldTypeTab:
				if pos < len(line) && line[pos] == '\t' {
					pos++
				}

			case FieldTypeSpace:
				if pos < len(line) && line[pos] == ' ' {
					pos++
				}

			case FieldTypeLiteral:
				if strings.HasPrefix(line[pos:], fieldCfg.Literal) {
					pos += len(fieldCfg.Literal)
				}

			}
		}

		if entry.Word == "" {
			continue
		}
		if hasPinyin {
			entry.CodeType = model.CodeTypePinyin
		} else if hasCode {
			entry.CodeType = model.CodeTypeWubi
		}

		entries = append(entries, entry)
		f.Debugf("%s\t%s\t%d\n", entry.Word, entry.Code, entry.Frequency)
	}

	return entries, nil
}

// extractField 从指定位置提取下一个字段（查找下一个分隔符）
func (f *CustomText) extractField(line string, start int) (string, int) {
	// 查找下一个 tab/space
	end := start
	for end < len(line) && line[end] != '\t' && line[end] != ' ' {
		end++
	}
	return line[start:end], end
}

// stripAffix 去除前后缀
func (f *CustomText) stripAffix(value, prefix, suffix string) string {
	if prefix != "" && strings.HasPrefix(value, prefix) {
		value = strings.TrimPrefix(value, prefix)
	}
	if suffix != "" && strings.HasSuffix(value, suffix) {
		value = strings.TrimSuffix(value, suffix)
	}
	return value
}

// Export 写入通用纯文本格式
func (f *CustomText) Export(entries []*model.Entry, w io.Writer) error {
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	if f.Encoding == utf16 {
		_, _ = writer.Write([]byte{0xFF, 0xFE})
	}

	ew := &encodingWriter{w: writer, enc: f.Encoding}

	if f.SortByCode {
		sort.SliceStable(entries, func(i, j int) bool {
			if entries[i].Code != nil && entries[j].Code != nil {
				return entries[i].Code.String() < entries[j].Code.String()
			}
			return false
		})
	}

	first := true
	for _, entry := range entries {
		if entry == nil {
			continue
		}
		if !first {
			ew.WriteString(util.LineBreak)
		}
		f.writeLine(ew, entry)
		first = false
	}

	ew.WriteString(util.LineBreak)

	return nil
}

// formatLine 将词条格式化为文本行并写入
func (f *CustomText) writeLine(w *encodingWriter, currEntry *model.Entry) {
	if currEntry == nil {
		return
	}

	for _, fieldCfg := range f.Fields {
		switch fieldCfg.Type {
		case FieldTypeWord:
			w.WriteString(currEntry.Word)

		case FieldTypePinyin:
			if currEntry.Code != nil {
				w.WriteString(fieldCfg.PinyinPrefix)
				if fieldCfg.PinyinSeparator != "" {
					w.WriteString(strings.Join(currEntry.Code.Strings(), fieldCfg.PinyinSeparator))
				} else {
					w.WriteString(strings.Join(currEntry.Code.Strings(), ""))
				}
				w.WriteString(fieldCfg.PinyinSuffix)
			}

		case FieldTypeCode:
			if currEntry.Code != nil {
				w.WriteString(currEntry.Code.String())
			}

		case FieldTypeFrequency:
			w.WriteString(strconv.Itoa(currEntry.Frequency))

		case FieldTypeRank:
			w.WriteString(strconv.Itoa(currEntry.Rank))

		case FieldTypeTab:
			w.WriteString("\t")

		case FieldTypeSpace:
			w.WriteString(" ")

		case FieldTypeLiteral:
			w.WriteString(fieldCfg.Literal)
		}
	}
}

type encodingWriter struct {
	w   io.Writer
	enc *util.Encoding
}

func (w *encodingWriter) WriteString(s string) {
	_, _ = w.w.Write(w.enc.Encode(s))
}
