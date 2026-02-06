package kafan_pinyin_bak

import (
	"bytes"
	"io"
	"strings"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

type Kafan struct {
	model.BaseFormat
	pyList []string
}

func New() *Kafan {
	f := &Kafan{
		BaseFormat: model.BaseFormat{
			ID:          "kafan_pinyin_bak",
			Name:        "卡饭拼音备份词库",
			Type:        model.FormatTypePinyin,
			Extension:   ".dict",
			Description: "卡饭拼音备份词库格式",
		},
	}

	f.pyList = pyList
	return f
}

func (f *Kafan) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	// 0x48 or 0x68
	r.Seek(0x48, io.SeekStart)
	head := string(r.ReadN(0x10))
	// 版本不匹配
	if !strings.HasPrefix(head, "ProtoDict1") {
		// 有的词库是在 0x68
		r.Seek(0x68, io.SeekStart)
		head = string(r.ReadN(0x10))
		if !strings.HasPrefix(head, "ProtoDict1") {
			f.Infof("文件头格式错误，不是有效的卡饭拼音备份文件\n")
			return nil, nil
		}
	}

	di := make([]*model.Entry, 0, 0xff)
	// 读取一个词
	for r.Len() > 0x28 {
		// 词库中间可能夹杂这个
		dictType := r.ReadN(0x10)
		if !bytes.HasPrefix(dictType, []byte("kf_pinyin")) {
			r.Seek(-0x10, io.SeekCurrent)
		}

		// 读取编码占用的字节
		codeBytes := make([]byte, 0, 0x28)
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

		// 转换为拼音
		pinyin := make([]string, 0, 2)
		// 每 0x28 个字节为一个音
		for i := len(codeBytes) % 0x28; i < len(codeBytes); i += 0x28 {
			idx := util.Bytes2Int(codeBytes[i : i+4])
			py := f.lookup(idx, r)
			if py == "" {
				f.Infof("codeBytes: %v\n", codeBytes)
			} else if py != " " {
				pinyin = append(pinyin, py)
			}
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
			f.Infof("读取词组错误, size: 0x%x, offset: 0x%x\n", size, int(r.Size())-r.Len())
			return nil, nil
		}

		// 下面读取词，词是按照8字节对齐的
		wordBytes := r.ReadN(size)
		if len(wordBytes)%8 != 0 {
			r.Seek(int64(8-len(wordBytes)%8), io.SeekCurrent)
		}
		word := string(wordBytes)
		// di = append(di, &Entry{
		// 	Word:   word,
		// 	Pinyin: pinyin,
		// 	Freq:   1,
		// })
		if py := f.filter(word, pinyin); len(py) > 0 {
			di = append(di, model.NewEntry(word).WithMultiCode(py...).WithFrequency(1))
			f.Debugf("词组: %s, 拼音: %v\n", word, py)
		}
	}
	return di, nil
}

func (k *Kafan) filter(word string, pinyin []string) []string {
	wordRunes := []rune(word)
	// 过滤单字
	if len(wordRunes) <= 1 {
		return nil
	}
	if len(wordRunes) == len(pinyin) {
		return pinyin
	}
	if len(wordRunes) < len(pinyin) {
		return pinyin[len(pinyin)-len(wordRunes):]
	}
	if len(wordRunes) > len(pinyin) {
		//! TODO
		// enc := encoder.NewPinyin()
		// pre := string(wordRunes[:len(wordRunes)-len(pinyin)])
		// res := append(enc.Encode(pre), pinyin...)
		// return res
	}
	return nil
}

func (k *Kafan) lookup(idx int, r *model.Reader) string {
	if len(k.pyList) <= idx {
		k.Infof(
			"index out of range: %d > %d, offset: 0x%x\n",
			idx,
			len(k.pyList)-1,
			int(r.Size())-r.Len(),
		)
		return ""
	}
	return k.pyList[idx]
}
