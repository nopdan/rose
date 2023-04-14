package rose

import (
	"log"

	util "github.com/flowerime/goutil"
)

const (
	_u16 = uint16(0)
	_u32 = uint32(0)
)

var (
	ReadUint16 = util.ReadUint16
	ReadUint32 = util.ReadUint32
	BytesToInt = util.BytesToInt

	Encode = util.Encode
	Decode = util.Decode
)

type Format interface {
	GetDict() *Dict
	Parse()
	GenFrom(*Dict) []byte
}

func Parse(path string, format string) *Dict {
	fm := NewFormat(format)
	d := fm.GetDict()
	d.path = path
	d.read(path)
	fm.Parse()
	return d
}

func Generate(src *Dict, format string) []byte {
	// 要转为的格式
	fm := NewFormat(format)
	d := fm.GetDict()
	if !d.IsPinyin {
		if src.IsPinyin {
			log.Panicln("不支持拼音词库转为", format)
		}
		return fm.GenFrom(src)
	}

	// 转为拼音
	if !src.IsPinyin {
		src.ToPyTable()
	}
	data := fm.GenFrom(src)
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
	case "sogou_bin":
		fm = NewSogouBin()
	case "ziguang_uwl", "uwl":
		fm = NewZiguangUwl()
	case "msudp_dat", "mspy_dat", "udp":
		fm = NewMsUDP()
	case "mspy_udl", "udl":
		fm = NewMspyUDL()
	// 纯文本拼音
	case "jiajia", "pyjj", "jj":
		fm = NewJiaJia()
	case "word_only", "w":
		fm = NewWordOnly()
	case "sogou", "sg":
		fm = NewCommonPyTable("sg")
	case "qq":
		fm = NewCommonPyTable("qq")
	case "baidu", "bd":
		fm = NewCommonPyTable("bd")
	case "google", "gg":
		fm = NewCommonPyTable("gg")
	case "rime":
		fm = NewCommonPyTable("rime")

	// 二进制字词码表
	case "mswb_lex", "lex":
		fm = NewMswbLex()
	case "baidu_def", "def":
		fm = NewBaiduDef()
	case "jidian_mb":
		fm = NewJidianMb()
	case "fcitx4_mb":
		fm = NewFcitx4Mb()
	// 字词的纯文本
	case "duoduo", "dd":
		fm = NewCommonTable("dd")
	case "bingling", "bl":
		fm = NewCommonTable("bl")
	case "jidian", "jd":
		fm = NewJidian()
	default:
		panic("输入格式不支持：" + format)
	}

	d := fm.GetDict()
	if !d.IsBinary {
		d.Suffix = "txt"
	}
	return fm
}
