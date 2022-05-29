package pinyin

import (
	"log"
	"os"
)

type PyEntry struct {
	Word  string
	Codes []string
	Freq  int
}

var (
	sogou     = GenRule{' ', '\'', "pw"}
	qq        = GenRule{' ', '\'', "cwf"}
	baidu     = GenRule{'\t', '\'', "wcf"}
	google    = GenRule{'\t', ' ', "wfc"}
	word_only = GenRule{' ', ' ', "w"}
)

func Parse(format, filepath string) []PyEntry {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		log.Panic("文件读取错误：", err)
	}
	switch format {
	case "baidu_bdict":
		return ParseBaiduBdict(f)
	case "baidu_bcd":
		return ParseBaiduBdict(f)
	case "sogou_scel":
		return ParseSogouScel(f)
	case "qq_qcel":
		return ParseSogouScel(f)
	case "ziguang_uwl":
		return ParseZiguangUwl(f)
	case "qq_qpyd":
		return ParseQqQpyd(f)
	case "mspy_dat":
		return ParseMspyDat(f)
	case "sogou_bin":
		return ParseSogouBin(f)
	case "sogou":
		return ParseGeneral(f, sogou)
	case "qq":
		return ParseGeneral(f, qq)
	case "baidu":
		return ParseGeneral(f, baidu)
	case "google":
		return ParseGeneral(f, google)
	case "word_only":
		return ParseGeneral(f, word_only)
	default:
		log.Panic("输入格式不支持：", format)
	}
	return []PyEntry{}
}

func Gen(format string, pe []PyEntry) []byte {
	switch format {
	case "sogou":
		return GenGeneral(pe, sogou)
	case "qq":
		return GenGeneral(pe, qq)
	case "baidu":
		return GenGeneral(pe, baidu)
	case "google":
		return GenGeneral(pe, google)
	case "word_only":
		return GenGeneral(pe, word_only)
	default:
		log.Panic("输出格式不支持：", format)
	}
	return []byte{}
}
