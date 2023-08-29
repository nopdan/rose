package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nopdan/rose/pkg/pinyin"
	"github.com/nopdan/rose/pkg/wubi"
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
	reader  *bytes.Reader

	Schema  string // 形码方案
	MbData  []byte // 自定义码表
	AABC    bool   // 三字词规则
	OFormat string // 输出格式
	OName   string // 保存文件名
	out     *Format
}

// 初始化
func (c *Config) init() {
	if c.IData == nil {
		data, err := os.ReadFile(c.IName)
		if err != nil {
			fmt.Printf("读取文件失败：%s\n", c.IName)
			panic(err)
		}
		c.IData = data
	}
	c.reader = bytes.NewReader(c.IData)

	c.in = matchFormat(c.IFormat)
	if c.in == nil {
		fmt.Printf("输入格式无效：%s\n", c.IFormat)
		os.Exit(1)
	}
	c.out = matchFormat(c.OFormat)
	if c.out == nil {
		fmt.Printf("输出格式无效：%s\n", c.OFormat)
		os.Exit(1)
	}
	if !c.out.CanMarshal {
		fmt.Printf("不支持导出该格式：%s\n", c.OFormat)
		os.Exit(1)
	}

	if c.Schema == "" {
		c.Schema = "86"
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

	var words []string
	switch c.in.Kind {
	// 输入拼音
	case pinyin.PINYIN:
		di := pinyin.New(c.IFormat).Unmarshal(c.reader)
		switch c.out.Kind {
		case pinyin.PINYIN:
			return pinyin.New(c.OFormat).Marshal(di)
		case wubi.WUBI:
			if c.Schema != "original" {
				words = PinyinToWords(di)
				return c.ToWubi(words, false)
			}
			new := make([]*wubi.Entry, 0, len(di))
			for _, e := range di {
				code := strings.Join(e.Pinyin, "")
				new = append(new, &wubi.Entry{
					Word: e.Word,
					Code: code,
				})
			}
			return wubi.New(c.OFormat).Marshal(new, false)
		case wubi.WORDS:
			words = PinyinToWords(di)
		}
	// 输入五笔
	case wubi.WUBI:
		f := wubi.New(c.IFormat)
		hasRank := f.GetHasRank()
		di := f.Unmarshal(c.reader)
		words = WubiToWords(di)
		switch c.out.Kind {
		case pinyin.PINYIN:
			return c.ToPinyin(words)
		case wubi.WUBI:
			if c.Schema == "original" {
				return wubi.New(c.OFormat).Marshal(di, hasRank)
			}
			words = WubiToWords(di)
			return c.ToWubi(words, hasRank)
		}
	// 输入纯词组
	case wubi.WORDS:
		words = wubi.NewWords().UnmarshalStr(c.reader)
		switch c.out.Kind {
		case pinyin.PINYIN:
			return c.ToPinyin(words)
		case wubi.WUBI:
			return c.ToWubi(words, false)
		}
	}
	return wubi.NewWords().MarshalStr(words)
}

func (c *Config) Save(data []byte) error {
	os.MkdirAll(filepath.Dir(c.OName), os.ModePerm)
	return os.WriteFile(c.OName, data, os.ModePerm)
}
