package rose

import (
	"bytes"
	"fmt"

	util "github.com/flowerime/goutil"
)

var sg_pinyin = []string{"a", "ai", "an", "ang", "ao", "ba", "bai", "ban",
	"bang", "bao", "bei", "ben", "beng", "bi", "bian", "biao", "bie", "bin",
	"bing", "bo", "bu", "ca", "cai", "can", "cang", "cao", "ce", "cen", "ceng",
	"cha", "chai", "chan", "chang", "chao", "che", "chen", "cheng", "chi",
	"chong", "chou", "chu", "chua", "chuai", "chuan", "chuang", "chui", "chun",
	"chuo", "ci", "cong", "cou", "cu", "cuan", "cui", "cun", "cuo", "da",
	"dai", "dan", "dang", "dao", "de", "dei", "den", "deng", "di", "dia",
	"dian", "diao", "die", "ding", "diu", "dong", "dou", "du", "duan", "dui",
	"dun", "duo", "e", "ei", "en", "eng", "er", "fa", "fan", "fang", "fei",
	"fen", "feng", "fiao", "fo", "fou", "fu", "ga", "gai", "gan", "gang", "gao",
	"ge", "gei", "gen", "geng", "gong", "gou", "gu", "gua", "guai", "guan",
	"guang", "gui", "gun", "guo", "ha", "hai", "han", "hang", "hao", "he",
	"hei", "hen", "heng", "hong", "hou", "hu", "hua", "huai", "huan", "huang",
	"hui", "hun", "huo", "ji", "jia", "jian", "jiang", "jiao", "jie", "jin",
	"jing", "jiong", "jiu", "ju", "juan", "jue", "jun", "ka", "kai", "kan",
	"kang", "kao", "ke", "kei", "ken", "keng", "kong", "kou", "ku", "kua", "kuai",
	"kuan", "kuang", "kui", "kun", "kuo", "la", "lai", "lan", "lang", "lao",
	"le", "lei", "leng", "li", "lia", "lian", "liang", "liao", "lie", "lin",
	"ling", "liu", "lo", "long", "lou", "lu", "luan", "lve", "lun", "luo", "lv",
	"ma", "mai", "man", "mang", "mao", "me", "mei", "men", "meng", "mi", "mian",
	"miao", "mie", "min", "ming", "miu", "mo", "mou", "mu", "na", "nai", "nan",
	"nang", "nao", "ne", "nei", "nen", "neng", "ni", "nian", "niang", "niao",
	"nie", "nin", "ning", "niu", "nong", "nou", "nu", "nuan", "nve", "nun", "nuo",
	"nv", "o", "ou", "pa", "pai", "pan", "pang", "pao", "pei", "pen", "peng",
	"pi", "pian", "piao", "pie", "pin", "ping", "po", "pou", "pu", "qi", "qia",
	"qian", "qiang", "qiao", "qie", "qin", "qing", "qiong", "qiu", "qu", "quan",
	"que", "qun", "ran", "rang", "rao", "re", "ren", "reng", "ri", "rong", "rou",
	"ru", "rua", "ruan", "rui", "run", "ruo", "sa", "sai", "san", "sang", "sao",
	"se", "sen", "seng", "sha", "shai", "shan", "shang", "shao", "she", "shei",
	"shen", "sheng", "shi", "shou", "shu", "shua", "shuai", "shuan", "shuang",
	"shui", "shun", "shuo", "si", "song", "sou", "su", "suan", "sui", "sun", "suo",
	"ta", "tai", "tan", "tang", "tao", "te", "ten", "teng", "ti", "tian", "tiao",
	"tie", "ting", "tong", "tou", "tu", "tuan", "tui", "tun", "tuo", "wa", "wai",
	"wan", "wang", "wei", "wen", "weng", "wo", "wu", "xi", "xia", "xian", "xiang",
	"xiao", "xie", "xin", "xing", "xiong", "xiu", "xu", "xuan", "xue", "xun", "ya",
	"yan", "yang", "yao", "ye", "yi", "yin", "ying", "yo", "yong", "you", "yu",
	"yuan", "yue", "yun", "za", "zai", "zan", "zang", "zao", "ze", "zei", "zen",
	"zeng", "zha", "zhai", "zhan", "zhang", "zhao", "zhe", "zhei", "zhen", "zheng",
	"zhi", "zhong", "zhou", "zhu", "zhua", "zhuai", "zhuan", "zhuang", "zhui",
	"zhun", "zhuo", "zi", "zong", "zou", "zu", "zuan", "zui", "zun", "zuo",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p",
	"q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5",
	"6", "7", "8", "9", "#"}

type SogouBin struct{ Dict }

func NewSogouBin() *SogouBin {
	d := new(SogouBin)
	d.Name = "搜狗拼音备份.bin"
	d.Suffix = "bin"
	return d
}

func (d *SogouBin) Parse() {
	wl := make([]Entry, 0, d.size>>8)

	r := bytes.NewReader(d.data)
	header := make([]byte, 4) // SGPU
	r.Read(header)
	if !bytes.Equal(header, []byte{0x53, 0x47, 0x50, 0x55}) {
		r.Seek(0, 0)
		d.ParseOld()
		return
	}

	r.Seek(12, 1)
	fileSize := ReadUint32(r) // file total size
	r.Seek(36, 1)
	idxBegin := ReadUint32(r) // index begin
	idxSize := ReadUint32(r)  // index size
	wordCount := ReadUint32(r)
	dictBegin := ReadUint32(r)     // = idxBegin + idxSize
	dictTotalSize := ReadUint32(r) // file total size - dictBegin
	dictSize := ReadUint32(r)      // effective dict size
	fmt.Printf("fileSize: 0x%x\n", fileSize)
	fmt.Printf("idxBegin: 0x%x\n", idxBegin)
	fmt.Printf("idxSize: 0x%x\n", idxSize)
	fmt.Printf("dictBegin: 0x%x\n", dictBegin)
	fmt.Printf("dictTotalSize: 0x%x\n", dictTotalSize)
	fmt.Printf("dictSize: 0x%x\n", dictSize)

	for i := _u32; i < wordCount; i++ {
		r.Seek(int64(idxBegin+4*i), 0)
		idx := ReadUint32(r)
		if idx == 0 && i != 0 {
			break
		}
		r.Seek(int64(idx+dictBegin), 0)
		freq := ReadUint32(r)
		r.Seek(5, 1) // 00 00 FE 07 02
		pyLen := ReadUint16(r) / 2
		py := make([]string, 0, pyLen)
		for j := _u16; j < pyLen; j++ {
			p := ReadUint16(r)
			py = append(py, sg_pinyin[p])
		}
		ReadUint16(r) // word size + code size（include idx）
		wordSize := ReadUint16(r)
		tmp := make([]byte, wordSize)
		r.Read(tmp)
		word, _ := util.Decode(tmp, "UTF-16LE")

		wl = append(wl, &PinyinEntry{word, py, int(freq)})
		// repeat code
	}
	d.WordLibrary = wl
}
