package rime

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/nopdan/rose/model"
)

// Format Rime格式实现
type Format struct {
	model.BaseFormat
}

// New 创建 Rime 格式处理器
func New() *Format {
	return &Format{
		BaseFormat: model.BaseFormat{
			ID:          "rime",
			Name:        "Rime输入法",
			Type:        model.FormatTypePinyin,
			Extension:   ".dict.yaml",
			Description: "Rime输入法词典格式",
		},
	}
}

// Info 返回格式信息
func (f *Format) Info() *model.BaseFormat {
	return &f.BaseFormat
}

// Import 从文件导入词条
func (f *Format) Import(src model.Source) ([]*model.Entry, error) {
	textReader, _, closeFn, err := model.OpenTextReader(src)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	var entries []*model.Entry
	scanner := bufio.NewScanner(textReader)
	inDataSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过注释和空行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 检查是否进入数据段
		if line == "..." {
			inDataSection = true
			continue
		}

		// 只处理数据段的内容
		if !inDataSection {
			continue
		}

		// 解析词条: word	pinyin	frequency
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		entry := model.NewEntry(parts[0]).WithSimpleCode(parts[1])

		// 如果有频率信息
		if len(parts) >= 3 && parts[2] != "" {
			entry.WithFrequency(1) // 简化处理，实际应解析数字
		}

		entries = append(entries, entry)
		f.Debugf("%s\t%s\t%d\n", entry.Word, entry.Code, entry.Frequency)
	}

	return entries, scanner.Err()
}

// Export 导出词条到文件
func (f *Format) Export(entries []*model.Entry, w io.Writer) error {
	// 写入文件头
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	fmt.Fprintln(writer, "# Rime dictionary")
	fmt.Fprintln(writer, "# encoding: utf-8")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "---")
	fmt.Fprintln(writer, "name: converted")
	fmt.Fprintln(writer, "version: \"1.0\"")
	fmt.Fprintln(writer, "sort: by_weight")
	fmt.Fprintln(writer, "use_preset_vocabulary: false")
	fmt.Fprintln(writer, "...")
	fmt.Fprintln(writer, "")

	// 写入词条数据
	for _, entry := range entries {
		// 检查编码是否为空
		if entry.Code == nil || entry.Code.IsEmpty() {
			// 如果没有编码，跳过或使用默认编码
			continue
		}

		// 格式: word	pinyin	frequency
		if entry.Frequency > 0 {
			fmt.Fprintf(writer, "%s\t%s\t%d\n", entry.Word, entry.Code.String(), entry.Frequency)
		} else {
			fmt.Fprintf(writer, "%s\t%s\n", entry.Word, entry.Code.String())
		}
	}

	return nil
}

// Detect 检测文件是否为Rime格式
func (f *Format) Detect(src model.Source) bool {
	textReader, _, closeFn, err := model.OpenTextReader(src)
	if err != nil {
		return false
	}
	defer closeFn()

	scanner := bufio.NewScanner(textReader)
	for i := 0; i < 10 && scanner.Scan(); i++ { // 只检查前10行
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "# Rime") || strings.Contains(line, "name:") || line == "---" {
			return true
		}
	}

	return false
}
