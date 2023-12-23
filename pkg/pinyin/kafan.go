package pinyin

import (
	"bytes"
	"fmt"
	"io"
	"strings"

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
	f.ID = "kfpybak,dict"

	f.pyList = []string{
		" ", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p",
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

	// 0x48 or 0x68
	r.Seek(0x48, io.SeekStart)
	head := string(ReadN(r, 0x10))
	// 版本不匹配
	if !strings.HasPrefix(head, "ProtoDict1") {
		// 有的词库是在 0x68
		r.Seek(0x68, io.SeekStart)
		head = string(ReadN(r, 0x10))
		if !strings.HasPrefix(head, "ProtoDict1") {
			fmt.Println("卡饭拼音备份.dict格式错误")
			return nil
		}
	}

	di := make([]*Entry, 0, 0xff)
	// 读取一个词
	for r.Len() > 0x28 {
		// 词库中间可能夹杂这个
		dictType := ReadN(r, 0x10)
		if !bytes.HasPrefix(dictType, []byte("kf_pinyin")) {
			r.Seek(-0x10, io.SeekCurrent)
		}

		// 读取编码占用的字节
		codeBytes := make([]byte, 0, 0x28)
		for {
			// 每次读取 8 个字节
			tmp := ReadN[int](r, 8)
			// 判断结束
			if bytes.Equal(tmp, []byte{4, 0, 0, 0, 3, 0, 1, 0}) {
				r.Seek(0x20, io.SeekCurrent)
				break
			} else if bytes.Equal(tmp, []byte{0, 0, 0, 0, 3, 0, 1, 0}) {
				r.Seek(0x18, io.SeekCurrent)
				break
			}
			codeBytes = append(codeBytes, tmp...)
		}

		// 转换为拼音
		pinyin := make([]string, 0, 2)
		// 每 0x28 个字节为一个音
		for i := len(codeBytes) % 0x28; i < len(codeBytes); i += 0x28 {
			idx := BytesToInt(codeBytes[i : i+4])
			py := f.lookup(idx, r)
			if py == "" {
				fmt.Printf("codeBytes: %v\n", codeBytes)
			} else if py != " " {
				pinyin = append(pinyin, py)
			}
		}

		// 跳过未知的4字节
		mark := ReadIntN(r, 4)
		if mark != 1 {
			r.Seek(8, io.SeekCurrent)
		}
		size := ReadIntN(r, 4)
		// 22	3	8
		// 2A	4
		// 32	5
		// 3A	6	8
		// 42	7	8
		// 4A	8	8
		// 52	9	16
		// 6A	12	16
		if size%0x10 == 2 {
			size = (size/0x10)*2 - 1
		} else if size%0x10 == 0xA {
			size = (size / 0x10) * 2
		} else {
			fmt.Printf("读取词组错误, size: 0x%x, offset: 0x%x\n", size, int(r.Size())-r.Len())
			return nil
		}

		// 下面读取词，词是按照8字节对齐的
		wordBytes := ReadN(r, size)
		if len(wordBytes)%8 != 0 {
			r.Seek(int64(8-len(wordBytes)%8), io.SeekCurrent)
		}
		word := string(wordBytes)
		// di = append(di, &Entry{
		// 	Word:   word,
		// 	Pinyin: pinyin,
		// 	Freq:   1,
		// })
		if py := f.filter(word, pinyin); len(py) > 0 {
			di = append(di, &Entry{
				Word:   word,
				Pinyin: py,
				Freq:   1,
			})
			fmt.Printf("词组: %s, 拼音: %v\n", word, py)
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

func (k *Kafan) lookup(idx int, r *bytes.Reader) string {
	if len(k.pyList) <= idx {
		fmt.Printf("index out of range: %d > %d, offset: 0x%x\n", idx, len(k.pyList)-1, int(r.Size())-r.Len())
		return ""
	}
	return k.pyList[idx]
}
