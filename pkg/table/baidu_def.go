package table

import (
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"

	"github.com/cxcn/dtool/pkg/util"
)

type BaiduDef struct{}

func (BaiduDef) Parse(filename string) Table {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(Table, 0, r.Len()>>8)
	var tmp []byte

	r.Seek(0x6D, 0) // 从 0x6D 开始读
	for r.Len() > 4 {
		codeLen, _ := r.ReadByte()  // 编码长度
		wordSize, _ := r.ReadByte() // 词长*2 + 2

		// 读编码
		tmp = make([]byte, int(codeLen))
		r.Read(tmp) // 编码切片
		code := string(tmp)
		spl := strings.Split(code, "=") // 直接删掉 = 号后的
		code = spl[0]

		// 读词
		tmp = make([]byte, int(wordSize)-2) // -2 后就是字节长度，没有考虑4字节的情况
		r.Read(tmp)
		word, _ := util.Decode(tmp, "UTF-16LE")
		// def = append(def, defEntry{word, code, order})
		ret = append(ret, Entry{word, code, 1})

		r.Seek(6, 1) // 6个00，1是相对当前位置
	}
	return ret
}

func (BaiduDef) Gen(table Table) []byte {
	jdt := ToJdTable(table)
	var buf bytes.Buffer
	// 首字母词条字节数统计
	lengthMap := make(map[byte]int)
	buf.Write(make([]byte, 0x6D))

	for _, v := range jdt {
		code := v.Code

		for i, word := range v.Words {
			if i != 0 { // 不在首选的写入位置信息，好像没什么用？
				code = v.Code + "=" + strconv.Itoa(i+1)
			}
			sliWord, _ := util.Encode([]byte(word), "UTF-16LE") // 转为utf-16le
			buf.WriteByte(byte(len(code)))                      // 写编码长度
			buf.WriteByte(byte(len(sliWord) + 2))               // 写词字节长+2
			buf.WriteString(code)                               // 写编码
			buf.Write(sliWord)                                  // 写词
			buf.Write(make([]byte, 6))                          // 写6个0

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
		currBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(currBytes, uint32(currNum))
		byteList = append(byteList, currBytes...)
	}
	// 替换文件头
	ret := buf.Bytes()
	copy(ret, byteList)
	return ret
}
