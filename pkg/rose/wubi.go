package rose

import (
	"bufio"
	"fmt"
	"strings"
)

type Wubi struct {
	Dict
	WordFirst bool
	Encoding  string
}

func NewWubi(format string) *Wubi {
	d := new(Wubi)
	if format == "dd" {
		d.Name = "多多.txt"
		d.WordFirst = true
		d.Encoding = "UTF-8"
	} else if format == "bl" {
		d.Name = "冰凌.txt"
		d.WordFirst = false
		d.Encoding = "UTF-16LE"
	}
	return d
}

func (d *Wubi) Parse() {
	wl := make([]Entry, 0, d.size>>8)

	scan := bufio.NewScanner(d.rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		word, code := entry[0], entry[1]
		if !d.WordFirst {
			word, code = code, word
		}
		if strings.HasPrefix(word, "$ddcmd") {
			fmt.Println("多多的命令" + word)
			continue
		}
		wl = append(wl, &WubiEntry{word, code, 1})
	}
	d.WordLibrary = wl
}

func (d *Wubi) GenFrom(wl WordLibrary) []byte {
	var sb strings.Builder
	sb.Grow(len(wl))
	for _, v := range wl {
		if d.WordFirst {
			sb.WriteString(v.GetWord())
			sb.WriteByte('\t')
			sb.WriteString(v.GetCode())
		} else {
			sb.WriteString(v.GetCode())
			sb.WriteByte('\t')
			sb.WriteString(v.GetWord())
		}
		sb.WriteByte('\r')
		sb.WriteByte('\n')
	}
	ret := EncodeMust(sb.String(), d.Encoding)
	if d.Encoding == "UTF-16LE" {
		ret = append([]byte{0xff, 0xfe}, ret...)
	}
	return ret
}
