package rose

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func outPyt(path string, dict *Dict) {
	os.Mkdir("out", 0666)
	os.WriteFile("out/"+filepath.Base(path)+".txt",
		Generate(dict.WordLibrary, "rime"), 0666)
}

func TestBdict(t *testing.T) {
	for _, path := range []string{
		"test/baidu.bdict", // 计算机硬件词汇 https://shurufa.baidu.com/dict
		"test/baidu.bcd",   // 计算机 https://mime.baidu.com/web/iw/index/
	} {
		dict := Parse(path, "bdict")
		outPyt(path, dict)
	}
}

func TestQpyd(t *testing.T) {
	path := "test/qq.qpyd"
	dict := Parse(path, "qpyd")
	outPyt(path, dict)
}

func TestScel(t *testing.T) {
	for _, path := range []string{
		"test/sogou.scel",
		"test/qq.qcel", // 网络流行新词【官方推荐】 http://cdict.qq.pinyin.cn/detail?dict_id=s4

		// "own/搜狗标准词库.scel",
	} {
		dict := Parse(path, "scel")
		outPyt(path, dict)
	}
}

func TestSogouBin(t *testing.T) {
	for _, path := range []string{
		"test/sogou_bak.bin", // 搜狗词库备份_2021_05_10.bin

		"test/sogou_bak_new.bin", // new
	} {
		dict := Parse(path, "sogou_bin")
		outPyt(path, dict)
	}
}

func TestUwl(t *testing.T) {
	for _, path := range []string{
		"test/music.uwl", // 来源，紫光内置

		// "own/sys7.uwl",
		// "own/大词库第六版.uwl",
	} {
		dict := Parse(path, "uwl")
		outPyt(path, dict)
	}
}

func TestMspyUDL(t *testing.T) {
	path := "test/ChsPinyinUDL.dat"
	dict := Parse(path, "mspy_udl")
	outPyt(path, dict)
}

// 拼音加加测试
func TestJiajia(t *testing.T) {
	path := "test/jiajia.txt"
	dict := Parse(path, "jj")
	outPyt(path, dict)
}

func TestWordOnly(t *testing.T) {
	path := "test/words.txt"
	dict := Parse(path, "w")
	outPyt(path, dict)
}

func TestPytOut(t *testing.T) {
	os.Mkdir("gen", 0666)
	wl := Parse("test/qq.qcel", "qcel").WordLibrary
	os.WriteFile("gen/msudppy.dat", Generate(wl, "msudp_dat"), 0666)
	os.WriteFile("gen/pyjj.txt", Generate(wl, "jj"), 0666)
	os.WriteFile("gen/word_only.txt", Generate(wl, "word_only"), 0666)
	os.WriteFile("gen/sogou.txt", Generate(wl, "sg"), 0666)
	os.WriteFile("gen/qq.txt", Generate(wl, "qq"), 0666)
	os.WriteFile("gen/baidu.txt", Generate(wl, "bd"), 0666)
	os.WriteFile("gen/google.txt", Generate(wl, "gg"), 0666)
	os.WriteFile("gen/rime.txt", Generate(wl, "rime"), 0666)
}

func tableOut(path string, dict *Dict) {
	os.Mkdir("out", 0666)
	os.WriteFile("out/"+filepath.Base(path)+".txt",
		Generate(dict.WordLibrary, "dd"), 0666)
}

func TestMsUDP(t *testing.T) {
	for _, path := range []string{
		"test/UserDefinedPhrase.dat",
		"test/ChsPinyinUDP.lex",
	} {
		dict := Parse(path, "msudp_dat")
		outPyt(path, dict)
	}
}

func TestMswbLex(t *testing.T) {
	path := "test/ChsWubi.lex"
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
	wl := Parse("test/duoduo.txt", "dd").WordLibrary
	os.WriteFile("gen/msudp.dat", Generate(wl, "msudp_dat"), 0666)
	os.WriteFile("gen/baidu.def", Generate(wl, "def"), 0666)
	os.WriteFile("gen/duoduo.txt", Generate(wl, "dd"), 0666)
	os.WriteFile("gen/bingling.txt", Generate(wl, "bl"), 0666)
}

func TestLexInitWordWeight(t *testing.T) {
	data, _ := os.ReadFile("assets/word_weight.txt")
	var buffer bytes.Buffer
	w := gzip.NewWriter(&buffer)
	w.Write(data)
	w.Flush()
	f, _ := os.OpenFile("assets/word_weight.bin", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	buffer.WriteTo(f)
}
