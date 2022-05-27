package zici

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func Test(t *testing.T) {

	// 哲哲豆词库 1w 多条
	format := "duoduo"
	data := Parse(format, "test/duoduo.txt")
	write(format, data)

	// // 哲哲豆词库
	// format = "baidu_def"
	// data = Parse(format, "test/baidu.def")
	// write(format, data)

	// 091 点儿词库
	format = "jidian_mb"
	data = Parse(format, "test/jidian.mb")
	write(format, data)

}

func write(filename string, data []ZcEntry) {
	var buf bytes.Buffer
	for _, v := range data {
		buf.WriteString(v.Word + "\t" + v.Code + "\r\n")
	}
	ioutil.WriteFile(fmt.Sprintf("out/%s.txt", filename), buf.Bytes(), 0777)
}
