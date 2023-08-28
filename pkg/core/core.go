package core

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/nopdan/rose/pinyin"
	"github.com/nopdan/rose/wubi"
)

func matchFormat(format string) *Format {
	for _, f := range FormatList {
		if f.ID == format {
			return f
		}
	}
	return nil
}

type Config struct {
	IName   string // 文件名
	IData   []byte // 输入，如果为空则读取IName
	IFormat string // 输入格式
	in      *Format

	Schema  string // 形码方案
	MbData  []byte // 自定义码表
	AABC    bool   // 三字词规则
	OFormat string // 输出格式
	OName   string // 保存文件名
	out     *Format
}

// 初始化
func (c *Config) init() {
	c.in = matchFormat(c.IFormat)
	c.out = matchFormat(c.OFormat)
	if c.in == nil || c.out == nil {
		panic("invalid format")
	}
	if !c.out.CanMarshal {
		panic("不支持该格式的导出")
	}

	if c.IData == nil {
		data, err := os.ReadFile(c.IName)
		if err != nil {
			panic(err)
		}
		c.IData = data
	}

	if c.OName == "" {
		c.OName = filepath.Base(c.IName) + "_" + c.out.Name
		s := strings.Split(c.out.Name, ".")
		if len(s) == 1 {
			c.OName += ".txt"
		}
	}
}

func (c *Config) Marshal() []byte {
	c.init()
	r := bytes.NewReader(c.IData)

	var data []byte
	switch c.out.Kind {
	// 输出为拼音
	case pinyin.PINYIN:
		var di []*pinyin.Entry
		// 拼音转拼音
		if c.in.Kind == pinyin.PINYIN {
			di = pinyin.New(c.IFormat).Unmarshal(r)
		} else { // 五笔转拼音
			src := wubi.New(c.IFormat).Unmarshal(r)
			di = ToPinyin(src)
		}
		data = pinyin.New(c.OFormat).Marshal(di)
	// 输出五笔
	case wubi.WUBI:
		var di []*wubi.Entry
		hasRank := false
		// 五笔转五笔
		if c.in.Kind == wubi.WUBI {
			f := wubi.New(c.IFormat)
			di = f.Unmarshal(r)
			// 保留原来的编码方案
			if c.Schema == "original" || c.Schema == "" {
				hasRank = f.GetHasRank()
			} else {
				di = c.Encode(di)
			}
		} else { // 拼音转五笔
			src := pinyin.New(c.IFormat).Unmarshal(r)
			di = c.ToWubi(src)
		}
		data = wubi.New(c.OFormat).Marshal(di, hasRank)
	// 纯词组
	case wubi.WORDS:
		var di []string
		if c.in.Kind == wubi.WUBI {
			src := wubi.New(c.IFormat).Unmarshal(r)
			for _, v := range src {
				di = append(di, v.Word)
			}
		} else {
			src := pinyin.New(c.IFormat).Unmarshal(r)
			for _, v := range src {
				di = append(di, v.Word)
			}
		}
		data = wubi.NewWords().MarshalStr(di)
	}
	return data
}

func (c *Config) Save(data []byte) error {
	os.MkdirAll(filepath.Dir(c.OName), os.ModePerm)
	return os.WriteFile(c.OName, data, os.ModePerm)
}
