package pinyin

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/cxcn/dtool/utils"
)

func ParseSogouScel(rd io.Reader) []PyEntry {
	ret := make([]PyEntry, 0, 0xff)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)
	var tmp []byte

	// 不展开的词条数
	r.Seek(0x120, 0)
	dictLen := ReadInt(r, 4)

	// 拼音表偏移量
	r.Seek(0x1540, 0)

	// 前两个字节是拼音表长度，413
	pyTableLen := ReadInt(r, 2)
	pyTable := make([]string, pyTableLen)
	// fmt.Println("拼音表长度", pyTableLen)

	// 丢掉两个字节
	r.Seek(2, 1)

	// 读拼音表
	for i := 0; i < pyTableLen; i++ {
		// 索引，2字节
		idx := ReadInt(r, 2)
		// 拼音长度，2字节
		pyLen := ReadInt(r, 2)
		// 拼音 utf-16le
		tmp = make([]byte, pyLen)
		r.Read(tmp)
		py := DecUtf16le(tmp)
		//
		pyTable[idx] = string(py)
	}

	// 读码表
	for j := 0; j < dictLen; j++ {
		// 重码数（同一串音对应多个词）
		repeat := ReadInt(r, 2)

		// 索引数组长
		codeLen := ReadInt(r, 2)

		// 读取编码
		var code []string
		for i := 0; i < codeLen/2; i++ {
			theIdx := ReadInt(r, 2)
			if theIdx >= pyTableLen {
				code = append(code, string(byte(theIdx-pyTableLen+97)))
				continue
			}
			code = append(code, pyTable[theIdx])
		}

		// 读取一个或多个词
		for i := 1; i <= repeat; i++ {
			// 词长
			wordLen := ReadInt(r, 2)

			// 读取词
			tmp = make([]byte, wordLen)
			r.Read(tmp)
			word := string(DecUtf16le(tmp))

			// 末尾的补充信息，作用未知
			extLen := ReadInt(r, 2)
			ext := make([]byte, extLen)
			r.Read(ext)

			ret = append(ret, PyEntry{word, code, 1})
		}
	}
	if r.Len() < 16 {
		return ret
	}

	// 黑名单
	r.Seek(12, 1)
	blackLen := ReadInt(r, 2)
	var black_list bytes.Buffer
	for i := 0; i < blackLen; i++ {
		wordLen := ReadInt(r, 2)
		tmp = make([]byte, wordLen*2)
		r.Read(tmp)
		word := string(DecUtf16le(tmp))
		black_list.WriteString(word)
		black_list.WriteByte('\n')
	}
	ioutil.WriteFile("black_list.txt", black_list.Bytes(), 0777)
	return ret
}
