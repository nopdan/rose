package wubi

import (
	"bytes"
)

type Entry struct {
	Word string
	Code string
	Rank int
}

type Format interface {
	GetName() string
	GetID() string
	GetCanMarshal() bool
	Marshal([]*Entry, bool) []byte
	GetHasRank() bool

	Unmarshal(*bytes.Reader) []*Entry
	GetKind() int
}

type Template struct {
	Name       string // 名称
	ID         string // 唯一 ID
	CanMarshal bool   // 是否可以序列化
	HasRank    bool   // 是否有 Rank
}

func (t *Template) GetName() string               { return t.Name }
func (t *Template) GetID() string                 { return t.ID }
func (t *Template) GetCanMarshal() bool           { return t.CanMarshal }
func (t *Template) Marshal([]*Entry, bool) []byte { return nil }
func (t *Template) GetHasRank() bool              { return t.HasRank }
func (t *Template) GetKind() int                  { return WUBI }

const WUBI = 2

var FormatList = make([]Format, 0, 20)

func GenRank(di []*Entry) []*Entry {
	codeMap := make(map[string]int)
	for i := range di {
		codeMap[di[i].Code]++
		di[i].Rank = codeMap[di[i].Code]
	}
	return di
}

func New(format string) Format {
	for _, v := range FormatList {
		if v.GetID() == format {
			return v
		}
	}
	return nil
}
