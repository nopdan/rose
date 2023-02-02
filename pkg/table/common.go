package table

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/imetool/goutil/util"
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
		ret = append(ret, Entry{word, code, 1})
	}
	return ret
}

func (c Common) Gen(table Table) []byte {
	var sb strings.Builder
	sb.Grow(len(table))
	for _, v := range table {
		if c.WordFirst {
			sb.WriteString(v.Word)
			sb.WriteByte('\t')
			sb.WriteString(v.Code)
		} else {
			sb.WriteString(v.Code)
			sb.WriteByte('\t')
			sb.WriteString(v.Word)
		}
		sb.WriteByte('\r')
		sb.WriteByte('\n')
	}
	ret, _ := util.Encode(sb.String(), c.Encoding)
	if c.Encoding == "UTF-16LE" {
		ret = append([]byte{0xff, 0xfe}, ret...)
	}
	return ret
}
