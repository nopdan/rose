package jiajia

import (
	"bufio"
	"io"
	"strings"

	"github.com/nopdan/rose/model"
)

type JiaJia struct {
	model.BaseFormat
}

func New() *JiaJia {
	return &JiaJia{
		BaseFormat: model.BaseFormat{
			ID:          "jiajia",
			Name:        "拼音加加",
			Type:        model.FormatTypePinyin,
			Extension:   ".txt",
			Description: "拼音加加格式",
		},
	}
}

func (f *JiaJia) Import(src model.Source) ([]*model.Entry, error) {
	textReader, _, closeFn, err := model.OpenTextReader(src)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	entries := make([]*model.Entry, 0)
	scan := bufio.NewScanner(textReader)
	for scan.Scan() {
		line := []rune(strings.TrimSpace(scan.Text()))
		// 注释
		if len(line) == 0 || line[0] == ';' {
			continue
		}
		wordRunes := make([]rune, 0, 1)
		pinyin := make([]string, 0, 1)
		for i := 0; i < len(line); {
			char := line[i]
			wordRunes = append(wordRunes, char)
			// 下一个是英文（拼音）
			if i+1 != len(line) && line[i+1] < 128 {
				j := 1 // 已匹配的字母数
				for i+j+1 < len(line) && line[i+j+1] < 128 {
					j++
				}
				pinyin = append(pinyin, string(line[i+1:i+j+1]))
				i += j + 1
				continue
			}
			// 自动为汉字生成拼音（简化处理，假设字符本身就是拼音）
			pinyin = append(pinyin, string(char))
			i++
		}

		word := string(wordRunes)
		entry := model.NewEntry(word).
			WithMultiCode(pinyin...).
			WithFrequency(1).
			WithCodeType(model.CodeTypeIncompletePinyin)
		entries = append(entries, entry)
		f.Debugf("%s\t%s\t%d\n", entry.Word, entry.Code, entry.Frequency)
	}
	return entries, scan.Err()
}

func (f *JiaJia) Export(entries []*model.Entry, w io.Writer) error {
	var buf strings.Builder
	for _, v := range entries {
		words := []rune(v.Word)
		pinyins := v.Code.Strings()
		if len(words) != len(pinyins) {
			continue
		}
		for i := range words {
			buf.WriteString(string(words[i]))
			buf.WriteString(pinyins[i])
		}
		buf.WriteString("\r\n")
	}
	_, err := io.WriteString(w, buf.String())
	return err
}
