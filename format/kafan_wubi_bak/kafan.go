package kafan_wubi_bak

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/nopdan/rose/model"
)

type Kafan struct {
	model.BaseFormat
}

func New() *Kafan {
	return &Kafan{
		BaseFormat: model.BaseFormat{
			ID:          "kafan_wubi_bak",
			Name:        "卡饭五笔备份词库",
			Type:        model.FormatTypeWubi,
			Extension:   ".dict",
			Description: "卡饭五笔备份词库格式",
		},
	}
}

func (f *Kafan) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	r.Seek(0x48, io.SeekStart)
	head := string(r.ReadN(0x10))
	// 版本不匹配
	if !strings.HasPrefix(head, "ProtoDict1") {
		return nil, fmt.Errorf("卡饭五笔备份词库格式错误, 版本: %s", head)
	}

	di := make([]*model.Entry, 0, 0xff)
	// 读取一个词
	for r.Len() > 0x28 {
		// 词库中间可能夹杂这个
		dictType := string(r.ReadN(8))
		if !strings.HasPrefix(dictType, "wubi86") {
			r.Seek(-8, io.SeekCurrent)
		}

		// 读取编码
		codeBytes := make([]byte, 0, 2)
		for {
			// 每次读取 8 个字节
			tmp := r.ReadN(8)
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
		mark := r.ReadIntN(4)
		if mark != 1 {
			r.Seek(8, io.SeekCurrent)
		}
		size := r.ReadIntN(4)
		// 22	3	8
		// 2A	4
		// 32	5
		// 3A	6	8
		// 42	7	8
		// 4A	8	8
		// 52	9	16
		// 6A	12	16
		switch size % 0x10 {
		case 2:
			size = (size/0x10)*2 - 1
		case 0xA:
			size = (size / 0x10) * 2
		default:
			return nil, fmt.Errorf("读取词组错误, size: 0x%x, offset: 0x%x", size, int(r.Size())-r.Len())
		}

		// 下面读取词，词是按照8字节对齐的
		wordBytes := r.ReadN(size)
		if len(wordBytes)%8 != 0 {
			r.Seek(int64(8-len(wordBytes)%8), io.SeekCurrent)
		}
		word := string(wordBytes)
		di = append(di, model.NewEntry(word).WithSimpleCode(string(codeB)))
		f.Debugf("词组: %s, 编码: %s\n", word, string(codeB))
	}
	return di, nil
}
