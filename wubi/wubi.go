package wubi

import (
	"bytes"
)

type Entry struct {
	Word string
	Code string
	Rank int
}

func (e *Entry) GetWord() string { return e.Word }

type Format interface {
	Unmarshal(*bytes.Reader) []*Entry
	Marshal([]*Entry, bool) []byte // dict, hasRank
	HasRank() bool
	GetName() string
	GetID() string
	GetType() int
}

type Template struct {
	Name string
	ID   string
	Rank bool
}

const WUBI = 2

func (t *Template) Marshal([]*Entry, bool) []byte { return nil }
func (t *Template) HasRank() bool                 { return t.Rank }
func (t *Template) GetName() string               { return t.Name }
func (t *Template) GetID() string                 { return t.ID }
func (t *Template) GetType() int                  { return WUBI }

func GenRank(di []*Entry) []*Entry {
	codeMap := make(map[string]int)
	for i := range di {
		codeMap[di[i].Code]++
		di[i].Rank = codeMap[di[i].Code]
	}
	return di
}

func New(format string) Format {
	switch format {
	// 微软用户自定义短语
	case "msudp_dat", "mspy_dat", "udp":
		return NewMsUDP()
	// 微软五笔
	case "mswb_lex", "lex":
		return NewMswbLex()
	// 百度手机自定义词库
	case "baidu_def", "def":
		return NewBaiduDef()
	// 极点五笔
	case "jidian", "jd":
		return NewJidian()
	case "jidian_mb", "jdmb":
		return NewJidianMb()
	// fcitx4.mb
	case "fcitx4_mb", "f4mb":
		return NewFcitx4Mb()
	// 多多生成器
	case "duodb":
		return NewDuoDB()
	case "dddmg", "dmg":
		return NewDuoduo()
	case "duoduo", "dd":
		return NewDuoduo()
	// 冰凌
	case "bingling", "bl":
		return NewBingling()
	case "words":
		return NewWords()
	}
	return nil
}
