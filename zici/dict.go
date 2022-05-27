package zici

import (
	"fmt"
	"os"
	"sort"

	. "github.com/cxcn/dtool/utils"
)

type ZcEntry struct {
	Word string
	Code string
}

// 解析词库，指定码表格式
func Parse(format, filepath string) interface{} {
	f, err := os.Open(filepath)
	if err != nil {
		panic("文件读取失败：" + filepath)
	}
	switch format {
	case "duoduo":
		rd, _ := ReadFile(f)
		return ParseDuoduo(rd)
	case "bingling":
		rd, _ := ReadFile(f)
		return ParseBingling(rd)
	case "jidian":
		rd, _ := ReadFile(f)
		return ParseJidian(rd)
	case "baidu_def":
		return ParseBaiduDef(f)
	case "jidian_mb":
		return ParseJidianMb(f)
	}
	panic("解析失败：" + filepath + format)
}

// 生成指定格式词库
func Gen(format string, d interface{}) []byte {
	switch format {
	case "duoduo":
		return GenDuoduo(ToZcEntries(d))
	case "bingling":
		return GenBingling(ToZcEntries(d))
	case "jidian":
		return GenJidian(ToCodeEntries(d))
	case "baidu_def":
		return GenBaiduDef(ToCodeEntries(d))
	case "jidian_mb":
		fmt.Println("不支持该格式的生成")
	}
	panic("生成失败：" + format)
}

// 一码多词
type CodeEntry struct {
	Code  string
	Words []string
}

// 转为 多多形式，一码一词
func ToZcEntries(dict interface{}) []ZcEntry {
	var ce []CodeEntry
	switch dict.(type) {
	case []ZcEntry:
		return dict.([]ZcEntry)
	case []CodeEntry:
		ce = dict.([]CodeEntry)
	default:
		return []ZcEntry{}
	}
	ret := make([]ZcEntry, len(ce)*3/2)
	for _, v := range ce {
		for _, word := range v.Words {
			ret = append(ret, ZcEntry{word, v.Code})
		}
	}
	return ret
}

// 转为 极点形式，一码多词
func ToCodeEntries(dict interface{}) []CodeEntry {
	var zc []ZcEntry
	switch dict.(type) {
	case []CodeEntry:
		return dict.([]CodeEntry)
	case []ZcEntry:
		zc = dict.([]ZcEntry)
	default:
		return []CodeEntry{}
	}

	codeMap := make(map[string][]string)
	for _, v := range zc {
		if _, ok := codeMap[v.Code]; !ok {
			codeMap[v.Code] = []string{v.Word}
			continue
		}
		codeMap[v.Code] = append(codeMap[v.Code], v.Word)
	}

	ret := make([]CodeEntry, len(codeMap))
	for k, v := range codeMap {
		ret = append(ret, CodeEntry{k, v})
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Code < ret[j].Code
	})
	return ret
}

func SortByCode(dict []ZcEntry) []ZcEntry {
	sort.Slice(dict, func(i, j int) bool {
		return dict[i].Code < dict[j].Code
	})
	return dict
}
