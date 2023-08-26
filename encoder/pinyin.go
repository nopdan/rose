package encoder

import (
	"github.com/nopdan/pinyin"
	"github.com/nopdan/rose/encoder/data"
)

type Pinyin struct {
	py *pinyin.Pinyin
}

func NewPinyin() *Pinyin {
	py := pinyin.New()
	py.AddReader(data.Duoyin)
	py.AddReader(data.Pinyin)
	py.AddReader(data.Correct)
	return &Pinyin{py: py}
}

// 自动注音
func (p Pinyin) Encode(word string) []string {
	return p.py.Match(word)
}
