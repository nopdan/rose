package rose

import (
	"bytes"
	"fmt"
	"time"

	util "github.com/flowerime/goutil"
)

type MspyUDL struct{ Dict }

func NewMspyUDL() *MspyUDL {
	d := new(MspyUDL)
	d.Name = "微软拼音自学习词汇.dat"
	d.Suffix = "dat"
	return d
}

// 自学习词库，纯汉字
func (d *MspyUDL) Parse() {
	r := bytes.NewReader(d.data)
	r.Seek(0xC, 0)
	count := ReadUint32(r)
	wl := make([]Entry, 0, count)

	r.Seek(4, 1)
	export_stamp := ReadUint32(r)
	export_time := MspyTime(export_stamp)
	fmt.Printf("时间: %v\n", export_time)

	for i := _u32; i < count; i++ {
		r.Seek(0x2400+60*int64(i), 0)
		data := make([]byte, 60)
		r.Read(data)

		// insert_stamp := util.BytesToInt(data[:4])
		// insert_time := MspyTime(uint32(insert_stamp))
		// jianpin := data[4:7]
		wordLen := data[10]
		p := 12 + int(wordLen)*2
		if p >= 60 {
			fmt.Println(p, r.Size()-int64(r.Len()), data)
		}
		wordSli := data[12:p]
		word, _ := Decode(wordSli, "UTF-16LE")

		py := make([]string, 0, wordLen)
		for j := 0; j < int(wordLen); j++ {
			idx := util.BytesToInt(data[p+2*j : p+2*(j+1)])
			if idx < len(mspy) {
				py = append(py, mspy[idx])
			} else {
				fmt.Println("idx > len(mspy)", idx, len(mspy))
			}
		}
		wl = append(wl, &PinyinEntry{word, py, 1})

		// fmt.Printf("时间: %v, 简拼: %s, 词: %s, 拼音: %s\n", insert_time, string(jianpin), word, py)
	}
	d.WordLibrary = wl
}

func (d *MspyUDL) GenFrom(wl WordLibrary) []byte {
	var buf bytes.Buffer
	buf.Grow(0x2400 + 60*len(wl))
	now := time.Now()
	timeBytes := MspyTimeTo(now)
	buf.Write([]byte{0x55, 0xAA, 0x88, 0x81, 0x02, 0x00, 0x60, 0x00, 0x55, 0xAA, 0x55, 0xAA})
	buf.Write(util.To4Bytes(uint32(len(wl))))
	buf.Write(make([]byte, 4))
	buf.Write(timeBytes)
	buf.Write(make([]byte, 0x2400-0x18))
	for i := range wl {
		b := make([]byte, 60)
		copy(b[:4], timeBytes)
		py := wl[i].GetPinyin()
		copy(b[4:7], d.jianpin(py)) // 3 bytes jianpin
		copy(b[7:10], []byte{0, 0, 4})
		w := wl[i].GetWord()
		word, _ := Encode(w, "UTF-16LE")
		b[10] = byte(len(word) / 2)
		b[11] = 0x5A
		copy(b[12:], word)
		copy(b[12+len(word):], MspyGetIndex(py))
		buf.Write(b)
	}
	// 补 0
	size := 1024 - (0x2400+60*len(wl))%1024
	buf.Write(make([]byte, size))
	return buf.Bytes()
}

func MspyGetIndex(py []string) []byte {
	ret := make([]byte, 0, len(py)/2)
	for _, v := range py {
		ret = append(ret, util.To2Bytes(mspyMap[v])...)
	}
	return ret
}

func MspyTime(stamp uint32) time.Time {
	return time.Unix(int64(stamp), 0).Add(946684800 * time.Second)
}

func MspyTimeTo(t time.Time) []byte {
	return util.To4Bytes(uint32(t.Add(-946684800 * time.Second).Unix()))
}

func (d *MspyUDL) jianpin(py []string) []byte {
	ret := make([]byte, 3)
	for i := 0; i < 3; i++ {
		if i >= len(py) {
			break
		}
		ret[i] = py[i][0]
	}
	return ret
}

var mspyMap = make(map[string]uint16)

func init() {
	for i, v := range mspy {
		mspyMap[v] = uint16(i)
	}
}

