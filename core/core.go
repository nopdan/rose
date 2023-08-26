package core

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/nopdan/rose/pinyin"
	"github.com/nopdan/rose/wubi"
)

func matchType(format string) int {
	if f := wubi.New(format); f != nil {
		return f.GetType()
	}
	if f := pinyin.New(format); f != nil {
		return f.GetType()
	}
	return 0
}

type Config struct {
	Name      string // 文件名
	InFormat  string // 输入格式
	OutFormat string // 输出格式
	OutName   string // 保存文件名
	Schema    string // 形码方案
}

func (c *Config) Marshal() {
	in := matchType(c.InFormat)
	out := matchType(c.OutFormat)
	if in == 0 || out == 0 {
		panic("invalid format")
	}

	b, err := os.ReadFile(c.Name)
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(b)
	var data []byte

	switch out {
	// 输出为拼音
	case pinyin.PINYIN:
		var di []*pinyin.Entry
		// 拼音转拼音
		if in == pinyin.PINYIN {
			di = pinyin.New(c.InFormat).Unmarshal(r)
		} else { // 五笔转拼音
			src := wubi.New(c.InFormat).Unmarshal(r)
			di = ToPinyin(src)
		}
		data = pinyin.New(c.OutFormat).Marshal(di)
	// 输出五笔
	case wubi.WUBI:
		var di []*wubi.Entry
		hasRank := false
		// 五笔转五笔
		if in == wubi.WUBI {
			f := wubi.New(c.InFormat)
			di = f.Unmarshal(r)
			// 保留原来的编码方案
			if c.Schema == "original" || c.Schema == "" {
				hasRank = f.HasRank()
			} else {
				di = Encode(di, c.Schema)
			}
		} else { // 拼音转五笔
			src := pinyin.New(c.InFormat).Unmarshal(r)
			di = ToWubi(src, c.Schema)
		}
		data = wubi.New(c.OutFormat).Marshal(di, hasRank)
	// 纯词组
	case wubi.WORDS:
		var di []string
		if in == wubi.WUBI {
			src := wubi.New(c.InFormat).Unmarshal(r)
			for _, v := range src {
				di = append(di, v.Word)
			}
		} else {
			src := pinyin.New(c.InFormat).Unmarshal(r)
			for _, v := range src {
				di = append(di, v.Word)
			}
		}
		data = wubi.NewWords().MarshalStr(di)
	}
	if len(data) == 0 {
		panic("不支持该格式的导出")
	}

	if c.OutName == "" {
		c.OutName = filepath.Base(c.Name) + ".txt"
	}
	os.MkdirAll(filepath.Dir(c.OutName), os.ModePerm)
	os.WriteFile(c.OutName, data, os.ModePerm)
}
