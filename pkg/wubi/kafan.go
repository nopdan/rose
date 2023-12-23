package wubi

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Kafan struct {
	Template
}

func init() {
	FormatList = append(FormatList, NewKafan())
}

func NewKafan() *Kafan {
	f := new(Kafan)
	f.Name = "卡饭五笔备份.dict"
	f.ID = "kfwbbak"
	return f
}

func (f *Kafan) Unmarshal(r *bytes.Reader) []*Entry {
	r.Seek(0x48, io.SeekStart)
	head := string(ReadN(r, 0x10))
	// 版本不匹配
	if !strings.HasPrefix(head, "ProtoDict1") {
		fmt.Println("卡饭五笔备份.dict格式错误")
		return nil
	}

	di := make([]*Entry, 0, 0xff)
	// 读取一个词
	for r.Len() > 0x28 {
		// 词库中间可能夹杂这个
		dictType := string(ReadN(r, 8))
		if !strings.HasPrefix(dictType, "wubi86") {
			r.Seek(-8, io.SeekCurrent)
		}

		// 读取编码
		codeBytes := make([]byte, 0, 2)
		for {
			// 每次读取 8 个字节
			tmp := ReadN[int](r, 8)
			// 判断结束
			if bytes.Equal(tmp, []byte{4, 0, 0, 0, 3, 0, 1, 0}) {
				r.Seek(0x20, io.SeekCurrent)
				break
			} else if bytes.Equal(tmp, []byte{0, 0, 0, 0, 3, 0, 1, 0}) {
				r.Seek(0x18, io.SeekCurrent)
				break
			}
			codeBytes = append(codeBytes, tmp...)
		}

		// 转换编码
		codeB := make([]byte, 0, 2)
		// 每 0x28 个字节
		for i := len(codeBytes) % 0x28; i < len(codeBytes); i += 0x28 {
			codeB = append(codeB, codeBytes[i])
		}

		// 跳过未知的4字节
		mark := ReadIntN(r, 4)
		if mark != 1 {
			r.Seek(8, io.SeekCurrent)
		}
		size := ReadIntN(r, 4)
		// 22	3	8
		// 2A	4
		// 32	5
		// 3A	6	8
		// 42	7	8
		// 4A	8	8
		// 52	9	16
		// 6A	12	16
		if size%0x10 == 2 {
			size = (size/0x10)*2 - 1
		} else if size%0x10 == 0xA {
			size = (size / 0x10) * 2
		} else {
			fmt.Printf("读取词组错误, size: 0x%x, offset: 0x%x\n", size, int(r.Size())-r.Len())
			return nil
		}

		// 下面读取词，词是按照8字节对齐的
		wordBytes := ReadN(r, size)
		if len(wordBytes)%8 != 0 {
			r.Seek(int64(8-len(wordBytes)%8), io.SeekCurrent)
		}
		word := string(wordBytes)

		di = append(di, &Entry{
			Word: word,
			Code: string(codeB),
		})
	}
	return di
}
