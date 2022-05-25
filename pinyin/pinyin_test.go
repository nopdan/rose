package pinyin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func Test(t *testing.T) {

	format := "baidu_bdict"
	data := Parse(format, "test/baidu.bdict")
	write(format, data)

	format = "baidu_bcd"
	data = Parse(format, "test/baidu.bcd")
	write(format, data)

	format = "sougou_scel"
	data = Parse(format, "test/sougou.scel")
	write(format, data)

	format = "qq_qcel"
	data = Parse(format, "test/含英文.qcel")
	write("qcel含英文", data)

	format = "ziguang_uwl"
	data = Parse(format, "test/music.uwl")
	write(format, data)

	format = "qq_qpyd"
	data = Parse(format, "test/qq.qpyd")
	write(format, data)
}

func write(filename string, data []Pinyin) {
	var buf bytes.Buffer
	for _, v := range data {
		buf.WriteString(v.Word + "\t" + strings.Join(v.Code, "'") + "\t" + strconv.Itoa(v.Freq) + "\r\n")
	}
	ioutil.WriteFile(fmt.Sprintf("out/%s.txt", filename), buf.Bytes(), 0777)
}
