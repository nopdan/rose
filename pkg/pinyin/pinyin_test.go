package pinyin

import (
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {

	// 计算机硬件词汇 https://shurufa.baidu.com/dict
	filename := "test/baidu.bdict"
	data := BaiduBdict{}.Parse(filename)
	write_out(filename, data)

	// 计算机 https://mime.baidu.com/web/iw/index/
	filename = "test/baidu.bcd"
	data = BaiduBdict{}.Parse(filename)
	write_out(filename, data)

	// 搜狗细胞词库
	filename = "test/sogou.scel"
	data = SogouScel{}.Parse(filename)
	write_out(filename, data)

	// 网络流行新词【官方推荐】 http://cdict.qq.pinyin.cn/detail?dict_id=s4
	filename = "test/qq.qcel"
	data = SogouScel{}.Parse(filename)
	write_out(filename, data)

	// 来源，紫光内置
	filename = "test/music.uwl"
	data = ZiguangUwl{}.Parse(filename)
	write_out(filename, data)

	// QQ
	filename = "test/qq.qpyd"
	data = QqQpyd{}.Parse(filename)
	write_out(filename, data)

	// 搜狗词库备份_2021_05_10.bin
	filename = "test/sogou_bak.bin"
	data = SogouBin{}.Parse(filename)
	write_out(filename, data)
}

// 拼音加加测试
func TestJiajia(t *testing.T) {
	filename := "own/拼音加加-305万大词库.txt"
	data := JiaJia{}.Parse(filename)
	write_out(filename, data)
}

func TestZiguangUwl(t *testing.T) {
	filename := "own/sys7.uwl"
	data := ZiguangUwl{}.Parse(filename)
	write_out(filename, data)

	filename = "own/大词库第六版.uwl"
	data = ZiguangUwl{}.Parse(filename)
	write_out(filename, data)
}

func TestSogouScel(t *testing.T) {
	filename := "own/搜狗标准词库.scel"
	data := SogouScel{}.Parse(filename)
	write_out(filename, data)
}

func TestSogouBin(t *testing.T) {
	filename := "own/搜狗词库备份_2021_05_10.bin"
	data := SogouBin{}.Parse(filename)
	write_out(filename, data)
}

func TestMspyUDL(t *testing.T) {
	filename := "own/ChsPinyinUDL.dat"
	data := MspyDat{}.Parse(filename)
	write_out(filename, data)

	filename = "own/SuperRime拓展词库 for Win10拼音版(600万词-含BetterRime)-v20.3.dat"
	data = MspyDat{}.Parse(filename)
	write_out(filename, data)
}

func write_out(filename string, data Dict) {
	err := os.WriteFile(fmt.Sprintf("%s_out.txt", filename),
		Common{'\t', '\'', "wcf"}.Gen(data), 0666)
	if err != nil {
		println(err)
	}
}

func TestGen(t *testing.T) {
	data := ZiguangUwl{}.Parse("own/sys.uwl")
	os.WriteFile("own/sys_pyjj.txt", JiaJia{}.Gen(data), 0666)
	os.WriteFile("own/sys_sogou.txt", TxtSogou.Gen(data), 0666)
	os.WriteFile("own/sys_baidu.txt", TxtBaidu.Gen(data), 0666)
	os.WriteFile("own/sys_google.txt", TxtGoogle.Gen(data), 0666)
	os.WriteFile("own/sys_qq.txt", TxtQQ.Gen(data), 0666)
	os.WriteFile("own/sys_word_only.txt", TxtWordOnly.Gen(data), 0666)
}
