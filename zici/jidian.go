package zici

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func ParseJidian(rd io.Reader) []ZcEntry {
	ret := make([]ZcEntry, 0, 1e5)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		for i := 1; i < len(entry); i++ {
			ret = append(ret, ZcEntry{entry[i], entry[0]})
		}
	}
	return ret
}

func GenJidian(dl []CodeEntry) []byte {
	var buf bytes.Buffer
	for _, v := range dl {
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Words, " "))
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}
