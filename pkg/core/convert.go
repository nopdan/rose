package core

import (
	"github.com/nopdan/rose/encoder"
	"github.com/nopdan/rose/pinyin"
	"github.com/nopdan/rose/wubi"
)

func (c *Config) Encode(di []*wubi.Entry) []*wubi.Entry {
	enc := encoder.New(c.Schema, c.AABC)
	new := make([]*wubi.Entry, 0, len(di))
	for _, v := range di {
		code := enc.Encode(v.Word)
		new = append(new, &wubi.Entry{
			Word: v.Word,
			Code: code,
		})
	}
	return new
}

// 拼音转为五笔词库
//
// 形码方案，wubi86/wubi98/wubi08/phrase
func (c *Config) ToWubi(di []*pinyin.Entry) []*wubi.Entry {
	enc := encoder.New(c.Schema, c.AABC)
	new := make([]*wubi.Entry, 0, len(di))
	for _, e := range di {
		code := enc.Encode(e.Word)
		new = append(new, &wubi.Entry{
			Word: e.Word,
			Code: code,
		})
	}
	return new
}

// 转换为拼音词库
func ToPinyin(di []*wubi.Entry) []*pinyin.Entry {
	enc := encoder.NewPinyin()
	new := make([]*pinyin.Entry, 0, len(di))
	for _, e := range di {
		new = append(new, &pinyin.Entry{
			Word:   e.Word,
			Pinyin: enc.Encode(e.Word),
			Freq:   1,
		})
	}
	return new
}
