package pinyin

import (
	"bytes"
	"io"
	"io/ioutil"

	"golang.org/x/text/encoding/unicode"
)

func ParseSougouScel(rd io.Reader) []Pinyin {
	ret := make([]Pinyin, 0, 1e5)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)

	// utf-16le 转换器
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()

	// 读词条数
	tmp := make([]byte, 4)
	r.Seek(0x124, 0)
	r.Read(tmp)
	dictLen := bytesToInt(tmp)

	// 拼音表偏移量
	r.Seek(0x1540, 0)

	// 前两个字节是拼音表长度，413
	tmp = make([]byte, 2)
	r.Read(tmp)
	pyTableLen := bytesToInt(tmp)
	pyTable := make([]string, pyTableLen)
	// fmt.Println("拼音表长度", pyTableLen)

	// 丢掉两个字节
	r.Read(tmp)

	// 读拼音表
	for i := 0; i < pyTableLen; i++ {
		// 索引
		tmp := make([]byte, 2)
		r.Read(tmp)
		idx := bytesToInt(tmp)

		// 拼音长度
		r.Read(tmp)
		pyLen := bytesToInt(tmp)

		// 拼音 utf-16le
		pySli := make([]byte, pyLen)
		r.Read(pySli)
		py, _ := decoder.Bytes(pySli)

		pyTable[idx] = string(py)
	}

	// 读码表
	for count := 0; count < dictLen; {
		// 重码数（同一串音对应多个词）
		tmp := make([]byte, 2)
		r.Read(tmp)
		repeat := bytesToInt(tmp)

		// 索引数组长
		r.Read(tmp)
		codeLen := bytesToInt(tmp)

		// 读取编码
		var code []string
		for i := 0; 2*i < codeLen; i++ {
			r.Read(tmp)
			theIdx := bytesToInt(tmp)
			if theIdx >= pyTableLen {
				code = append(code, string(byte(theIdx-pyTableLen+97)))
				continue
			}
			code = append(code, pyTable[theIdx])
		}

		// 读取一个或多个词
		count += repeat
		for i := 1; i <= repeat; i++ {
			// 词长
			r.Read(tmp)
			wordLen := bytesToInt(tmp)

			// 读取词
			wordSli := make([]byte, wordLen)
			r.Read(wordSli)
			wordSli, _ = decoder.Bytes(wordSli)
			word := string(wordSli)
			ret = append(ret, Pinyin{word, code, 1})

			// 末尾的补充信息
			r.Read(tmp)
			infoLen := bytesToInt(tmp)
			info := make([]byte, infoLen)
			r.Read(info)
		}
	}
	return ret
}
