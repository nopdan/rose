package pinyin

import (
	"bytes"
	"fmt"
	"io"

	"github.com/nopdan/rose/pkg/util"
)

type Kafan struct {
	Template
	pyList []string
}

func NewKafan() *Kafan {
	f := new(Kafan)
	f.Name = "卡饭拼音备份.bin"
	f.ID = "kfpybin"

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

func (k *Kafan) Unmarshal(r *bytes.Reader) []*Entry {
	wl := make([]*Entry, 0, 0xff)
	for r.Len() > 8 {
		check := make([]byte, 8)
		r.Read(check)
		if string(check) == "ProtoDic" {
			r.Seek(8, io.SeekCurrent)
			break
		}
	}

	// 跳过 kf_pinyin 之前的内容
	for {
		tmp := make([]byte, 2)
		r.Read(tmp)
		if string(tmp) == "kf" {
			r.Seek(0x10-2, io.SeekCurrent)
			break
		}
	}

	for r.Len() > 0x28 {
		pinyin := make([]string, 0, 2)
		for {
			tmp := make([]byte, 8)
			_, err := r.Read(tmp)
			if err != nil {
				break
			}
			util.PrintHex(tmp)
			idx := BytesToInt(tmp[:4])

			if idx == 0 {
				idx2 := BytesToInt(tmp[4:])
				if idx2 != 0 && idx2 < 512 {
					r.Seek(-4, io.SeekCurrent)
				}
				continue
			}

			if bytes.Equal(tmp, []byte{0x04, 0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x00}) {
				fmt.Println("break")
				r.Seek(0x30-12, io.SeekCurrent)
				break
			}

			if BytesToInt(tmp[4:]) != 0 {
				continue
			}

			if idx < 512 {
				pinyin = append(pinyin, k.Lookup(idx))
				fmt.Println(idx, k.Lookup(idx))
			}

			// tmp := make([]byte, 0x28)
			// if idx > 512 || idx == 0 {
			// 	r.Seek(4-0x28, io.SeekCurrent)
			// 	continue
			// }

			// if tmp[32] != 0 {
			// 	r.Seek(8, io.SeekCurrent)
			// }

			// if !bytes.Equal(tmp[4:8], []byte{0x00, 0x00, 0x00, 0x00}) {
			// 	fmt.Println("continue")
			// 	continue
			// }
		}
		freq := ReadUint32(r)
		wordBytes := make([]byte, 0, 4)
		for {
			tmp := make([]byte, 4)
			r.Read(tmp)
			wordBytes = append(wordBytes, tmp...)
			if tmp[3] == 0 {
				break
			}
		}
		wl = append(wl, &Entry{
			Word:   string(wordBytes),
			Pinyin: pinyin,
			Freq:   int(freq),
		})
		fmt.Println(Entry{Pinyin: pinyin, Freq: int(freq), Word: string(wordBytes)})
	}

	return wl
}

func (k *Kafan) NewUnmarshal(r *bytes.Reader) []Entry {
	wl := make([]Entry, 0, 0xff)
	for r.Len() > 8 {
		check := make([]byte, 8)
		r.Read(check)
		if string(check) == "ProtoDic" {
			r.Seek(8, io.SeekCurrent)
			break
		}
	}
	return wl
}

func (k *Kafan) Test(r *bytes.Reader) {
	for r.Len() > 8 {
		check := make([]byte, 8)
		r.Read(check)
		if string(check) == "ProtoDic" {
			r.Seek(8, io.SeekCurrent)
			break
		}
	}
	for r.Len() > 8 {
		k.readPinyin(r)
	}
}

func (k *Kafan) readPinyin(r *bytes.Reader) []string {
	pinyin := make([]string, 0, 2)
	for {
		tmp := make([]byte, 8)
		_, err := r.Read(tmp)
		if err != nil {
			break
		}
		// PrintHex(tmp)
		idx := BytesToInt(tmp[:4])

		if idx == 0 {
			idx2 := BytesToInt(tmp[4:])
			if idx2 != 0 && idx2 < 512 {
				r.Seek(-4, io.SeekCurrent)
			}
			continue
		}

		if bytes.Equal(tmp, []byte{0x04, 0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x00}) {
			// fmt.Println("break")
			r.Seek(-8, io.SeekCurrent)
			break
		}

		if BytesToInt(tmp[4:]) != 0 {
			continue
		}

		if idx < 512 {
			pinyin = append(pinyin, k.Lookup(idx))
		}
	}

	for {
		tmp := make([]byte, 8)
		r.Read(tmp)
		if BytesToInt(tmp[4:]) == 0 {
			r.Seek(-8, io.SeekCurrent)
			break
		}
		if BytesToInt(tmp[:4]) == 0 {
			r.Seek(-4, io.SeekCurrent)
			continue
		}

		if bytes.Equal(tmp, []byte{0x04, 0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x00}) {
			r.Seek(0x30-8, io.SeekCurrent)
			wordBytes := make([]byte, 0, 4)
			for {
				tmp := make([]byte, 4)
				r.Read(tmp)
				wordBytes = append(wordBytes, tmp...)
				if tmp[3] == 0 {
					break
				}
			}
			fmt.Println(string(wordBytes))
		}

	}

	fmt.Println(pinyin)
	fmt.Println()
	return pinyin
}

func (k *Kafan) Lookup(idx int) string {
	if len(k.pyList) <= idx {
		fmt.Println("index out of range: ", idx, ">", len(k.pyList)-1)
		return ""
	}
	return k.pyList[idx]
}
