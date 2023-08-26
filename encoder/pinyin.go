package encoder

import (
	"github.com/nopdan/pinyin"
)

type Pinyin struct {
	py *pinyin.Pinyin
}

func NewPinyin() *Pinyin {
	py := pinyin.New()
	py.AddFile("./data/duoyin.txt")
	py.AddFile("./data/pinyin.txt")
	py.AddFile("./data/correct.txt")
	return &Pinyin{py: py}
}

// 自动注音
func (p Pinyin) Encode(word string) []string {
	return p.py.Match(word)
}
