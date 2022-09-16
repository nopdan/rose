package pinyin

import (
	"os"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) {

	// 计算机硬件词汇 https://shurufa.baidu.com/dict
	path := "test/baidu.bdict"
	dict := BaiduBdict{}.Parse(path)
	write_out(path, dict)

	// 计算机 https://mime.baidu.com/web/iw/index/
	path = "test/baidu.bcd"
	dict = BaiduBdict{}.Parse(path)
	write_out(path, dict)

	// 搜狗细胞词库
	path = "test/sogou.scel"
	dict = SogouScel{}.Parse(path)
	write_out(path, dict)

	// 网络流行新词【官方推荐】 http://cdict.qq.pinyin.cn/detail?dict_id=s4
	path = "test/qq.qcel"
	dict = SogouScel{}.Parse(path)
	write_out(path, dict)

	// 来源，紫光内置
	path = "test/music.uwl"
	dict = ZiguangUwl{}.Parse(path)
	write_out(path, dict)

	// QQ
	path = "test/qq.qpyd"
	dict = QqQpyd{}.Parse(path)
	write_out(path, dict)

	// 搜狗词库备份_2021_05_10.bin
	path = "test/sogou_bak.bin"
	dict = SogouBin{}.Parse(path)
	write_out(path, dict)
}

// 拼音加加测试
func TestJiajia(t *testing.T) {
	path := "own/拼音加加-305万大词库.txt"
	dict := JiaJia{}.Parse(path)
	write_out(path, dict)
}

func TestZiguangUwl(t *testing.T) {
	path := "own/sys7.uwl"
	dict := ZiguangUwl{}.Parse(path)
	write_out(path, dict)

	path = "own/大词库第六版.uwl"
	dict = ZiguangUwl{}.Parse(path)
	write_out(path, dict)
}

func TestSogouScel(t *testing.T) {
	path := "own/搜狗标准词库.scel"
	dict := SogouScel{}.Parse(path)
	write_out(path, dict)
}

func TestSogouBin(t *testing.T) {
	path := "own/搜狗词库备份_2021_05_10.bin"
	dict := SogouBin{}.Parse(path)
	write_out(path, dict)
}

func TestMspyUDL(t *testing.T) {
	path := "own/ChsPinyinUDL.dat"
	dict := MspyUDL{}.Parse(path)
	write_out(path, dict)
}

func write_out(path string, dict Dict) {
	os.WriteFile(filepath.Dir(path)+"/out_"+filepath.Base(path)+".txt",
		Common{'\t', '\'', "wcf"}.Gen(dict), 0666)
}

func TestGen(t *testing.T) {
	dict := SogouScel{}.Parse("test/qq.qcel")
	os.WriteFile("test/gen_pyjj.txt", JiaJia{}.Gen(dict), 0666)
	os.WriteFile("test/gen_word_only.txt", WordOnly{}.Gen(dict), 0666)
	os.WriteFile("test/gen_sogou.txt", Sogou.Gen(dict), 0666)
	os.WriteFile("test/gen_baidu.txt", Baidu.Gen(dict), 0666)
	os.WriteFile("test/gen_google.txt", Google.Gen(dict), 0666)
	os.WriteFile("test/gen_qq.txt", QQ.Gen(dict), 0666)
}
