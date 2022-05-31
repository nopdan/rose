package pinyin

import (
	"bufio"
	"bytes"
	"io"
	"log"

	encoder "github.com/cxcn/dtool/encoders"
	. "github.com/cxcn/dtool/utils"
)

func ParseWordOnly(rd io.Reader) []PyEntry {
	rd, err := DecodeIO(rd)
	if err != nil {
		log.Panic("编码格式未知")
	}
	ret := make([]PyEntry, 0, 0xff)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		ret = append(ret, PyEntry{scan.Text(), encoder.GetPinyin(scan.Text()), 1})
	}
	return ret
}

func GenWordOnly(pe []PyEntry) []byte {
	var buf bytes.Buffer
	for _, v := range pe {
		buf.WriteString(v.Word)
		buf.WriteString(LineBreak)
	}
	return buf.Bytes()
}
