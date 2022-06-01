package dtool

import (
	"log"

	. "github.com/cxcn/dtool/pinyin"
	. "github.com/cxcn/dtool/utils"
	. "github.com/cxcn/dtool/zici"
)

// 拼音
func PinyinParse(format, filename string) WpfDict {
	switch format {
	case "baidu_bdict":
		return ParseBaiduBdict(filename)
	case "baidu_bcd":
		return ParseBaiduBdict(filename)
	case "sogou_scel":
		return ParseSogouScel(filename)
	case "qq_qcel":
		return ParseSogouScel(filename)
	case "ziguang_uwl":
		return ParseZiguangUwl(filename)
	case "qq_qpyd":
		return ParseQqQpyd(filename)
	case "mspy_dat":
		return ParseMspyDat(filename)
	case "mspy_udl":
		return ParseMspyUDL(filename)
	case "sogou_bin":
		return ParseSogouBin(filename)
	// 纯文本拼音
	case "pyjj":
		return ParseJiaJia(filename)
	case "word_only":
		return ParseWordOnly(filename)
	case "sogou":
		return ParsePinyin(filename, R_sogou)
	case "qq":
		return ParsePinyin(filename, R_qq)
	case "baidu":
		return ParsePinyin(filename, R_baidu)
	case "google":
		return ParsePinyin(filename, R_google)
	default:
		log.Panic("输入格式不支持：", format)
	}
	return WpfDict{}
}

func PinyinGen(format string, wpfd WpfDict) []byte {
	switch format {
	case "pyjj":
		return GenJiaJia(wpfd)
	case "word_only":
		return GenWordOnly(wpfd)
	case "sogou":
		return GenPinyin(wpfd, R_sogou)
	case "qq":
		return GenPinyin(wpfd, R_qq)
	case "baidu":
		return GenPinyin(wpfd, R_baidu)
	case "google":
		return GenPinyin(wpfd, R_google)
	default:
		log.Panic("输出格式不支持：", format)
	}
	return []byte{}
}

// 字词
func ZiciParse(format, filename string) WcTable {
	switch format {
	// 字词的
	case "baidu_def":
		return ParseBaiduDef(filename)
	case "jidian_mb":
		return ParseJidianMb(filename)
	case "fcitx4_mb":
		return ParseFcitx4Mb(filename)
	// 字词的纯文本
	case "duoduo":
		return ParseDuoduo(filename)
	case "bingling":
		return ParseBingling(filename)
	case "jidian":
		return ParseJidian(filename)
	default:
		log.Panic("输入格式不支持：", format)
	}
	return WcTable{}
}

func ZiciGen(format string, wct WcTable) []byte {
	switch format {
	case "duoduo":
		return GenDuoduo(wct)
	case "bingling":
		return GenBingling(wct)
	case "jidian":
		return GenJidian(wct)
	case "baidu_def":
		return GenBaiduDef(wct)
	default:
		log.Panic("输出格式不支持：", format)
	}
	return []byte{}
}
