package pinyin

import "os"

type Pinyin struct {
	Word string
	Code []string
	Freq int
}

func Parse(format, filepath string) []Pinyin {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	switch format {
	case "baidu_bdict":
		return ParseBaiduBdict(f)
	case "baidu_bcd":
		return ParseBaiduBdict(f)
	case "sougou_scel":
		return ParseSougouScel(f)
	case "ziguang_uwl":
		return ParseZiguangUwl(f)
	case "qq_qpyd":
		return ParseQqQpyd(f)
	}
	return []Pinyin{}
}

// 字节（小端）转为整数
func bytesToInt(b []byte) int {
	ret := int(b[1])*0x100 + int(b[0])
	if len(b) >= 3 {
		ret += int(b[2]) * 0x10000
	}
	if len(b) >= 4 {
		ret += int(b[3]) * 0x1000000
	}
	return ret
}
