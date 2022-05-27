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
func Parse(format, filepath string) []ZcEntry {
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
	return []ZcEntry{}
}

// 生成指定格式词库
func Gen(format string, d []CodeEntry) []byte {
	switch format {
	case "duoduo":
		return GenDuoduo(d)
	case "bingling":
		return GenBingling(d)
	case "jidian":
		return GenJidian(d)
	case "baidu_def":
		return GenBaiduDef(d)
	case "jidian_mb":
		fmt.Println("不支持该格式的生成")
	}
	return []byte{}
}

// 一码多词
type CodeEntry struct {
	Code  string
	Words []string
}

// 转为 极点形式的码表，一码多词
func ToCodeEntries(dict []ZcEntry) []CodeEntry {
	codeMap := make(map[string][]string)
	for _, v := range dict {
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
