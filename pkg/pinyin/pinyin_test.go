package pinyin

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func do(f Format, name string) {
	os.Mkdir("test", 0666)
	b, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(b)
	di := f.Unmarshal(r)
	f = NewRime()
	data := f.Marshal(di)
	// get filename without ext
	s := "test/" + strings.TrimSuffix(filepath.Base(name), filepath.Ext(name)) + ".txt"
	os.WriteFile(s, data, 0666)
}

func TestScel(t *testing.T) {
	f := NewSogouScel()
	for _, path := range []string{
		"sample/sogou.scel",
		"sample/qq.qcel", // 网络流行新词【官方推荐】 http://cdict.qq.pinyin.cn/detail?dict_id=s4

		// "own/搜狗标准词库.scel",
	} {
		do(f, path)
	}
}

func TestSogouBin(t *testing.T) {
	f := NewSogouBak()
	for _, path := range []string{
		"sample/sogou_bak_v2.bin", // 搜狗词库备份_2021_05_10.bin

		"sample/sogou_bak_v3.bin", // new
	} {
		do(f, path)
	}
}

func TestQpyd(t *testing.T) {
	f := NewQqQpyd()
	path := "sample/qq.qpyd"
	do(f, path)
}

func TestBdict(t *testing.T) {
	f := NewBaiduBdict()
	for _, path := range []string{
		"sample/baidu.bdict", // 计算机硬件词汇 https://shurufa.baidu.com/dict
		"sample/baidu.bcd",   // 计算机 https://mime.baidu.com/web/iw/index/
	} {
		do(f, path)
	}
}

func TestUwl(t *testing.T) {
	f := NewZiguangUwl()
	for _, path := range []string{
		"sample/music.uwl", // 来源，紫光内置

		// "own/sys7.uwl",
		// "own/大词库第六版.uwl",
	} {
		do(f, path)
	}
}

func TestMspyUDL(t *testing.T) {
	f := NewMspyUDL()
	path := "sample/ChsPinyinUDL.dat"
	do(f, path)
}

// 拼音加加测试
func TestJiajia(t *testing.T) {
	f := NewJiaJia()
	path := "sample/jiajia.txt"
	do(f, path)
}

func TestMarshal(t *testing.T) {
	f := NewSogouScel()
	r, _ := os.ReadFile("sample/qq.qcel")
	di := f.Unmarshal(bytes.NewReader(r))
	os.WriteFile("test/to_msudl.dat", New("udl").Marshal(di), 0666)
	os.WriteFile("test/to_pyjj.txt", New("jiajia").Marshal(di), 0666)
	os.WriteFile("test/to_sogou.txt", New("sg").Marshal(di), 0666)
	os.WriteFile("test/to_qq.txt", New("qq").Marshal(di), 0666)
	os.WriteFile("test/to_baidu.txt", New("bd").Marshal(di), 0666)
	os.WriteFile("test/to_google.txt", New("gg").Marshal(di), 0666)
	os.WriteFile("test/to_rime.txt", New("rime").Marshal(di), 0666)
}

func TestFormatList(t *testing.T) {
	for _, v := range FormatList {
		if !v.GetCanMarshal() {
			fmt.Printf("不")
		}
		fmt.Printf("可导出")
		fmt.Printf(" %s \t %s\n", v.GetID(), v.GetName())
	}
}
