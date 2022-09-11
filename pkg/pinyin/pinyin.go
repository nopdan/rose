package pinyin

import "github.com/cxcn/dtool/pkg/util"

var (
	ReadUint16 = util.ReadUint16
	ReadUint32 = util.ReadUint32
	BytesToInt = util.BytesToInt
)

// 拼音词库
type Dict []Entry

// 词，拼音，词频
type Entry struct {
	Word   string
	Pinyin []string
	Freq   int
}

type Parser interface {
	Parse(string) Dict
}

type Generator interface {
	Gen(Dict) []byte
}

// 拼音
func Parse(format, filename string) Dict {
	var p Parser
	switch format {
	case "baidu_bdict", "baidu_bcd":
		p = BaiduBdict{}
	case "sogou_scel", "qq_qcel":
		p = SogouScel{}
	case "ziguang_uwl":
		p = ZiguangUwl{}
	case "qq_qpyd":
		p = QqQpyd{}
	case "mspy_dat":
		p = MspyDat{}
	case "mspy_udl":
		p = MspyUDL{}
	case "sogou_bin":
		p = SogouBin{}
	// 纯文本拼音
	case "pyjj":
		p = JiaJia{}
	case "word_only":
		p = WordOnly{}
	case "sogou":
		p = Sogou
	case "qq":
		p = QQ
	case "baidu":
		p = Baidu
	case "google":
		p = Google
	default:
		panic("输入格式不支持：" + format)
	}
	return p.Parse(filename)
}

func Generate(format string, dict Dict) []byte {
	var g Generator
	switch format {
	case "pyjj":
		g = JiaJia{}
	case "word_only":
		g = WordOnly{}
	case "sogou":
		g = Sogou
	case "qq":
		g = QQ
	case "baidu":
		g = Baidu
	case "google":
		g = Google
	default:
		panic("输出格式不支持：" + format)
	}
	return g.Gen(dict)
}
