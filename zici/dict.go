package zici

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
)

// 一码多词
type Dict map[string][]string

func (d Dict) insert(code, word string) {
	if _, ok := d[code]; !ok {
		d[code] = []string{word}
	} else {
		d[code] = append(d[code], word)
	}
}

type codeAndWords struct {
	code  string
	words []string
}

// map 转为切片，方便按编码排序
func (d Dict) toSlice() []codeAndWords {
	ret := make([]codeAndWords, len(d))
	for k, v := range d {
		ret = append(ret, codeAndWords{k, v})
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].code < ret[j].code
	})
	return ret
}

// 解析词库，指定码表格式
func Parse(format, filepath string) Dict {
	f, err := os.Open(filepath)
	if err != nil {
		panic("文件读取失败：" + filepath)
	}
	switch format {
	case "duoduo":
		rd, _ := conv(f)
		return ParseDuoduo(rd)
	case "bingling":
		rd, _ := conv(f)
		return ParseBingling(rd)
	case "jidian":
		rd, _ := conv(f)
		return ParseJidian(rd)
	case "baidu_def":
		return ParseBaiduDef(f)
	case "jidian_mb":
		return ParseJidianMb(f)
	}
	return Dict{}
}

// 生成指定格式词库
func Gen(format string, d Dict) []byte {
	dl := d.toSlice()
	switch format {
	case "duoduo":
		return GenDuoduo(dl)
	case "bingling":
		return GenBingling(dl)
	case "jidian":
		return GenJidian(dl)
	case "baidu_def":
		return GenBaiduDef(dl)
	case "jidian_mb":
		fmt.Println("不支持该格式的生成")
	}
	return []byte{}
}

// 将 io流 转换为 utf-8
func conv(f io.Reader) (io.Reader, error) {

	brd := bufio.NewReader(f)
	buf, _ := brd.Peek(1024)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf) // 检测编码格式
	if err != nil {
		return brd, err
	}
	if cs.Confidence != 100 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}
	rd, err := charset.NewReaderLabel(cs.Charset, brd) // 转换字节流
	return rd, err
}
