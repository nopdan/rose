package table

import (
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	// 哲哲豆词库 1w 多条
	filename := "test/duoduo.txt"
	wct := DuoDuo.Parse(filename)
	write_out(filename, wct)

	// 091 点儿词库
	filename = "test/jidian.mb"
	wct = JidianMb{}.Parse(filename)
	write_out(filename, wct)
}

func TestBaiduDef(t *testing.T) {
	// 哲哲豆词库
	filename := "own/baidu.def"
	wct := BaiduDef{}.Parse(filename)
	write_out(filename, wct)
}

func write_out(filename string, table Table) {
	os.WriteFile(fmt.Sprintf("%s_out.txt", filename), DuoDuo.Gen(table), 0)
}
