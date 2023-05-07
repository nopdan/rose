package rose

import (
	"strings"
)

type Format interface {
	GetDict() *Dict
	Parse()
	GenFrom(WordLibrary) []byte
}

func Parse(path, format string) *Dict {
	fm := DetectFormat(path, format)
	return FParse(path, fm)
}

func FParse(path string, fm Format) *Dict {
	d := fm.GetDict()
	d.read(path)
	fm.Parse()
	return d
}

func Generate(wl WordLibrary, format string) []byte {
	// 要转为的格式
	fm := NewFormat(format)
	// d := fm.GetDict()
	data := fm.GenFrom(wl)
	return data
}

func NewFormat(format string) Format {
	var fm Format
	switch format {
	// 二进制拼音词库
	case "baidu_bdict", "baidu_bcd", "bdict", "bcd":
		fm = NewBaiduBdict()
	case "qq_qpyd", "qpyd":
		fm = NewQqQpyd()
	case "sogou_scel", "qq_qcel", "scel", "qcel":
		fm = NewSogouScel()
	case "sogou_bin", "sgbin":
		fm = NewSogouBin()
	case "ziguang_uwl", "uwl":
		fm = NewZiguangUwl()
	case "mspy_udl", "udl":
		fm = NewMspyUDL()
	// 纯文本拼音
	case "jiajia", "pyjj", "jj":
		fm = NewJiaJia()
	case "word_only", "w":
		fm = NewWordOnly()
	case "sogou", "sg":
		fm = NewPinyin("sg")
	case "qq":
		fm = NewPinyin("qq")
	case "baidu", "bd":
		fm = NewPinyin("bd")
	case "google", "gg":
		fm = NewPinyin("gg")
	case "rime":
		fm = NewPinyin("rime")

	// 二进制字词码表
	case "msudp_dat", "mspy_dat", "udp":
		fm = NewMsUDP()
	case "mswb_lex", "lex":
		fm = NewMswbLex()
	case "baidu_def", "def":
		fm = NewBaiduDef()
	case "jidian_mb", "jdmb":
		fm = NewJidianMb()
	case "fcitx4_mb", "f4mb":
		fm = NewFcitx4Mb()
	// 字词的纯文本
	case "duoduo", "dd":
		fm = NewWubi("dd")
	case "bingling", "bl":
		fm = NewWubi("bl")
	case "jidian", "jd":
		fm = NewJidian()
	default:
		panic("输入格式不支持：" + format)
	}
	return fm
}

func DetectFormat(path string, format string) Format {
	var fm Format
	if format != "" {
		return NewFormat(format)
	}

	tmp := strings.Split(path, ".")
	if len(tmp) < 1 {
		return NewPinyin("rime")
	}
	suffix := tmp[len(tmp)-1]
	switch suffix {
	case "bdict", "bcd", "qpyd", "scel", "qcel", "uwl", "lex", "def":
		fm = NewFormat(suffix)
	case "bin":
		fm = NewSogouBin()
	case "dat":
		fm = NewMsUDP()
	case "mb":
		fm = NewJidianMb()
	default:
		fm = NewPinyin("rime")
	}
	return fm
}
