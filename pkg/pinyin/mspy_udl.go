package pinyin

import (
	"bytes"
	"fmt"
	"slices"
	"time"

	"github.com/nopdan/rose/pkg/util"
)

type MspyUDL struct {
	Template
	pyList []string
	pyMap  map[string]uint16
}

func init() {
	FormatList = append(FormatList, NewMspyUDL())
}
func NewMspyUDL() *MspyUDL {
	f := new(MspyUDL)
	f.Name = "微软拼音自学习词汇.dat"
	f.ID = "udl"
	f.CanMarshal = true

	f.pyList = []string{
		"a",
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
	f.pyMap = make(map[string]uint16)
	for idx, py := range f.pyList {
		f.pyMap[py] = uint16(idx)
	}
	return f
}

// 自学习词库，纯汉字
func (f *MspyUDL) Unmarshal(r *bytes.Reader) []*Entry {
	r.Seek(0xC, 0)
	count := ReadIntN(r, 4)
	di := make([]*Entry, 0, count)

	r.Seek(4, 1)
	export_stamp := ReadUint32(r)
	export_time := util.MsToTime(export_stamp)
	fmt.Printf("时间: %v\n", export_time)

	for i := 0; i < count; i++ {
		r.Seek(0x2400+60*int64(i), 0)
		data := ReadN(r, 60)
		// insert_stamp := BytesToInt(data[:4])
		// insert_time := MspyTime(uint32(insert_stamp))
		// jianpin := data[4:7]
		wordLen := int(data[10])
		p := 12 + wordLen*2
		if wordLen > 12 {
			fmt.Println(p, r.Size()-int64(r.Len()), data)
			continue
		}
		wordBytes := data[12:p]
		word := util.DecodeMust(wordBytes, "UTF-16LE")

		py := make([]string, 0, wordLen)
		for j := 0; j < wordLen; j++ {
			idx := BytesToInt(data[p+2*j : p+2*(j+1)])
			if idx < len(f.pyList) {
				py = append(py, f.pyList[idx])
			} else {
				fmt.Println("idx > len(mspy)", idx, len(f.pyList))
			}
		}
		di = append(di, &Entry{word, py, 1})
		// fmt.Printf("时间: %v, 简拼: %s, 词: %s, 拼音: %s\n", insert_time, string(jianpin), word, py)
	}
	return di
}

func (f *MspyUDL) Marshal(di []*Entry) []byte {
	var buf bytes.Buffer
	buf.Grow(0x2400 + 60*len(di))

	now := time.Now()
	timeBytes := util.MsTimeTo(now)
	buf.Write([]byte{0x55, 0xAA, 0x88, 0x81, 0x02, 0x00, 0x60, 0x00, 0x55, 0xAA, 0x55, 0xAA})
	count := len(di)
	buf.Write(make([]byte, 8))
	buf.Write(timeBytes)
	buf.Write(make([]byte, 0x2400-0x18))

	type udlEntry struct {
		Word    string
		Pinyin  []string
		Jianpin []byte
	}
	dict := make([]*udlEntry, 0, count)
	for i := range di {
		dict = append(dict, &udlEntry{Word: di[i].Word,
			Pinyin:  di[i].Pinyin,
			Jianpin: f.jianpin(di[i].Pinyin),
		})
	}
	// 按照简拼排序
	slices.SortStableFunc(dict, func(a, b *udlEntry) int {
		return bytes.Compare(a.Jianpin, b.Jianpin)
	})
	for _, v := range dict {
		b := make([]byte, 60)
		copy(b[:4], timeBytes)
		copy(b[4:7], v.Jianpin) // 3 bytes jianpin
		copy(b[7:10], []byte{0, 0, 4})
		wordBytes := util.EncodeMust(v.Word, "UTF-16LE")
		b[10] = byte(len(wordBytes) / 2)
		b[11] = 0x5A
		copy(b[12:], wordBytes)
		if len(wordBytes)/2 > 12 {
			fmt.Println("这个词太长了：", v.Word)
			count--
			continue
		} else {
			copy(b[12+len(wordBytes):], f.GetIndex(v.Pinyin))
		}
		buf.Write(b)
	}
	// 补 0
	size := 0x400 - len(buf.Bytes())%0x400
	buf.Write(make([]byte, size))
	b := buf.Bytes()
	copy(b[12:16], To4Bytes(count))
	b[0xE6C] = 1
	b[0xE98] = 2
	return b
}

func (f *MspyUDL) GetIndex(py []string) []byte {
	b := make([]byte, 0, len(py)/2)
	for _, v := range py {
		b = append(b, To2Bytes(f.pyMap[v])...)
	}
	return b
}

// 三位简拼
func (f *MspyUDL) jianpin(py []string) []byte {
	ret := make([]byte, 3)
	for i := 0; i < 3; i++ {
		if i >= len(py) {
			break
		}
		ret[i] = py[i][0]
	}
	return ret
}
