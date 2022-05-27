package pinyin

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/cxcn/dtool/utils"
)

var uwlSm = []string{
	"", "b", "c", "ch", "d", "f", "g", "h", "j", "k", "l", "m", "n",
	"p", "q", "r", "s", "sh", "t", "w", "x", "y", "z", "zh",
}

var uwlYm = []string{
	"ang", "a", "ai", "an", "ang", "ao", "e", "ei", "en", "eng", "er",
	"i", "ia", "ian", "iang", "iao", "ie", "in", "ing", "iong", "iu",
	"o", "ong", "ou", "u",
	"ua", "uai", "uan", "uang", "ue", "ui", "un", "uo", "v",
}

func ParseZiguangUwl(rd io.Reader) []PyEntry {
	ret := make([]PyEntry, 0, 1e5)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)
	var tmp []byte

	r.Seek(0xC10, 0)
	for r.Len() > 7 {
		// 读2个字节
		space := make([]byte, 2)
		r.Read(space)
		// 2字节都为空
		if space[0] == 0 && space[1] == 0 {
			continue
		}
		// if space[0]%2 == 0 || // 1字节是偶数
		// 	space[1]>>4%2 != 0 || // 2字节前4位是奇数
		// 	space[1]%0x10 == 0 { // 2字节后4位是0
		// 	r.Seek(14, 1)
		// 	continue
		// }
		r.Seek(-2, 1) // 回退2字节

		// 读 16 字节
		tmp = make([]byte, 16)
		r.Read(tmp)
		flag := true // 是否丢弃
		for i := 0; i < 4; i++ {
			// 后两个字节相差不超过1
			if tmp[4*i+2]-tmp[4*i+3] < 2 {
				continue
			} else {
				flag = false
			}
		}
		if flag {
			continue
		}
		r.Seek(-16, 1)

		// 正式读
		head := make([]byte, 4)
		r.Read(head)
		// 词长 * 2
		wordLen := head[0]%0x80 - 1
		// 拼音长
		codeLen := head[1]<<4>>4*2 + head[0]/0x80

		// 频率
		freq := BytesToInt(head[2:])
		// fmt.Println(freqSli, freq)

		// 拼音
		code := make([]string, 0, codeLen)
		for i := 0; i < int(codeLen); i++ {
			bsm, _ := r.ReadByte()
			bym, _ := r.ReadByte()
			smIdx := bsm & 0x1F
			ymIdx := (bsm >> 5) + (bym << 3)
			// fmt.Println(bsm, bym, smIdx, ymIdx)
			if bym >= 0x10 || smIdx >= 24 || ymIdx >= 34 {
				break
			}
			code = append(code, uwlSm[smIdx]+uwlYm[ymIdx])
			// fmt.Println(smIdx, ymIdx, uwlSm[smIdx]+uwlYm[ymIdx])
		}

		// 词
		tmp = make([]byte, wordLen)
		r.Read(tmp)
		word := string(DecUtf16le(tmp))
		// fmt.Println(string(word))

		ret = append(ret, PyEntry{word, code, freq})
	}
	return ret
}
