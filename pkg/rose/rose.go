package rose

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	util "github.com/flowerime/goutil"
	"github.com/flowerime/rose/pkg/zhuyin"
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

type TableEntry struct {
	Word string
	Code string
	Pos  int // 重码顺序
}

type CodeEntry struct {
	Code  string
	Words []string
}

type PinyinEntry struct {
	Word   string
	Pinyin []string
	Freq   int // 词频
}

type Table = []*TableEntry
type PyTable = []*PinyinEntry
type CodeTable = []*CodeEntry

type Dict struct {
	IsPinyin bool
	pyt      PyTable
	table    Table
	codet    CodeTable

	IsBinary bool
	data     []byte
	size     int64
	rd       io.Reader

	Name   string
	Suffix string
	path   string
}

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

func genErr(name string) []byte {
	fmt.Println("不支持生成", name)
	return []byte{}
}

func (d *Dict) read(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	fi, _ := f.Stat()
	d.size = fi.Size()
	if d.IsBinary {
		d.data, _ = io.ReadAll(f)
	} else {
		d.rd = util.NewReader(f)
	}
}

// 按码表顺序生成候选位置
func (d *Dict) GenPos() {
	t := d.table
	count := make(map[string]int)
	for i := range t {
		count[t[i].Code] += 1
		t[i].Pos = count[t[i].Code]
	}
}

// 转为 极点形式，一码多词
func (d *Dict) ToCodeTable() {
	table := d.table
	sort.Slice(table, func(i, j int) bool {
		return table[i].Pos < table[j].Pos
	})
	codeMap := make(map[string][]string)
	for _, wc := range table {
		if _, ok := codeMap[wc.Code]; !ok {
			codeMap[wc.Code] = []string{wc.Word}
			continue
		}
		codeMap[wc.Code] = append(codeMap[wc.Code], wc.Word)
	}

	codet := make(CodeTable, 0, len(codeMap))
	for k, v := range codeMap {
		codet = append(codet, &CodeEntry{k, v})
	}
	sort.Slice(codet, func(i, j int) bool {
		return codet[i].Code < codet[j].Code
	})
	d.codet = codet
}

// CodeTable => Table
func (d *Dict) ToTable() {
	table := make(Table, 0, len(d.codet))
	for _, entry := range d.codet {
		for j, word := range entry.Words {
			table = append(table, &TableEntry{word, entry.Code, j + 1})
		}
	}
	d.table = table
}

// Table => PyTable
func (d *Dict) ToPyTable() {
	pyt := make(PyTable, 0, len(d.table))
	for _, entry := range d.table {
		pyt = append(pyt, &PinyinEntry{entry.Word, zhuyin.Get(entry.Word), 1})
	}
	fmt.Println("Table => PyTable")
	d.pyt = pyt
}

func (d *Dict) PyToTable(sep string) {
	table := make(Table, 0, len(d.pyt))
	for i := range d.pyt {
		table = append(table, &TableEntry{d.pyt[i].Word, strings.Join(d.pyt[i].Pinyin, sep), 1})
	}
	d.table = table
}

func (d *Dict) GetTable() Table {
	return d.table
}
