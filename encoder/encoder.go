package encoder

import (
	"bufio"
	"strings"

	"github.com/nopdan/rose/util"
)

type Encoder interface {
	Encode(string) string // 编码一个词，可能有多个编码
}

func New(schema string, isAABC bool) Encoder {
	if schema == "phrase" {
		return NewPhrase()
	}
	if w := NewWubi(schema, isAABC); w != nil {
		return w
	}
	w := &Wubi{
		Char:   make(map[rune]string),
		IsAABC: isAABC,
	}
	r, err := util.Read(schema)
	if err != nil {
		panic(err)
	}
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		fields := strings.Split(scan.Text(), "\t")
		if len(fields) < 2 {
			continue
		}
		word := []rune(fields[0])
		// 跳过词组
		if len(word) != 1 {
			continue
		}
		char := word[0]
		w.Char[char] = fields[1]
	}
	return w
}

type Phrase struct {
	enc *Pinyin
}

func NewPhrase() *Phrase {
	enc := NewPinyin()
	return &Phrase{enc: enc}
}

func (p Phrase) Encode(word string) string {
	return strings.Join(p.enc.Encode(word), "")
}
