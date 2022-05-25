package zici

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func ParseBingling(rd io.Reader) Dict {
	ret := make(Dict, 1e5)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		code, word := entry[0], entry[1]
		if _, ok := ret[code]; !ok {
			ret[code] = []string{word}
			continue
		}
		ret[code] = append(ret[code], word)
	}
	return ret
}

func GenBingling(dl []codeAndWords) []byte {
	var buf bytes.Buffer
	for _, v := range dl {
		for _, word := range v.words {
			buf.WriteString(v.code)
			buf.WriteByte('\t')
			buf.WriteString(word)
			buf.WriteByte('\r')
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
}
