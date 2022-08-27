package pinyin

import (
	"bufio"
	"bytes"
	"log"
	"os"

	"github.com/cxcn/dtool/pkg/encoder"
	. "github.com/cxcn/dtool/pkg/util"
)

func ParseWordOnly(filename string) WpfDict {
	f, _ := os.Open(filename)
	defer f.Close()
	rd, err := DecodeIO(f)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make(WpfDict, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		ret = append(ret, WordPyFreq{scan.Text(), encoder.GetPinyin(scan.Text()), 1})
	}
	return ret
}

func GenWordOnly(dict WpfDict) []byte {
	var buf bytes.Buffer
	for _, v := range dict {
		buf.WriteString(v.Word)
		buf.WriteString(LineBreak)
	}
	return buf.Bytes()
}
