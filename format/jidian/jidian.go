package jidian

import (
	"bufio"
	"bytes"
	"io"
	"slices"
	"strings"

	"github.com/nopdan/rose/model"
)

type Jidian struct {
	model.BaseFormat
}

func New() *Jidian {
	return &Jidian{
		BaseFormat: model.BaseFormat{
			ID:          "jidian",
			Name:        "极点码表",
			Type:        model.FormatTypeWubi,
			Extension:   ".txt",
			Description: "极点五笔码表格式",
		},
	}
}

func (f *Jidian) Import(src model.Source) ([]*model.Entry, error) {
	textReader, _, closeFn, err := model.OpenTextReader(src)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	entries := make([]*model.Entry, 0)
	scan := bufio.NewScanner(textReader)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), " ")
		if len(entry) < 2 {
			continue
		}
		for i := 1; i < len(entry); i++ {
			entries = append(entries, model.NewEntry(entry[i]).WithSimpleCode(entry[0]).WithRank(i))
			f.Debugf("%s\t%s\t%d\n", entry[i], entry[0], i)
		}
	}
	return entries, scan.Err()
}

func (f *Jidian) Export(entries []*model.Entry, w io.Writer) error {
	var buf bytes.Buffer
	buf.Grow(len(entries) * 10) // 估算大小

	// 按编码排序
	entriesCopy := make([]*model.Entry, len(entries))
	copy(entriesCopy, entries)
	slices.SortStableFunc(entriesCopy, func(a, b *model.Entry) int {
		return strings.Compare(a.Code.String(), b.Code.String())
	})

	// 记录上一个编码
	lastCode := ""
	for i, entry := range entriesCopy {
		code := entry.Code.String()
		if lastCode != code {
			if i != 0 {
				buf.WriteByte('\r')
				buf.WriteByte('\n')
			}
			buf.WriteString(code)
			buf.WriteByte('\t')
			buf.WriteString(entry.Word)
			lastCode = code
			continue
		}
		buf.WriteByte(' ')
		buf.WriteString(entry.Word)
		lastCode = code
	}
	_, err := w.Write(buf.Bytes())
	return err
}
