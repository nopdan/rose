package table

import "github.com/imetool/goutil/util"

const (
	_u16 = uint16(0)
	_u32 = uint32(0)
)

var (
	ReadUint16 = util.ReadUint16
	ReadUint32 = util.ReadUint32
	BytesToInt = util.BytesToInt
)

// 多多形式码表
type Table []Entry

// 词，编码
type Entry struct {
	Word string
	Code string
	Pos  int // 在候选中的位置
}

// 按码表顺序生成候选位置
func (t Table) GenPos() {
	count := make(map[string]int)
	for i := range t {
		count[t[i].Code] += 1
		t[i].Pos = count[t[i].Code]
	}
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
	case "pang":
		p = Pang{}
	case "pang_assoc":
		p = PangAssoc{}
	// 字词的纯文本
	case "duoduo", "dd":
		p = DuoDuo
	case "bingling", "bl":
		p = Bingling
	case "jidian", "jd":
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
	case "duoduo", "dd":
		g = DuoDuo
	case "bingling", "bl":
		g = Bingling
	case "jidian", "jd":
		g = Jidian{}
	case "baidu_def":
		g = BaiduDef{}
	default:
		panic("输出格式不支持：" + format)
	}
	return g.Gen(table)
}