var mspy = []string{"a",
	"ai",
	"an",
	"ang",
	"ao",
	"ba",
	"bai",
	"ban",
	"bang",
	"bao",
	"bei",
	"ben",
	"beng",
	"bi",
	"bian",
	"biao",
	"bie",
	"bin",
	"bing",
	"bo",
	"bu",
	"ca",
	"cai",
	"can",
	"cang",
	"cao",
	"ce",
	"cen",
	"ceng",
	"cha",
	"chai",
	"chan",
	"chang",
	"chao",
	"che",
	"chen",
	"cheng",
	"chi",
	"chong",
	"chou",
	"chu",
	"chua",
	"chuai",
	"chuan",
	"chuang",
	"chui",
	"chun",
	"chuo",
	"ci",
	"cong",
	"cou",
	"cu",
	"cuan",
	"cui",
	"cun",
	"cuo",
	"da",
	"dai",
	"dan",
	"dang",
	"dao",
	"de",
	"dei",
	"den",
	"deng",
	"di",
	"dia",
	"dian",
	"diao",
	"die",
	"ding",
	"diu",
	"dong",
	"dou",
	"du",
	"duan",
	"dui",
	"dun",
	"duo",
	"e",
	"ei",
	"en",
	"eng",
	"er",
	"fa",
	"fan",
	"fang",
	"fei",
	"fen",
	"feng",
	"fiao",
	"fo",
	"fou",
	"fu",
	"ga",
	"gai",
	"gan",
	"gang",
	"gao",
	"ge",
	"gei",
	"gen",
	"geng",
	"gong",
	"gou",
	"gu",
	"gua",
	"guai",
	"guan",
	"guang",
	"gui",
	"gun",
	"guo",
	"ha",
	"hai",
	"han",
	"hang",
	"hao",
	"he",
	"hei",
	"hen",
	"heng",
	"hong",
	"hou",
	"hu",
	"hua",
	"huai",
	"huan",
	"huang",
	"hui",
	"hun",
	"huo",
	"ji",
	"jia",
	"jian",
	"jiang",
	"jiao",
	"jie",
	"jin",
	"jing",
	"jiong",
	"jiu",
	"ju",
	"juan",
	"jue",
	"jun",
	"ka",
	"kai",
	"kan",
	"kang",
	"kao",
	"ke",
	"kei",
	"ken",
	"keng",
	"kong",
	"kou",
	"ku",
	"kua",
	"kuai",
	"kuan",
	"kuang",
	"kui",
	"kun",
	"kuo",
	"la",
	"lai",
	"lan",
	"lang",
	"lao",
	"le",
	"lei",
	"leng",
	"li",
	"lia",
	"lian",
	"liang",
	"liao",
	"lie",
	"lin",
	"ling",
	"liu",
	"lo",
	"long",
	"lou",
	"lu",
	"luan",
	"lve", // lue
	"lun",
	"luo",
	"lv",
	"ma",
	"mai",
	"man",
	"mang",
	"mao",
	"me",
	"mei",
	"men",
	"meng",
	"mi",
	"mian",
	"miao",
	"mie",
	"min",
	"ming",
	"miu",
	"mo",
	"mou",
	"mu",
	"na",
	"nai",
	"nan",
	"nang",
	"nao",
	"ne",
	"nei",
	"nen",
	"neng",
	"ni",
	"nian",
	"niang",
	"niao",
	"nie",
	"nin",
	"ning",
	"niu",
	"nong",
	"nou",
	"nu",
	"nuan",
	"nve", // nue
	"nun",
	"nuo",
	"nv",
	"o",
	"ou",
	"pa",
	"pai",
	"pan",
	"pang",
	"pao",
	"pei",
	"pen",
	"peng",
	"pi",
	"pian",
	"piao",
	"pie",
	"pin",
	"ping",
	"po",
	"pou",
	"pu",
	"qi",
	"qia",
	"qian",
	"qiang",
	"qiao",
	"qie",
	"qin",
	"qing",
	"qiong",
	"qiu",
	"qu",
	"quan",
	"que",
	"qun",
	"ran",
	"rang",
	"rao",
	"re",
	"ren",
	"reng",
	"ri",
	"rong",
	"rou",
	"ru",
	"rua",
	"ruan",
	"rui",
	"run",
	"ruo",
	"sa",
	"sai",
	"san",
	"sang",
	"sao",
	"se",
	"sen",
	"seng",
	"sha",
	"shai",
	"shan",
	"shang",
	"shao",
	"she",
	"shei",
	"shen",
	"sheng",
	"shi",
	"shou",
	"shu",
	"shua",
	"shuai",
	"shuan",
	"shuang",
	"shui",
	"shun",
	"shuo",
	"si",
	"song",
	"sou",
	"su",
	"suan",
	"sui",
	"sun",
	"suo",
	"ta",
	"tai",
	"tan",
	"tang",
	"tao",
	"te",
	"tei",
	"teng",
	"ti",
	"tian",
	"tiao",
	"tie",
	"ting",
	"tong",
	"tou",
	"tu",
	"tuan",
	"tui",
	"tun",
	"tuo",
	"wa",
	"wai",
	"wan",
	"wang",
	"wei",
	"wen",
	"weng",
	"wo",
	"wu",
	"xi",
	"xia",
	"xian",
	"xiang",
	"xiao",
	"xie",
	"xin",
	"xing",
	"xiong",
	"xiu",
	"xu",
	"xuan",
	"xue",
	"xun",
	"ya",
	"yan",
	"yang",
	"yao",
	"ye",
	"yi",
	"yin",
	"ying",
	"yo",
	"yong",
	"you",
	"yu",
	"yuan",
	"yue",
	"yun",
	"za",
	"zai",
	"zan",
	"zang",
	"zao",
	"ze",
	"zei",
	"zen",
	"zeng",
	"zha",
	"zhai",
	"zhan",
	"zhang",
	"zhao",
	"zhe",
	"zhei",
	"zhen",
	"zheng",
	"zhi",
	"zhong",
	"zhou",
	"zhu",
	"zhua",
	"zhuai",
	"zhuan",
	"zhuang",
	"zhui",
	"zhun",
	"zhuo",
	"zi",
	"zong",
	"zou",
	"zu",
	"zuan",
	"zui",
	"zun",
	"zuo",
}
