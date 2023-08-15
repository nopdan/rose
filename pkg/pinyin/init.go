package pinyin

import (
	"github.com/nopdan/pinyin"
	"github.com/nopdan/rose/pkg/data"
)

var py = pinyin.New()
var Match = py.Match

func init() {
	py.AddReader(data.Decompress(data.Pinyin))
	py.AddReader(data.Decompress(data.Duoyin))
	py.AddReader(data.Decompress(data.Correct))
}
