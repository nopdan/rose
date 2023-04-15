package rose

import (
	"bufio"
	"bytes"
	"strings"
)

type Jidian struct{ Dict }

func NewJidian() *Jidian {
	d := new(Jidian)
	d.Name = "极点码表.txt"
	return d
}

func (d *Jidian) Parse() {
	wl := make([]Entry, 0, d.size>>8)

	scan := bufio.NewScanner(d.rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), " ")
		if len(entry) < 2 {
			continue
		}
		for i := 1; i < len(entry); i++ {
			wl = append(wl, &WubiEntry{entry[i], entry[0], i})
		}
	}
	d.WordLibrary = wl
}

func (Jidian) GenFrom(wl WordLibrary) []byte {
	ct := wl.ToCodeTable()
	var buf bytes.Buffer
	for _, v := range ct {
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Words, " "))
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}
