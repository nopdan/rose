package core

import (
	"github.com/nopdan/rose/encoder"
	"github.com/nopdan/rose/pinyin"
	"github.com/nopdan/rose/wubi"
)

func Encode(di []*wubi.Entry, schema string) []*wubi.Entry {
	enc := encoder.New(schema)
	new := make([]*wubi.Entry, 0, len(di))
	for _, v := range di {
		codes := enc.Encode(v.Word)
		for _, c := range codes {
			new = append(new, &wubi.Entry{
				Word: v.Word,
				Code: c,
			})
		}
	}
	return new
}

// 拼音转为五笔词库
//
// 形码方案，wubi86/wubi98/wubi08/phrase
func ToWubi(di []*pinyin.Entry, schema string) []*wubi.Entry {
	enc := encoder.New(schema)
	new := make([]*wubi.Entry, 0, len(di))
	for _, e := range di {
		codes := enc.Encode(e.Word)
		for _, c := range codes {
			new = append(new, &wubi.Entry{
				Word: e.Word,
				Code: c,
			})
		}
	}
	new = wubi.GenRank(new)
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
