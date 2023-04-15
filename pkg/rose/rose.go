package rose

import (
	"bytes"
	"fmt"

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

func DecodeY(b []byte, e string) string {
	v, _ := Decode(b, e)
	return v
}

func PrintInfo(r *bytes.Reader, size uint32, info string) {
	tmp := make([]byte, size)
	r.Read(tmp)
	fmt.Printf("%s%s\n", info, DecodeY(tmp, "UTF-16LE"))
}

type Format interface {
	GetDict() *Dict
	Parse()
	GenFrom(WordLibrary) []byte
}

func Parse(path string, format string) *Dict {
	fmt.Println("正在解析词库：", path)
	fm := NewFormat(format)
	d := fm.GetDict()
	d.read(path)
	fmt.Printf("> 词库格式：%s -> %s\n", format, d.Name)
	fm.Parse()
	fmt.Printf("> 解析成功！词条数：%d\n\n", len(d.WordLibrary))
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
	case "sogou_bin":
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
	case "jidian_mb":
		fm = NewJidianMb()
	case "fcitx4_mb":
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
