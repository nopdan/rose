package table

import "github.com/cxcn/dtool/pkg/util"

var (
	ReadUint16 = util.ReadUint16
	ReadUint32 = util.ReadUint32
	BytesToInt = util.BytesToInt
)

// 多多形式码表
type Table []Entry

// 词，编码
type Entry struct {
	Word  string
	Code  string
	Order byte
}

type Parser interface {
	Parse(string) Table
}

type Generator interface {
	Gen(Table) []byte
}

// 字词
func Parse(format, filename string) Table {
	var p Parser
	switch format {
	// 字词的
	case "msudp_dat":
		p = MsUDP{}
	case "mswb_lex":
		p = MswbLex{}
	case "baidu_def":
		p = BaiduDef{}
	case "jidian_mb":
		p = JidianMb{}
	case "fcitx4_mb":
		p = Fcitx4Mb{}
	// 字词的纯文本
	case "duoduo":
		p = DuoDuo
	case "bingling":
		p = Bingling
	case "jidian":
		p = Jidian{}
	default:
		panic("输入格式不支持：" + format)
	}
	return p.Parse(filename)
}

func Generate(format string, table Table) []byte {
	var g Generator
	switch format {
	case "msudp_dat":
		g = MsUDP{}
	case "duoduo":
		g = DuoDuo
	case "bingling":
		g = Bingling
	case "jidian":
		g = Jidian{}
	case "baidu_def":
		g = BaiduDef{}
	default:
		panic("输出格式不支持：" + format)
	}
	return g.Gen(table)
}
