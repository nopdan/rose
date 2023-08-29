package pinyin

import (
	"bytes"
	"fmt"
	"io"

	"github.com/nopdan/rose/pkg/encoder"
)

type Kafan struct {
	Template
	pyList []string
}

func init() {
	FormatList = append(FormatList, NewKafan())
}

func NewKafan() *Kafan {
	f := new(Kafan)
	f.Name = "卡饭拼音备份.dict"
	f.ID = "kfpybak"

	f.pyList = []string{
		"", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p",
		"q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
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
		"lun",
		"luo",
		"lv",
		"lve", // lue
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
		"nia",
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
		"nun",
		"nuo",
		"nv",
		"nve",
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
	return f
}

func (f *Kafan) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, 0xff)
	for r.Len() > 8 {
		check := make([]byte, 8)
		r.Read(check)
		if string(check) == "ProtoDic" {
			r.Seek(8, io.SeekCurrent)
			break
		}
	}
	for r.Len() > 0x28 {
		tmp := ReadN(r, 4)
		// kf_pinyin
		if bytes.Equal(tmp, []byte{0x6B, 0x66, 0x5F, 0x70}) {
			r.Seek(12, io.SeekCurrent)
			continue
		}
		if BytesToInt(tmp) == 0 {
			continue
		}
		r.Seek(-4, io.SeekCurrent)
		pinyin := make([]string, 0, 2)
		var word string
		for {
			// 每40个字节为一个音
			tmp := ReadN[int](r, 0x28) // 40
			if bytes.Equal(tmp[:8], []byte{4, 0, 0, 0, 3, 0, 1, 0}) {
				r.Seek(8, io.SeekCurrent) // 未知
				wordBytes := make([]byte, 0, 4)
				for {
					b := ReadN[int](r, 4)
					wordBytes = append(wordBytes, b...)
					if b[3] == 0 {
						break
					}
				}
				// 去除末尾的 0
				for i := len(wordBytes) - 1; i >= 0 && wordBytes[i] == 0; i-- {
					wordBytes = wordBytes[:i]
				}
				word = string(wordBytes)
				break
			}
			idx := BytesToInt(tmp[:4])
			pinyin = append(pinyin, f.lookup(idx))
		}
		if py := f.filter(word, pinyin); len(py) > 0 {
			di = append(di, &Entry{
				Word:   word,
				Pinyin: py,
				Freq:   1,
			})
		}
	}
	return di
}

func (k *Kafan) filter(word string, pinyin []string) []string {
	wordRunes := []rune(word)
	// 过滤单字
	if len(wordRunes) <= 1 {
		return nil
	}
	if len(wordRunes) == len(pinyin) {
		return pinyin
	}
	if len(wordRunes) < len(pinyin) {
		return pinyin[len(pinyin)-len(wordRunes):]
	}
	if len(wordRunes) > len(pinyin) {
		enc := encoder.NewPinyin()
		pre := string(wordRunes[:len(wordRunes)-len(pinyin)])
		res := append(enc.Encode(pre), pinyin...)
		return res
	}
	return nil
}

func (k *Kafan) lookup(idx int) string {
	if len(k.pyList) <= idx {
		fmt.Println("index out of range: ", idx, ">", len(k.pyList)-1)
		return ""
	}
	return k.pyList[idx]
}
