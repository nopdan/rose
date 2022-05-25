package zici

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func ParseJidian(rd io.Reader) Dict {
	ret := make(Dict, 1e5)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		code, words := entry[0], strings.Split(entry[1], " ")
		ret[code] = words
	}
	return ret
}

func GenJidian(dl []codeAndWords) []byte {
	var buf bytes.Buffer
	for _, v := range dl {
		buf.WriteString(v.code)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.words, " "))
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}
