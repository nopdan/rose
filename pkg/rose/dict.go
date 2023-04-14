package rose

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	util "github.com/flowerime/goutil"
	"github.com/flowerime/rose/pkg/zhuyin"
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

func (d *Dict) GetTable() Table {
	return d.table
}

func (d *Dict) GetDict() *Dict {
	return d
}

// 默认不支持生成
func (d *Dict) GenFrom(src *Dict) []byte {
	fmt.Println("不支持生成", d.Name)
	return []byte{}
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
