package table

import (
	"bufio"
	"bytes"
	"log"
	"sort"
	"strings"

	"github.com/imetool/goutil/util"
)

// 极点形式码表
type JdTable []CodeWords

// 一码多词
type CodeWords struct {
	Code  string
	Words []string
}

type Jidian struct{}

func (Jidian) Parse(filename string) Table {
	rd, err := util.Read(filename)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make(Table, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), " ")
		if len(entry) < 2 {
			continue
		}
		for i := 1; i < len(entry); i++ {
			ret = append(ret, Entry{entry[i], entry[0], byte(i)})
		}
	}
	return ret
}

func (Jidian) Gen(table Table) []byte {
	jdt := ToJdTable(table)
	var buf bytes.Buffer
	for _, v := range jdt {
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Words, " "))
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

// 转为 极点形式，一码多词
func ToJdTable(table Table) JdTable {
	codeMap := make(map[string][]string)
	for _, wc := range table {
		if _, ok := codeMap[wc.Code]; !ok {
			codeMap[wc.Code] = []string{wc.Word}
			continue
		}
		codeMap[wc.Code] = append(codeMap[wc.Code], wc.Word)
	}
	ret := make(JdTable, 0, len(codeMap))
	for k, v := range codeMap {
		ret = append(ret, CodeWords{k, v})
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Code < ret[j].Code
	})
	return ret
}
