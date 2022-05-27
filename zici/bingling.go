package zici

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func ParseBingling(rd io.Reader) []ZcEntry {
	ret := make([]ZcEntry, 0, 1e5)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		word, code := entry[1], entry[0]
		ret = append(ret, ZcEntry{word, code})
	}
	return ret
}

func GenBingling(ce []CodeEntry) []byte {
	var buf bytes.Buffer
	for _, v := range ce {
		for _, word := range v.Words {
			buf.WriteString(v.Code)
			buf.WriteByte('\t')
			buf.WriteString(word)
			buf.WriteByte('\r')
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
}
