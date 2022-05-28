package pinyin

import (
	"bytes"
	"strconv"
	"strings"
)

func GenGenersal(pe []PyEntry) []byte {
	var buf bytes.Buffer
	for _, v := range pe {
		buf.WriteString(v.Word)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Codes, "'"))
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(v.Freq))
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}
