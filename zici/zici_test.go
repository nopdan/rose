package zici

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/cxcn/dtool/utils"
)

func Test(t *testing.T) {
	// 哲哲豆词库 1w 多条
	filename := "test/duoduo.txt"
	wct := ParseDuoduo(filename)
	write_out(filename, wct)

	// 091 点儿词库
	filename = "test/jidian.mb"
	wct = ParseJidianMb(filename)
	write_out(filename, wct)
}

func TestBaiduDef(t *testing.T) {
	// 哲哲豆词库
	filename := "own/baidu.def"
	wct := ParseBaiduDef(filename)
	write_out(filename, wct)
}

func write_out(filename string, wct WcTable) {
	ioutil.WriteFile(fmt.Sprintf("%s_out.txt", filename), GenDuoduo(wct), 0777)
}
