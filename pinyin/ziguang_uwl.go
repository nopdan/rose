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
	ret := make([]PyEntry, 0, 0xff)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)
	r.Seek(2, 0)
	// 编码格式，08 为 GBK，09 为 UTF-16LE
	encoding, _ := r.ReadByte()

	// 分段
	r.Seek(0x48, 0)
	partLen := ReadUint32(r)
	for i := 0; i < partLen; i++ {
		r.Seek(0xC00+int64(i)<<10, 0)
		ret = parseZgUwlPart(r, ret, encoding)
	}

	return ret
}

func parseZgUwlPart(r *bytes.Reader, ret []PyEntry, e byte) []PyEntry {
	r.Seek(12, 1)
	// 词条占用字节数
	max := ReadUint32(r)
	// 当前字节
	curr := 0
	for curr < max {
		head := make([]byte, 4)
		r.Read(head)
		// 词长 * 2
		wordLen := head[0]%0x80 - 1
		// 拼音长
		codeLen := head[1]<<4>>4*2 + head[0]/0x80
		// 频率
		freq := BytesToInt(head[2:])
		// fmt.Println(freqSli, freq)
		curr += int(4 + wordLen + codeLen*2)

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
		tmp := make([]byte, wordLen)
		r.Read(tmp)
		var word string
		switch e {
		case 0x08:
			word, _ = Decode(tmp, "gbk")
		case 0x09:
			word, _ = Decode(tmp, "utf16")
		}
		// fmt.Println(string(word))
		ret = append(ret, PyEntry{word, code, freq})
	}
	return ret
}
