package pinyin

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/cxcn/dtool/utils"
)

func Test(t *testing.T) {

	// 计算机硬件词汇 https://shurufa.baidu.com/dict
	filename := "test/baidu.bdict"
	data := ParseBaiduBdict(filename)
	write_out(filename, data)

	// 计算机 https://mime.baidu.com/web/iw/index/
	filename = "test/baidu.bcd"
	data = ParseBaiduBdict(filename)
	write_out(filename, data)

	// 搜狗细胞词库
	filename = "test/sogou.scel"
	data = ParseSogouScel(filename)
	write_out(filename, data)

	// 网络流行新词【官方推荐】 http://cdict.qq.pinyin.cn/detail?dict_id=s4
	filename = "test/qq.qcel"
	data = ParseSogouScel(filename)
	write_out(filename, data)

	// 来源，紫光内置
	filename = "test/music.uwl"
	data = ParseZiguangUwl(filename)
	write_out(filename, data)

	// QQ
	filename = "test/qq.qpyd"
	data = ParseQqQpyd(filename)
	write_out(filename, data)

	// 搜狗词库备份_2021_05_10.bin
	filename = "test/sogou_bak.bin"
	data = ParseSogouBin(filename)
	write_out(filename, data)
}

// 拼音加加测试
func TestJiajia(t *testing.T) {
	filename := "own/拼音加加-305万大词库.txt"
	data := ParseJiaJia(filename)
	write_out(filename, data)
}

func TestZiguangUwl(t *testing.T) {
	filename := "own/sys7.uwl"
	data := ParseZiguangUwl(filename)
	write_out(filename, data)

	filename = "own/大词库第六版.uwl"
	data = ParseZiguangUwl(filename)
	write_out(filename, data)
}

func TestSogouScel(t *testing.T) {
	filename := "own/搜狗标准词库.scel"
	data := ParseSogouScel(filename)
	write_out(filename, data)
}

func TestSogouBin(t *testing.T) {
	filename := "own/搜狗词库备份_2021_05_10.bin"
	data := ParseSogouBin(filename)
	write_out(filename, data)
}

func TestMspyUDL(t *testing.T) {
	filename := "own/ChsPinyinUDL.dat"
	data := ParseMspyUDL(filename)
	write_out(filename, data)

	filename = "own/SuperRime拓展词库 for Win10拼音版(600万词-含BetterRime)-v20.3.dat"
	data = ParseMspyUDL(filename)
	write_out(filename, data)
}

func write_out(filename string, data WpfDict) {
	err := ioutil.WriteFile(fmt.Sprintf("%s_out.txt", filename),
		GenPinyin(data, PinyinRule{'\t', '\'', "wcf"}), 0777)
	if err != nil {
		println(err)
	}
}

func TestGen(t *testing.T) {
	data := ParseZiguangUwl("own/sys.uwl")
	ioutil.WriteFile("own/sys_pyjj.txt", GenJiaJia(data), 0777)
	ioutil.WriteFile("own/sys_sogou.txt", GenPinyin(data, R_sogou), 0777)
	ioutil.WriteFile("own/sys_baidu.txt", GenPinyin(data, R_baidu), 0777)
	ioutil.WriteFile("own/sys_google.txt", GenPinyin(data, R_google), 0777)
	ioutil.WriteFile("own/sys_qq.txt", GenPinyin(data, R_qq), 0777)
	ioutil.WriteFile("own/sys_word_only.txt", GenPinyin(data, R_word_only), 0777)
}
