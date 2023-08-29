package core

import (
	"cmp"
	"os"
	"slices"

	"github.com/nopdan/rose/pkg/pinyin"
	"github.com/nopdan/rose/pkg/wubi"
	"github.com/olekukonko/tablewriter"
)

type Format struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	CanMarshal bool   `json:"canMarshal"`
	Kind       int    `json:"kind"`
}

var FormatList = make([]*Format, 0, 20)

func init() {
	for _, f := range pinyin.FormatList {
		FormatList = append(FormatList, &Format{
			Name:       f.GetName(),
			ID:         f.GetID(),
			CanMarshal: f.GetCanMarshal(),
			Kind:       f.GetKind(),
		})
	}
	for _, f := range wubi.FormatList {
		FormatList = append(FormatList, &Format{
			Name:       f.GetName(),
			ID:         f.GetID(),
			CanMarshal: f.GetCanMarshal(),
			Kind:       f.GetKind(),
		})
	}
	slices.SortStableFunc(FormatList, func(a, b *Format) int {
		return cmp.Compare(a.ID, b.ID)
	})
}

func PrintFormatList() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "格式", "可导出"})
	for _, f := range FormatList {
		canMarshal := ""
		if f.CanMarshal {
			canMarshal = "是"
		}
		table.Append([]string{f.ID, f.Name, canMarshal})
	}
	table.Render()
}
