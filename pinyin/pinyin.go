package pinyin

import (
	"bytes"
)

type Entry struct {
	Word   string
	Pinyin []string
	Freq   int
}

func (e *Entry) GetWord() string { return e.Word }

type Format interface {
	Unmarshal(*bytes.Reader) []*Entry
	Marshal([]*Entry) []byte
	GetName() string
	GetID() string
	GetType() int
}

type Template struct {
	Name string
	ID   string
}

const PINYIN = 1

func (t *Template) Marshal([]*Entry) []byte { return nil }
func (t *Template) GetName() string         { return t.Name }
func (t *Template) GetID() string           { return t.ID }
func (t *Template) GetType() int            { return PINYIN }

func New(format string) Format {
	switch format {
	// 搜狗输入法、qq输入法
	case "sogou", "sg":
		return NewSogou()
	case "sogou_scel", "scel", "qq_qcel", "qcel":
		return NewSogouScel()
	case "sogou_bin", "sgbin":
		return NewSogouBak()
	case "qq":
		return NewQQ()
	case "qq_qpyd", "qpyd":
		return NewQqQpyd()
	// 百度输入法
	case "baidu", "bd":
		return NewBaidu()
	case "baidu_bdict", "bdict", "baidu_bcd", "bcd":
		return NewBaiduBdict()
	// 紫光拼音
	case "ziguang_uwl", "uwl":
		return NewZiguangUwl()
	// 拼音加加
	case "jiajia", "pyjj", "jj":
		return NewJiaJia()
	// 微软自学习
	case "mspy_udl", "udl":
		return NewMspyUDL()
	// 谷歌输入法
	case "google", "gg":
		return NewGoogle()
	// rime 拼音
	case "rime":
		return NewRime()
	}
	return nil
}
