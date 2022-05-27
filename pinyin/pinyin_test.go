package pinyin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func Test(t *testing.T) {

	// 计算机硬件词汇 https://shurufa.baidu.com/dict
	format := "baidu_bdict"
	data := Parse(format, "test/baidu.bdict")
	write(format, data)

	// 计算机 https://mime.baidu.com/web/iw/index/
	format = "baidu_bcd"
	data = Parse(format, "test/baidu.bcd")
	write(format, data)

	// 搜狗细胞词库
	format = "sogou_scel"
	data = Parse(format, "test/sogou.scel")
	write(format, data)

	// 网络流行新词【官方推荐】 http://cdict.qq.pinyin.cn/detail?dict_id=s4
	format = "qq_qcel"
	data = Parse(format, "test/qq.qcel")
	write(format, data)

	// 来源，紫光内置
	format = "ziguang_uwl"
	data = Parse(format, "test/music.uwl")
	write(format, data)

	format = "ziguang_uwl"
	data = Parse(format, "none/sys.uwl")
	write("ziguang_uwl_sys", data)

	// QQ
	format = "qq_qpyd"
	data = Parse(format, "test/qq.qpyd")
	write(format, data)
}

func write(filename string, data []PyEntry) {
	var buf bytes.Buffer
	for _, v := range data {
		buf.WriteString(v.Word + "\t" + strings.Join(v.Codes, "'") + "\t" + strconv.Itoa(v.Freq) + "\r\n")
	}
	ioutil.WriteFile(fmt.Sprintf("out/%s.txt", filename), buf.Bytes(), 0777)
}
