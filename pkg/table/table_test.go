package table

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMswbLex(t *testing.T) {
	path := "own/ChsWubiNew.lex"
	dict := MswbLex{}.Parse(path)
	write_out(path, dict)
}

func TestMsUDP(t *testing.T) {
	path := "own/UserDefinedPhrase.dat"
	dict := MsUDP{}.Parse(path)
	write_out(path, dict)
}

func TestJidian(t *testing.T) {
	// 091 点儿词库
	path := "test/jidian.mb"
	table := JidianMb{}.Parse(path)
	write_out(path, table)
}

func TestBaiduDef(t *testing.T) {
	// 哲哲豆词库
	path := "own/baidu.def"
	table := BaiduDef{}.Parse(path)
	write_out(path, table)
}

func TestFcitx4Mb(t *testing.T) {
	// 98 五笔
	path := "own/98wb_ci.mb"
	table := Fcitx4Mb{}.Parse(path)
	write_out(path, table)
}

func TestPang(t *testing.T) {
	path := "own/binary.main.bin"
	table := Pang{}.Parse(path)
	// write_out(path, table)
	os.WriteFile(filepath.Dir(path)+"/out_"+filepath.Base(path)+".txt",
		Jidian{}.Gen(table), 0666)
}

func TestPangAssoc(t *testing.T) {
	path := "own/binary.assoc.bin"
	table := PangAssoc{}.Parse(path)
	// write_out(path, table)
	os.WriteFile(filepath.Dir(path)+"/out_"+filepath.Base(path)+".txt",
		DuoDuo.Gen(table), 0666)
}

func write_out(path string, table Table) {
	os.WriteFile(filepath.Dir(path)+"/out_"+filepath.Base(path)+".txt",
		DuoDuo.Gen(table), 0666)
}

func TestGen(t *testing.T) {
	// 哲哲豆词库 1w 多条
	table := DuoDuo.Parse("test/duoduo.txt")
	os.WriteFile("test/gen_duoduo.txt", DuoDuo.Gen(table), 0666)
	os.WriteFile("test/gen_bingling.txt", Bingling.Gen(table), 0666)
	os.WriteFile("test/gen_baidu.def", BaiduDef{}.Gen(table), 0666)
	os.WriteFile("test/gen_msudp.dat", MsUDP{}.Gen(table), 0666)
}
