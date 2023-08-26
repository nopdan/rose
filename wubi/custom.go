package wubi

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"slices"

	"golang.org/x/net/html/charset"
)

type Custom struct {
	Template
	sep      string
	rule     []string
	encoding string
}

// 规则用"|"分割，第一项是分隔符
//
// t: tab, s: space; w: word, c: code, r: rank
//
// 多多 t|w|c，冰凌 t|c|w
func NewCustom(rule string, e string) *Custom {
	f := new(Custom)
	f.encoding = e
	s := strings.Split(rule, "|")
	switch s[0] {
	case "t":
		f.sep = "\t"
	case "s":
		f.sep = " "
	default:
		f.sep = s[0]
	}
	f.rule = s[1:]
	return f
}

func NewDuoduo() *Custom {
	f := NewCustom("t|w|c", "UTF-16LE")
	f.Name = "多多.txt"
	f.ID = "duoduo"
	return f
}
func NewBingling() *Custom {
	f := NewCustom("t|c|w", "UTF-8")
	f.Name = "冰凌.txt"
	f.ID = "bingling"
	return f
}

func (f *Custom) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)

	scan := bufio.NewScanner(r)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), f.sep)
		var word, code string
		var rank int
		for i := range f.rule {
			if i >= len(entry) {
				break
			}
			switch f.rule[i] {
			case "w":
				word = entry[i]
			case "c":
				code = entry[i]
			case "r":
				rank, _ = strconv.Atoi(entry[i])
			}
		}
		di = append(di, &Entry{word, code, rank})
	}
	if slices.Contains(f.rule, "r") {
		f.Rank = true
	}
	return di
}

func (f *Custom) Marshal(di []*Entry, hasRank bool) []byte {
	var buf bytes.Buffer
	buf.Grow(len(di))
	// 生成 Rank
	if slices.Contains(f.rule, "r") && !hasRank {
		di = GenRank(di)
	}

	e, name := charset.Lookup(f.encoding)
	// bom
	switch name {
	case "utf-16le":
		buf.Write([]byte{0xff, 0xfe})
	case "utf-16be":
		buf.Write([]byte{0xfe, 0xff})
	}
	w := e.NewEncoder().Writer(&buf)

	for _, v := range di {
		for i := range f.rule {
			switch f.rule[i] {
			case "w":
				w.Write([]byte(v.Word))
			case "c":
				w.Write([]byte(v.Code))
			case "r":
				w.Write([]byte(strconv.Itoa(v.Rank)))
			}
			if i != len(f.rule)-1 {
				w.Write([]byte(f.sep))
			}
		}
		w.Write([]byte{'\r', '\n'})
	}
	return buf.Bytes()
}
