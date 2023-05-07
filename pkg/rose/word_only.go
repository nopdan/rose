package rose

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/flowerime/rose/pkg/pinyin"
	"github.com/nopdan/ku"
)

// 纯词词库
type WordEntry struct{ Word string }

func (e *WordEntry) GetWord() string     { return e.Word }
func (e *WordEntry) GetCode() string     { return strings.Join(e.GetPinyin(), "") }
func (e *WordEntry) GetPos() int         { return 1 }
func (e *WordEntry) GetPinyin() []string { return pinyin.Match(e.Word) }
func (e *WordEntry) GetFreq() int        { return 1 }

type WordOnly struct{ Dict }

func NewWordOnly() *WordOnly {
	d := new(WordOnly)
	d.Name = "纯汉字.txt"
	return d
}

func (d *WordOnly) Parse() {
	wl := make([]Entry, 0, 0xff)
	scan := bufio.NewScanner(d.rd)
	for scan.Scan() {
		wl = append(wl, &WordEntry{scan.Text()})
	}
	d.WordLibrary = wl
}

func (WordOnly) GenFrom(wl WordLibrary) []byte {
	var buf bytes.Buffer
	for _, entry := range wl {
		buf.WriteString(entry.GetWord())
		buf.WriteString(ku.LineBreak)
	}
	return buf.Bytes()
}
