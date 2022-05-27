package zici

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	. "github.com/cxcn/dtool/utils"
)

func ParseBaiduDef(rd io.Reader) []ZcEntry {
	// type defEntry struct {
	// 	word  string
	// 	code  string
	// 	order int
	// }
	// def := make([]defEntry, 0, 1e5) // 初始化
	ret := make([]ZcEntry, 0, 1e5) // 初始化
	data, _ := ioutil.ReadAll(rd)  // 全部读到内存
	r := bytes.NewReader(data)
	var tmp []byte

	r.Seek(0x6D, 0) // 从 0x6D 开始读
	for r.Len() > 4 {
		codeLen, _ := r.ReadByte() // 编码长度
		wordLen, _ := r.ReadByte() // 词长*2 + 2

		// 读编码
		tmp = make([]byte, int(codeLen))
		r.Read(tmp) // 编码切片
		code := string(tmp)
		spl := strings.Split(code, "=") // 直接删掉 = 号后的
		code = spl[0]
		// 位置
		// order := 1
		// if len(cao) > 1 {
		// 	order, _ = strconv.Atoi(cao[1])
		// }

		// 读词
		tmp = make([]byte, int(wordLen)-2) // -2 后就是字节长度，没有考虑4字节的情况
		r.Read(tmp)
		word := string(DecUtf16le(tmp))
		// def = append(def, defEntry{word, code, order})
		ret = append(ret, ZcEntry{word, code})

		r.Seek(6, 1) // 6个00，1是相对当前位置
	}
	// sort.Slice(def, func(i, j int) bool {
	// 	return def[i].order < def[j].order
	// })
	// for _, v := range def {
	// 	ret = append(ret, ZcEntry{v.word, v.code})
	// }
	return ret
}

func GenBaiduDef(ce []CodeEntry) []byte {
	var buf bytes.Buffer
	// 首字母词条字节数统计
	lengthMap := make(map[byte]int)
	buf.Write(make([]byte, 0x6D, 0x6D))

	for _, v := range ce {
		code := v.Code

		for i, word := range v.Words {
			if i != 0 { // 不在首选的写入位置信息，好像没什么用？
				code = v.Code + "=" + strconv.Itoa(i+1)
			}
			sliWord := ToUtf16le([]byte(word))    // 转为utf-16le
			buf.WriteByte(byte(len(code)))        // 写编码长度
			buf.WriteByte(byte(len(sliWord) + 2)) // 写词字节长+2
			buf.WriteString(code)                 // 写编码
			buf.Write(sliWord)                    // 写词
			buf.Write([]byte{0, 0, 0, 0, 0, 0})   // 写6个0

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
