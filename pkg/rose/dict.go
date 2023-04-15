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

type Dict struct {
	WordLibrary

	Name   string
	Suffix string // 后缀为空默认 txt 纯文本
	data   []byte
	size   int64
	rd     io.Reader
}

func (d *Dict) read(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	fi, _ := f.Stat()
	d.size = fi.Size()
	if d.Suffix != "" {
		d.data, _ = io.ReadAll(f)
	} else {
		d.rd = util.NewReader(f)
	}
}

func (d *Dict) GetDict() *Dict {
	return d
}

// 默认不支持生成
func (d *Dict) GenFrom(wl WordLibrary) []byte {
	fmt.Println("不支持生成", d.Name)
	return []byte{}
}

type WordLibrary []Entry
type Entry interface {
	GetWord() string
	GetCode() string
	GetPos() int

	GetPinyin() []string
	GetFreq() int
}

// 类五笔码表
type WubiTable []*WubiEntry
type WubiEntry struct {
	Word string
	Code string
	Pos  int // 重码顺序
}

func (e *WubiEntry) GetWord() string     { return e.Word }
func (e *WubiEntry) GetCode() string     { return e.Code }
func (e *WubiEntry) GetPos() int         { return e.Pos }
func (e *WubiEntry) GetPinyin() []string { return zhuyin.Get(e.Word) }
func (e *WubiEntry) GetFreq() int        { return 1 }

// 拼音输入法词库
type PinyinTable []*PinyinEntry
type PinyinEntry struct {
	Word   string
	Pinyin []string
	Freq   int // 词频
}

func (e *PinyinEntry) GetWord() string     { return e.Word }
func (e *PinyinEntry) GetCode() string     { return strings.Join(e.Pinyin, "") }
func (e *PinyinEntry) GetPos() int         { return 1 }
func (e *PinyinEntry) GetPinyin() []string { return e.Pinyin }
func (e *PinyinEntry) GetFreq() int        { return e.Freq }

// 类极点，一码多词
type CodeTable []*CodeEntry
type CodeEntry struct {
	Code  string
	Words []string
}

func (wl WordLibrary) ToWubiTable() WubiTable {
	wt := make(WubiTable, 0, len(wl))
	for _, entry := range wl {
		wt = append(wt, &WubiEntry{entry.GetWord(), entry.GetCode(), entry.GetPos()})
	}
	return wt
}

func (wl WordLibrary) ToCodeTable() CodeTable {
	t := wl.ToWubiTable()
	sort.Slice(t, func(i, j int) bool {
		return t[i].Pos < t[j].Pos
	})
	codeMap := make(map[string][]string)
	for _, wc := range t {
		if _, ok := codeMap[wc.Code]; !ok {
			codeMap[wc.Code] = []string{wc.Word}
			continue
		}
		codeMap[wc.Code] = append(codeMap[wc.Code], wc.Word)
	}

	ct := make(CodeTable, 0, len(codeMap))
	for k, v := range codeMap {
		ct = append(ct, &CodeEntry{k, v})
	}
	sort.Slice(ct, func(i, j int) bool {
		return ct[i].Code < ct[j].Code
	})
	return ct
}

func (wl WordLibrary) ToPinyinTable() PinyinTable {
	pyt := make(PinyinTable, 0, len(wl))
	for _, entry := range wl {
		pyt = append(pyt, &PinyinEntry{entry.GetWord(), entry.GetPinyin(), 1})
	}
	return pyt
}
