package encoder

import "strings"

func New(schema string) Encoder {
	switch schema {
	case "wubi86":
	case "wubi98":
	case "wubi08":
	case "phrase":
		return Phrase{}
	case "pinyin":
		return NewPinyin()
	}
	return nil
}

type Encoder interface {
	Encode(string) []string // 编码一个词，可能有多个编码
}

type Phrase struct {
	enc Encoder
}

func NewPhrase() *Phrase {
	enc := NewPinyin()
	return &Phrase{enc: enc}
}

func (p Phrase) Encode(word string) []string {
	return []string{
		strings.Join(p.enc.Encode(word), ""),
	}
}
