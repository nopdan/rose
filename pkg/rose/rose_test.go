package rose

import (
	"os"
	"path/filepath"
	"testing"
)

func outPyt(path string, dict *Dict) {
	os.Mkdir("out", 0666)
	os.WriteFile("out/"+filepath.Base(path)+".txt",
		Generate(dict, "rime"), 0666)
}

func TestBdict(t *testing.T) {
	// 计算机硬件词汇 https://shurufa.baidu.com/dict
	path := "test/baidu.bdict"
	dict := Parse(path, "bdict")
	outPyt(path, dict)

	// 计算机 https://mime.baidu.com/web/iw/index/
	path = "test/baidu.bcd"
	dict = Parse(path, "bcd")
	outPyt(path, dict)
}

func TestQpyd(t *testing.T) {
	path := "test/qq.qpyd"
	dict := Parse(path, "qpyd")
	outPyt(path, dict)
}

func TestScel(t *testing.T) {
	path := "test/sogou.scel"
	dict := Parse(path, "scel")
	outPyt(path, dict)

	path = "own/搜狗标准词库.scel"
	dict = Parse(path, "scel")
	outPyt(path, dict)

	// 网络流行新词【官方推荐】 http://cdict.qq.pinyin.cn/detail?dict_id=s4
	path = "test/qq.qcel"
	dict = Parse(path, "qcel")
	outPyt(path, dict)
}

func TestSogouBin(t *testing.T) {
	// 搜狗词库备份_2021_05_10.bin
	path := "test/sogou_bak.bin"
	dict := Parse(path, "sogou_bin")
	outPyt(path, dict)

	// new
	path = "own/sogou-bin/搜狗词库备份_2023_4_5.bin"
	dict = Parse(path, "sogou_bin")
	outPyt(path, dict)
}

func TestUwl(t *testing.T) {
	// 来源，紫光内置
	path := "test/music.uwl"
	dict := Parse(path, "uwl")
	outPyt(path, dict)

	path = "own/sys7.uwl"
	dict = Parse(path, "uwl")
	outPyt(path, dict)

	path = "own/大词库第六版.uwl"
	dict = Parse(path, "uwl")
	outPyt(path, dict)
}

func TestMspyUDL(t *testing.T) {
	path := "own/ChsPinyinUDL.dat"
	dict := Parse(path, "mspy_udl")
	outPyt(path, dict)
}

func TestMsUDP(t *testing.T) {
	path := "own/UserDefinedPhrase.dat"
	dict := Parse(path, "msudp_dat")
	outPyt(path, dict)
}

// 拼音加加测试
func TestJiajia(t *testing.T) {
	path := "own/拼音加加-305万大词库.txt"
	dict := Parse(path, "jj")
	outPyt(path, dict)
}

func TestPytOut(t *testing.T) {
	os.Mkdir("gen", 0666)
	dict := Parse("test/qq.qcel", "qcel")
	os.WriteFile("gen/pyjj.txt", Generate(dict, "jj"), 0666)
	os.WriteFile("gen/word_only.txt", Generate(dict, "word_only"), 0666)
	os.WriteFile("gen/sogou.txt", Generate(dict, "sg"), 0666)
	os.WriteFile("gen/qq.txt", Generate(dict, "qq"), 0666)
	os.WriteFile("gen/baidu.txt", Generate(dict, "bd"), 0666)
	os.WriteFile("gen/google.txt", Generate(dict, "gg"), 0666)
	os.WriteFile("gen/rime.txt", Generate(dict, "rime"), 0666)
}

func tableOut(path string, dict *Dict) {
	os.Mkdir("out", 0666)
	os.WriteFile("out/"+filepath.Base(path)+".txt",
		Generate(dict, "dd"), 0666)
}

func TestMswbLex(t *testing.T) {
	path := "own/ChsWubiNew.lex"
	table := Parse(path, "mswb_lex")
	tableOut(path, table)
}

func TestJidian(t *testing.T) {
	// 091 点儿词库
	path := "test/jidian.mb"
	table := Parse(path, "jidian_mb")
	tableOut(path, table)
}

func TestBaiduDef(t *testing.T) {
	// 哲哲豆词库
	path := "own/baidu.def"
	table := Parse(path, "baidu_def")
	tableOut(path, table)
}

func TestFcitx4Mb(t *testing.T) {
	// 98 五笔
	path := "own/98wb_ci.mb"
	table := Parse(path, "fcitx4_mb")
	tableOut(path, table)
}

func TestGen(t *testing.T) {
	os.Mkdir("gen", 0666)
	// 哲哲豆词库 1w 多条
	table := Parse("test/duoduo.txt", "dd")
	os.WriteFile("gen/msudp.dat", Generate(table, "msudp_dat"), 0666)
	os.WriteFile("gen/baidu.def", Generate(table, "def"), 0666)
	os.WriteFile("gen/duoduo.txt", Generate(table, "dd"), 0666)
	os.WriteFile("gen/bingling.txt", Generate(table, "bl"), 0666)
}
