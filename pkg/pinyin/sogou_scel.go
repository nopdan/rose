package pinyin

import (
	"bytes"
	"fmt"

	"github.com/nopdan/rose/pkg/util"
)

type SogouScel struct{ Template }

func init() {
	FormatList = append(FormatList, NewSogouScel(), NewQqQcel())
}
func NewSogouScel() *SogouScel {
	f := new(SogouScel)
	f.Name = "搜狗细胞词库.scel"
	f.ID = "scel"
	return f
}
func NewQqQcel() *SogouScel {
	f := NewSogouScel()
	f.Name = "QQ拼音v6以上.qcel"
	f.ID = "qcel"
	return f
}

func (f *SogouScel) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)

	// 不展开的词条数
	r.Seek(0x120, 0)
	dictLen := ReadUint32(r)
	fmt.Printf("未展开的词条数: %d\n", dictLen)

	// 拼音表偏移量
	r.Seek(0x1540, 0)
	// 前两个字节是拼音表长度，413
	pyTableLen := ReadUint16(r)
	pyTable := make([]string, pyTableLen)
	fmt.Printf("拼音表长度: %d\n", pyTableLen)
	// 丢掉两个字节
	r.Seek(2, 1)

	// 读拼音表
	for i := _u16; i < pyTableLen; i++ {
		// 索引，2字节
		idx := ReadUint16(r)
		// 拼音长度，2字节
		pyLen := ReadUint16(r)
		// 拼音 utf-16le
		pyBytes := ReadN(r, pyLen)
		py := util.DecodeMust(pyBytes, "UTF-16LE")
		pyTable[idx] = string(py)
	}

	// 读码表
	for j := _u32; j < dictLen; j++ {
		// 重码数（同一串音对应多个词）
		repeat := ReadUint16(r)
		// 索引数组长
		pinyinSize := ReadUint16(r)
		// 读取编码
		var pinyin []string
		for i := _u16; i < pinyinSize/2; i++ {
			theIdx := ReadUint16(r)
			if theIdx >= pyTableLen {
				pinyin = append(pinyin, string(byte(theIdx-pyTableLen+97)))
				continue
			}
			pinyin = append(pinyin, pyTable[theIdx])
		}

		// 读取一个或多个词
		for i := _u16 + 1; i <= repeat; i++ {
			// 词长
			wordSize := ReadUint16(r)
			// 读取词
			wordBytes := ReadN(r, wordSize)
			word := util.DecodeMust(wordBytes, "UTF-16LE")
			// 末尾的补充信息，作用未知
			extSize := ReadUint16(r)
			ext := make([]byte, extSize)
			r.Read(ext)

			di = append(di, &Entry{word, pinyin, 1})
		}
	}
	if r.Len() > 16 {
		// 黑名单
		r.Seek(12, 1)
		blackLen := ReadUint16(r)
		fmt.Printf("以下是词库自带的黑名单: \n\n")
		for i := _u16; i < blackLen; i++ {
			wordLen := ReadUint16(r)
			tmp := make([]byte, wordLen*2)
			r.Read(tmp)
			word := util.DecodeMust(tmp, "UTF-16LE")
			fmt.Println(word)
		}
	}
	return di
}
