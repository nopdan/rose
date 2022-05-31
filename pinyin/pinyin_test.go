package pinyin

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test(t *testing.T) {

	// 计算机硬件词汇 https://shurufa.baidu.com/dict
	format := "baidu_bdict"
	data := Parse(format, "test/baidu.bdict")
	write("out/"+format, data)

	// 计算机 https://mime.baidu.com/web/iw/index/
	format = "baidu_bcd"
	data = Parse(format, "test/baidu.bcd")
	write("out/"+format, data)

	// 搜狗细胞词库
	format = "sogou_scel"
	data = Parse(format, "test/sogou.scel")
	write("out/"+format, data)

	// 网络流行新词【官方推荐】 http://cdict.qq.pinyin.cn/detail?dict_id=s4
	format = "qq_qcel"
	data = Parse(format, "test/qq.qcel")
	write("out/"+format, data)

	// 来源，紫光内置
	format = "ziguang_uwl"
	data = Parse(format, "test/music.uwl")
	write("out/"+format, data)

	// QQ
	format = "qq_qpyd"
	data = Parse(format, "test/qq.qpyd")
	write("out/"+format, data)

	// 搜狗词库备份_2021_05_10.bin
	format = "sogou_bin"
	data = Parse(format, "test/sogou_bak.bin")
	write("out/"+format, data)

}

// 拼音加加测试
func TestPyjj(t *testing.T) {
	format := "pyjj"
	fp := "own/拼音加加-305万大词库.txt"
	data := Parse(format, fp)
	write(fp, data)
}

func TestZgUwl(t *testing.T) {
	format := "ziguang_uwl"
	fp := "own/sys7.uwl"
	data := Parse(format, fp)
	write(fp, data)

	fp = "own/大词库第六版.uwl"
	data = Parse(format, fp)
	write(fp, data)
}

func TestSgScel(t *testing.T) {
	format := "sogou_scel"
	fp := "own/成语俗语大全.qcel"
	data := Parse(format, fp)
	write(fp, data)
}

func TestSogouBin(t *testing.T) {
	format := "sogou_bin"
	fp := "own/搜狗词库备份_2021_05_10.bin"
	data := Parse(format, fp)
	write(fp, data)
}

func TestMsUDL(t *testing.T) {
	format := "mspy_dat"
	fp := "own/ChsPinyinUDL.dat"
	data := Parse(format, fp)
	write(fp, data)

	fp = "own/SuperRime拓展词库 for Win10拼音版(600万词-含BetterRime)-v20.3.dat"
	data = Parse(format, fp)
	write(fp, data)
}

func write(filename string, data []PyEntry) {
	err := ioutil.WriteFile(fmt.Sprintf("%s.txt", filename), GenGeneral(data, baidu), 0777)
	if err != nil {
		println(err)
	}
}

func TestGen(t *testing.T) {
	format := "ziguang_uwl"
	fp := "own/sys.uwl"
	data := Parse(format, fp)
	ioutil.WriteFile("own/sys_pyjj.txt", GenPyJiaJia(data), 0777)
	ioutil.WriteFile("own/sys_sogou.txt", Gen("sogou", data), 0777)
	ioutil.WriteFile("own/sys_baidu.txt", Gen("baidu", data), 0777)
	ioutil.WriteFile("own/sys_google.txt", Gen("google", data), 0777)
	ioutil.WriteFile("own/sys_qq.txt", Gen("qq", data), 0777)
	ioutil.WriteFile("own/sys_word_only.txt", Gen("word_only", data), 0777)
}
