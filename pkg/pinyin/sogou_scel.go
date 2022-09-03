package pinyin

import (
	"bytes"
	"os"

	. "github.com/cxcn/dtool/pkg/util"
)

func ParseSogouScel(filename string) WpfDict {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(WpfDict, 0, r.Len()>>8)
	var tmp []byte

	// 不展开的词条数
	r.Seek(0x120, 0)
	dictLen := ReadUint32(r)

	// 拼音表偏移量
	r.Seek(0x1540, 0)

	// 前两个字节是拼音表长度，413
	pyTableLen := ReadUint16(r)
	pyTable := make([]string, pyTableLen)
	// fmt.Println("拼音表长度", pyTableLen)

	// 丢掉两个字节
	r.Seek(2, 1)

	// 读拼音表
	for i := 0; i < pyTableLen; i++ {
		// 索引，2字节
		idx := ReadUint16(r)
		// 拼音长度，2字节
		pyLen := ReadUint16(r)
		// 拼音 utf-16le
		tmp = make([]byte, pyLen)
		r.Read(tmp)
		py, _ := Decode(tmp, "utf16")
		//
		pyTable[idx] = string(py)
	}

	// 读码表
	for j := 0; j < dictLen; j++ {
		// 重码数（同一串音对应多个词）
		repeat := ReadUint16(r)

		// 索引数组长
		pinyinSize := ReadUint16(r)

		// 读取编码
		var pinyin []string
		for i := 0; i < pinyinSize/2; i++ {
			theIdx := ReadUint16(r)
			if theIdx >= pyTableLen {
				pinyin = append(pinyin, string(byte(theIdx-pyTableLen+97)))
				continue
			}
			pinyin = append(pinyin, pyTable[theIdx])
		}

		// 读取一个或多个词
		for i := 1; i <= repeat; i++ {
			// 词长
			wordSize := ReadUint16(r)

			// 读取词
			tmp = make([]byte, wordSize)
			r.Read(tmp)
			word, _ := Decode(tmp, "utf16")

			// 末尾的补充信息，作用未知
			extSize := ReadUint16(r)
			ext := make([]byte, extSize)
			r.Read(ext)

			ret = append(ret, WordPyFreq{word, pinyin, 1})
		}
	}
	if r.Len() < 16 {
		return ret
	}

	// 黑名单
	r.Seek(12, 1)
	blackLen := ReadUint16(r)
	var black_list bytes.Buffer
	for i := 0; i < blackLen; i++ {
		wordLen := ReadUint16(r)
		tmp = make([]byte, wordLen*2)
		r.Read(tmp)
		word, _ := Decode(tmp, "utf16")
		black_list.WriteString(word)
		black_list.WriteByte('\n')
	}
	// os.WriteFile("black_list.txt", black_list.Bytes(), 0777)
	return ret
}
