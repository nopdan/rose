package pinyin

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
}

func TestOwn(t *testing.T) {

	format := "ziguang_uwl"
	fp := "own/sys.uwl"
	data := Parse(format, fp)
	write("own/"+"sys.uwl", data)

	format = "ziguang_uwl"
	fp = "own/大词库第六版.uwl"
	data = Parse(format, fp)
	write(fp, data)

	format = "sogou_scel"
	fp = "own/成语俗语大全.qcel"
	data = Parse(format, fp)
	write(fp, data)

	// format = "mspy_dat"
	// fp = "own/SuperRime拓展词库 for Win10拼音版(600万词-含BetterRime)-v20.3.dat"
	// data = Parse(format, fp)
	// write(fp, data)

}

func TestUDL(t *testing.T) {
	f, _ := os.Open("own/ChsPinyinUDL.dat")
	data := ParseMspyUDL(f)
	ioutil.WriteFile("own/ChsPinyinUDL.dat.txt", []byte(strings.Join(data, "\n")), 0777)
}

func TestGenernal(t *testing.T) {
	fp := "own/sys.uwl.txt"
	f, _ := os.Open(fp)
	data := ParseGeneral(f, GenRule{'\t', '\'', "wcf"})
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
	ioutil.WriteFile("own/sys_sogou.txt", Gen("sogou", data), 0777)
	ioutil.WriteFile("own/sys_baidu.txt", Gen("baidu", data), 0777)
	ioutil.WriteFile("own/sys_google.txt", Gen("google", data), 0777)
	ioutil.WriteFile("own/sys_qq.txt", Gen("qq", data), 0777)
	ioutil.WriteFile("own/sys_word_only.txt", Gen("word_only", data), 0777)
}
