package zici

import (
	"log"
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
	defer f.Close()
	if err != nil {
		log.Panic("文件读取失败：" + filepath)
	}

	switch format {
	case "duoduo":
		rd, _ := DecodeIO(f)
		return ParseDuoduo(rd)
	case "bingling":
		rd, _ := DecodeIO(f)
		return ParseBingling(rd)
	case "jidian":
		rd, _ := DecodeIO(f)
		return ParseJidian(rd)
	case "baidu_def":
		return ParseBaiduDef(f)
	case "jidian_mb":
		return ParseJidianMb(f)
	case "fcitx4_mb":
		return ParseFcitx4Mb(f)
	default:
		log.Panic("输入格式不支持：", format)
	}
	return []ZcEntry{}
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
	default:
		log.Panic("输出格式不支持：", format)
	}
	return []byte{}
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
		log.Panic("内部码表格式错误：", dict)
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
		log.Panic("内部码表格式错误：", dict)
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
