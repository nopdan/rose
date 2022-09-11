package table

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/cxcn/dtool/pkg/util"
)

type Common struct {
	WordFirst bool
	Encoding  string
}

var (
	DuoDuo   = Common{true, "UTF-8"}
	Bingling = Common{false, "UTF-16LE"}
)

func (c Common) Parse(filename string) Table {
	rd, err := util.Read(filename)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make(Table, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		word, code := entry[0], entry[1]
		if !c.WordFirst {
			word, code = code, word
		}
		if strings.HasPrefix(word, "$ddcmd") {
			fmt.Println("多多的命令" + word)
			continue
		}
		ret = append(ret, Entry{word, code})
	}
	return ret
}

func (c Common) Gen(table Table) []byte {
	var buf bytes.Buffer
	for _, v := range table {
		if c.WordFirst {
			buf.WriteString(v.Word)
			buf.WriteByte('\t')
			buf.WriteString(v.Code)
		} else {
			buf.WriteString(v.Code)
			buf.WriteByte('\t')
			buf.WriteString(v.Word)
		}
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
	ret, _ := util.Encode(buf.Bytes(), c.Encoding)
	if c.Encoding == "UTF-16LE" {
		ret = append([]byte{0xff, 0xfe}, ret...)
	}
	return ret
}
