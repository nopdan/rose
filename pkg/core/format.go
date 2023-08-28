package core

import (
	"fmt"

	"github.com/nopdan/rose/pinyin"
	"github.com/nopdan/rose/wubi"
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
}

func PrintFormatList() {
	for _, f := range FormatList {
		if !f.CanMarshal {
			fmt.Printf("不")
		}
		fmt.Printf("可导出")
		fmt.Printf(" %s \t %s\n", f.ID, f.Name)
	}
}
