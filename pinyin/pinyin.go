package pinyin

import (
	"bytes"
)

type Entry struct {
	Word   string
	Pinyin []string
	Freq   int
}

func (e *Entry) GetWord() string { return e.Word }

type Format interface {
	GetName() string
	GetID() string
	GetCanMarshal() bool
	Marshal([]*Entry) []byte

	Unmarshal(*bytes.Reader) []*Entry
	GetKind() int
}

type Template struct {
	Name       string // 名称
	ID         string // 唯一 ID
	CanMarshal bool   // 是否可以序列化
}

func (t *Template) GetName() string         { return t.Name }
func (t *Template) GetID() string           { return t.ID }
func (t *Template) GetCanMarshal() bool     { return t.CanMarshal }
func (t *Template) Marshal([]*Entry) []byte { return nil }
func (t *Template) GetKind() int            { return PINYIN }

const PINYIN = 1

var FormatList = make([]Format, 0, 20)

func New(format string) Format {
	for _, v := range FormatList {
		if v.GetID() == format {
			return v
		}
	}
	return nil
}
