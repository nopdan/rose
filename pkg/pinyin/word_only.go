package pinyin

import (
	"bufio"
	"bytes"
	"log"

	"github.com/cxcn/dtool/pkg/encoder"
	"github.com/cxcn/dtool/pkg/util"
)

type WordOnly struct{}

func (WordOnly) Parse(filename string) Dict {
	rd, err := util.Read(filename)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make(Dict, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		ret = append(ret, Entry{scan.Text(), encoder.GetPinyin(scan.Text()), 1})
	}
	return ret
}

func (WordOnly) Gen(dict Dict) []byte {
	var buf bytes.Buffer
	for _, v := range dict {
		buf.WriteString(v.Word)
		buf.WriteString(util.LineBreak)
	}
	return buf.Bytes()
}
