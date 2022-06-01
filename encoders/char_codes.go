package encoder

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"

	. "github.com/cxcn/dtool/utils"
)

// 读取单字码表
func ReadCharCodes(data []byte) CharCodes {
	ret := make(CharCodes)
	rd := bytes.NewReader(data)
	scan := bufio.NewScanner(rd)

	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		ret[[]rune(entry[0])[0]] = entry[1:]
	}
	return ret
}
