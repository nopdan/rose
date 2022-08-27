package zici

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"

	. "github.com/cxcn/dtool/pkg/util"
)

func ParseJidian(filename string) WcTable {
	f, _ := os.Open(filename)
	defer f.Close()
	rd, err := DecodeIO(f)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make(WcTable, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), " ")
		if len(entry) < 2 {
			continue
		}
		for i := 1; i < len(entry); i++ {
			ret = append(ret, WordCode{entry[i], entry[0]})
		}
	}
	return ret
}

func GenJidian(wct WcTable) []byte {
	cwt := ToCwsTable(wct)
	var buf bytes.Buffer
	for _, v := range cwt {
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Words, " "))
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}
