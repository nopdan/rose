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

	format = "mspy_dat"
	fp = "own/SuperRime拓展词库 for Win10拼音版(600万词-含BetterRime)-v20.3.dat"
	data = Parse(format, fp)
	write(fp, data)

}

func write(filename string, data []PyEntry) {
	ioutil.WriteFile(fmt.Sprintf("%s.txt", filename), GenGenersal(data), 0777)
}
