package wubi

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
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
	f = NewCustom("t|w|c|r", "utf-8")
	data := f.Marshal(di, false)
	// get filename without ext
	s := "test/" + strings.TrimSuffix(filepath.Base(name), filepath.Ext(name)) + ".txt"
	os.WriteFile(s, data, 0666)
}

func TestMsudp(t *testing.T) {
	f := NewMsUDP()
	for _, path := range []string{
		"./sample/UserDefinedPhrase.dat",
		"./sample/ChsPinyinUDP.lex",
	} {
		do(f, path)
	}
}

func TestMswbLex(t *testing.T) {
	path := "./sample/ChsWubi.lex"
	f := NewMswbLex()
	do(f, path)
}

func TestBaiduDef(t *testing.T) {
	// 哲哲豆词库
	path := "./sample/baidu.def"
	f := NewBaiduDef()
	do(f, path)
}

func TestJidianMb(t *testing.T) {
	// 091 点儿词库
	path := "./sample/jidian.mb"
	f := NewJidianMb()
	do(f, path)
}

func TestFcitx4Mb(t *testing.T) {
	path := "./sample/98wb_ci.mb"
	f := NewFcitx4Mb()
	do(f, path)
}

func TestDuoDB(t *testing.T) {
	path := "./sample/main.duodb"
	f := NewDuoDB()
	do(f, path)
}

func TestMarshal(t *testing.T) {
	// 哲哲豆词库 1w 多条
	path := "./sample/duoduo.txt"
	f := NewDuoduo()
	data, _ := os.ReadFile(path)
	r := bytes.NewReader(data)
	di := f.Unmarshal(r)
	di = slices.DeleteFunc(di, func(e *Entry) bool {
		return e.Code == ""
	})

	os.WriteFile("test/to_bingling.txt", New("bl").Marshal(di, false), 0666)
	os.WriteFile("test/to_baidu_def.def", New("def").Marshal(di, false), 0666)
	os.WriteFile("test/to_msudp.dat", New("udp").Marshal(di, false), 0666)
	os.WriteFile("test/to_lex.lex", New("lex").Marshal(di, false), 0666)
}
