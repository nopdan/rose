package core

import (
	"github.com/nopdan/rose/pkg/encoder"
	"github.com/nopdan/rose/pkg/pinyin"
	"github.com/nopdan/rose/pkg/wubi"
)

// 拼音转为五笔词库
//
// 形码方案，wubi86/wubi98/wubi08/phrase
func (c *Config) ToWubi(di []string, hasRank bool) []byte {
	enc := encoder.New(c.Schema, c.AABC)
	new := make([]*wubi.Entry, 0, len(di))
	for _, word := range di {
		code := enc.Encode(word)
		new = append(new, &wubi.Entry{
			Word: word,
			Code: code,
		})
	}
	return wubi.New(c.OFormat).Marshal(new, hasRank)
}

// 转换为拼音词库
func (c *Config) ToPinyin(di []string) []byte {
	enc := encoder.NewPinyin()
	new := make([]*pinyin.Entry, 0, len(di))
	for _, word := range di {
		new = append(new, &pinyin.Entry{
			Word:   word,
			Pinyin: enc.Encode(word),
			Freq:   1,
		})
	}
	return pinyin.New(c.OFormat).Marshal(new)
}

func WubiToWords(di []*wubi.Entry) []string {
	new := make([]string, 0, len(di))
	for _, e := range di {
		new = append(new, e.Word)
	}
	return new
}

func PinyinToWords(di []*pinyin.Entry) []string {
	new := make([]string, 0, len(di))
	for _, e := range di {
		new = append(new, e.Word)
	}
	return new
}
