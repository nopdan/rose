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
			// 每40个字节为一个字母
			tmp := ReadN[int](r, 0x28) // 40
			// 判断前8个字节决定是否结束
			if bytes.Equal(tmp[:8], []byte{4, 0, 0, 0, 3, 0, 1, 0}) {
				break
			}
			codeBytes = append(codeBytes, tmp[0])
		}

		// 跳过未知的8字节
		r.Seek(8, io.SeekCurrent)
		// 下面读取词，词是按照8字节对齐的
		wordBytes := make([]byte, 0, 8)
		for {
			// 每次读8字节
			b := ReadN[int](r, 8)
			wordBytes = append(wordBytes, b...)
			// 如果最后一个字节是0则结束
			if b[7] == 0 {
				break
			}
		}
		// 去除末尾的 0
		for i := len(wordBytes) - 1; i >= 0 && wordBytes[i] == 0; i-- {
			wordBytes = wordBytes[:i]
		}
		word := string(wordBytes)

		di = append(di, &Entry{
			Word: word,
			Code: string(codeBytes),
		})
	}
	return di
}
