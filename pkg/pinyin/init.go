package pinyin

import (
	"github.com/flowerime/pinyin"
	"github.com/flowerime/rose/pkg/data"
)

var py = pinyin.New()
var Match = py.Match

func init() {
	py.AddReader(data.Decompress(data.Pinyin))
	py.AddReader(data.Decompress(data.Duoyin))
	py.AddReader(data.Decompress(data.Correct))
}
