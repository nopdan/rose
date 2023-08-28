package encoder

import (
	"bufio"
	"strings"

	"github.com/nopdan/rose/util"
)

type Wubi struct {
	Char   map[rune]string
	IsAABC bool
}

func NewWubi(schema string, isAABC bool) *Wubi {
	r, err := util.Read("./data/CJK.txt")
	if err != nil {
		panic(err)
	}
	w := &Wubi{
		Char:   make(map[rune]string),
		IsAABC: isAABC,
	}
	idx := 0
	switch schema {
	case "wubi86", "86":
		idx = 2
	case "wubi98", "98":
		idx = 3
	case "wubi06", "06":
		idx = 4
	default:
		return nil
	}
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		fields := strings.Split(scan.Text(), ",")
		if idx > len(fields)-1 {
			continue
		}
		char := []rune(fields[1])[0]
		w.Char[char] = fields[idx]
	}
	return w
}

// 生成唯一编码
func (w *Wubi) Encode(word string) string {
	wordRunes := []rune(word)
	wordLen := len(wordRunes)
	var code string
	if wordLen == 1 {
		code = w.Char[wordRunes[0]]
	} else if wordLen == 2 {
		a := w.Char[wordRunes[0]]
		b := w.Char[wordRunes[1]]
		code = cut(a, 2) + cut(b, 2)
	} else if wordLen == 3 {
		a := w.Char[wordRunes[0]]
		b := w.Char[wordRunes[1]]
		c := w.Char[wordRunes[2]]
		if w.IsAABC {
			code = cut(a, 2) + cut(b, 1) + cut(c, 1)
		} else {
			code = cut(a, 1) + cut(b, 1) + cut(c, 2)
		}
	} else if wordLen >= 4 {
		a := w.Char[wordRunes[0]]
		b := w.Char[wordRunes[1]]
		c := w.Char[wordRunes[2]]
		z := w.Char[wordRunes[wordLen-1]]
		code = cut(a, 1) + cut(b, 1) + cut(c, 1) + cut(z, 1)
	}
	return code
}

func cut(code string, length int) string {
	if len(code) < length {
		return ""
	}
	return string(code[:length])
}
