package gboard

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/nopdan/rose/model"
)

// Gboard Gboard 个人字典二进制格式
type Gboard struct {
	model.BaseFormat
}

// New 创建Gboard格式处理器
func New() *Gboard {
	return &Gboard{
		BaseFormat: model.BaseFormat{
			ID:          "gboard",
			Name:        "Gboard 个人字典",
			Type:        model.FormatTypePinyin,
			Extension:   ".zip",
			Description: "Gboard 个人字典二进制格式",
		},
	}
}

func (f *Gboard) Import(src model.Source) ([]*model.Entry, error) {
	// 读取文件，解压 zip
	rd, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}
	zrd, err := zip.NewReader(rd, rd.Size())
	if err != nil {
		return nil, err
	}
	file, err := zrd.Open("dictionary.txt")
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(file)
	entries := make([]*model.Entry, 0)
	start := false
	for scan.Scan() {
		line := scan.Text()
		if line == "# Gboard Dictionary version:1" {
			start = true
			continue
		}
		if !start {
			continue
		}
		items := strings.Split(line, "\t")
		if len(items) < 3 {
			continue
		}
		code := items[0]
		word := items[1]
		locale := items[2]
		entry := model.NewEntry(word).WithSimpleCode(code)
		if locale == "zh-CN" {
			entry = entry.WithCodeType(model.CodeTypePinyinString)
		} else {
			entry = entry.WithCodeType(model.CodeTypeWubi)
		}
		entries = append(entries, entry)
		f.Debugf("%s\n", line)
	}
	return entries, nil
}

func (f *Gboard) Export(entries []*model.Entry, w io.Writer) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	writer, err := zipWriter.Create("dictionary.txt")
	if err != nil {
		return err
	}
	writer.Write([]byte("# Gboard Dictionary version:1\n"))
	for _, entry := range entries {
		line := fmt.Sprintf("%s\t%s\tzh-CN\n", entry.Code.String(), entry.Word)
		_, err := writer.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	zipWriter.Flush()
	return nil
}
