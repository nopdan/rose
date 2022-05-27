package zici

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func ParseDuoduo(rd io.Reader) []ZcEntry {
	ret := make([]ZcEntry, 0, 1e5)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		word, code := entry[0], entry[1]
		if strings.HasPrefix(word, "$ddcmd") {
			fmt.Println(code)
			continue
		}
		ret = append(ret, ZcEntry{word, code})
	}
	return ret
}

func GenDuoduo(ce []CodeEntry) []byte {
	var buf bytes.Buffer
	for _, v := range ce {
		for _, word := range v.Words {
			buf.WriteString(word)
			buf.WriteByte('\t')
			buf.WriteString(v.Code)
			buf.WriteByte('\r')
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
}
