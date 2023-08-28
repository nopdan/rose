package wubi

import (
	"bytes"
	"io"
	"sort"
	"strings"

	"github.com/nopdan/rose/pkg/util"
)

type BaiduDef struct{ Template }

func init() {
	FormatList = append(FormatList, NewBaiduDef())
}
func NewBaiduDef() *BaiduDef {
	f := new(BaiduDef)
	f.Name = "百度手机自定义方案.def"
	f.ID = "def"
	f.CanMarshal = true
	return f
}

func (BaiduDef) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)

	r.Seek(0x6D, io.SeekStart) // 从 0x6D 开始读
	for r.Len() > 4 {
		codeLen, _ := r.ReadByte()  // 编码长度
		wordSize, _ := r.ReadByte() // 词长*2 + 2

		// 读编码
		codeBytes := ReadN(r, codeLen) // 编码切片
		code := string(codeBytes)
		code = strings.Split(code, "=")[0] // 直接删掉 = 号后的
		// 读词
		wordBytes := ReadN(r, wordSize-2) // -2 后就是字节长度，没有考虑4字节的情况
		word := util.DecodeMust(wordBytes, "UTF-16LE")

		di = append(di, &Entry{
			Word: word,
			Code: code,
		})
		r.Seek(6, io.SeekCurrent) // 6个00
	}
	return di
}

func (BaiduDef) Marshal(di []*Entry, hasRank bool) []byte {
	// 按编码排序
	sort.SliceStable(di, func(i, j int) bool {
		return di[i].Code < di[j].Code
	})
	// 首字母词条字节数统计
	lengthMap := make(map[byte]int)
	var buf bytes.Buffer
	buf.Write(make([]byte, 0x6D))
	for _, v := range di {
		word, code := v.Word, v.Code
		wordBytes := util.EncodeMust(word, "UTF-16LE")
		buf.WriteByte(byte(len(code)))          // 写编码长度
		buf.WriteByte(byte(len(wordBytes) + 2)) // 写词字节长+2
		buf.WriteString(code)                   // 写编码
		buf.Write(wordBytes)                    // 写词
		buf.Write(make([]byte, 6))              // 写6个0

		// 编码长度 + 词字节长 + 6，不包括长度本身占的2个字节
		lengthMap[code[0]] += len(code) + len(wordBytes) + 2 + 6
	}

	// 文件头
	byteList := make([]byte, 0, 0x6D)
	byteList = append(byteList, 0) // 第一个字节可能是最大码长？
	// 长度累加
	var currNum int
	for i := 0; i <= 26; i++ {
		currNum += lengthMap[byte(i+0x60)]
		currBytes := To4Bytes(currNum)
		byteList = append(byteList, currBytes...)
	}
	// 替换文件头
	data := buf.Bytes()
	copy(data, byteList)
	return data
}
