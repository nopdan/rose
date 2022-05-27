package pinyin

import (
	"os"
)

type PyEntry struct {
	Word  string
	Codes []string
	Freq  int
}

func Parse(format, filepath string) []PyEntry {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
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
	}
	return []PyEntry{}
}
