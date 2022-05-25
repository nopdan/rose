package zici

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/unicode"
)

func ParseBaiduDef(rd io.Reader) Dict {
	ret := make(Dict, 1e5)       // 初始化
	tmp, _ := ioutil.ReadAll(rd) // 全部读到内存
	r := bytes.NewReader(tmp)
	r.Seek(0x6D, 0) // 从 0x6D 开始读
	// utf-16le 转换
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	for {
		codeLen, err := r.ReadByte() // 编码长度
		wordLen, err := r.ReadByte() // 词长*2 + 2
		if err != nil {
			break
		}
		sliCode := make([]byte, int(codeLen))
		sliWord := make([]byte, int(wordLen)-2) // -2 后就是字节长度，没有考虑4字节的情况

		r.Read(sliCode) // 编码切片
		r.Read(sliWord)

		code := string(sliCode)
		word, _ := decoder.Bytes(sliWord)
		ret.insert(strings.Split(code, "=")[0], string(word))

		r.Seek(6, 1) // 6个00，1是相对当前位置
	}
	return ret
}

func GenBaiduDef(dl []codeAndWords) []byte {
	var buf bytes.Buffer
	// 首字母词条字节数统计
	lengthMap := make(map[byte]int)
	buf.Write(make([]byte, 0x6D, 0x6D))
	// utf-16le 转换
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	for _, v := range dl {
		code := v.code
		for i, word := range v.words {
			if i != 0 { // 不在首选的写入位置信息，好像没什么用？
				code = v.code + "=" + strconv.Itoa(i+1)
			}
			sliWord, _ := encoder.Bytes([]byte(word)) // 转为utf-16le
			buf.WriteByte(byte(len(code)))            // 写编码长度
			buf.WriteByte(byte(len(sliWord) + 2))     // 写词字节长+2
			buf.WriteString(code)                     // 写编码
			buf.Write(sliWord)                        // 写词
			buf.Write([]byte{0, 0, 0, 0, 0, 0})       // 写6个0

			// 编码长度 + 词字节长 + 6，不包括长度本身占的2个字节
			lengthMap[code[0]] += len(code) + len(sliWord) + 2 + 6
		}
	}

	// 文件头
	byteList := make([]byte, 0, 0x6D)
	byteList = append(byteList, 0) // 第一个字节可能是最大码长？
	// 长度累加
	var currNum int
	for i := 0; i <= 26; i++ {
		currNum += lengthMap[byte(i+0x60)]
		// 不知道怎么来的，反正就这样算
		currBytes := []byte{byte(currNum % 0x100), byte((currNum / 0x100) % 0x100),
			byte((currNum / 0x10000) % 0x100), byte((currNum / 0x1000000) % 0x100)}
		byteList = append(byteList, currBytes...)
	}
	// 替换文件头
	ret := buf.Bytes()
	for i := 0; i < len(byteList); i++ {
		ret[i] = byteList[i]
	}
	return ret
}
