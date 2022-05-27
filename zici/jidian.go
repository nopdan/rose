package zici

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func ParseJidian(rd io.Reader) []CodeEntry {
	ret := make([]CodeEntry, 0, 1e5)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), " ")
		if len(entry) < 2 {
			continue
		}
		ret = append(ret, CodeEntry{entry[0], entry[1:]})
	}
	return ret
}

func GenJidian(ce []CodeEntry) []byte {
	var buf bytes.Buffer
	for _, v := range ce {
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Words, " "))
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}
