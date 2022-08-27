package zici

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	. "github.com/cxcn/dtool/pkg/util"
)

func ParseDuoduo(filename string) WcTable {
	return parseWcTable(filename, true)
}
func GenDuoduo(wct WcTable) []byte {
	return genWcTable(wct, true)
}

func ParseBingling(filename string) WcTable {
	return parseWcTable(filename, false)
}
func GenBingling(wct WcTable) []byte {
	return genWcTable(wct, false)
}

func parseWcTable(filename string, word_first bool) WcTable {
	f, _ := os.Open(filename)
	defer f.Close()
	rd, err := DecodeIO(f)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make(WcTable, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		word, code := entry[0], entry[1]
		if !word_first {
			word, code = code, word
		}
		if strings.HasPrefix(word, "$ddcmd") {
			fmt.Println("多多的命令" + word)
			continue
		}
		ret = append(ret, WordCode{word, code})
	}
	return ret
}

func genWcTable(wct WcTable, word_first bool) []byte {
	var buf bytes.Buffer
	for _, v := range wct {
		if word_first {
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
	return buf.Bytes()
}
