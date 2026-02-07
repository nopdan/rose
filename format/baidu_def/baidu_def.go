package baidu_def

import (
	"bytes"
	"io"
	"slices"
	"strings"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var utf16 = util.NewEncoding("UTF-16LE")

type BaiduDef struct {
	model.BaseFormat
}

func New() *BaiduDef {
	return &BaiduDef{
		BaseFormat: model.BaseFormat{
			ID:          "def",
			Name:        "百度手机自定义方案",
			Type:        model.FormatTypeWubi,
			Extension:   ".def",
			Description: "百度手机五笔自定义方案格式",
		},
	}
}

func (f *BaiduDef) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	entries := make([]*model.Entry, 0, r.Size()>>8)

	r.Seek(0x6D, io.SeekStart) // 从 0x6D 开始读
	for r.Len() > 4 {
		codeLen, _ := r.ReadByte()  // 编码长度
		wordSize, _ := r.ReadByte() // 词组占用字节数，包含末尾 \0

		// 读编码
		code := r.ReadString(int(codeLen))
		code = strings.Split(code, "=")[0] // 直接删掉 = 号后的
		// 读词
		word := r.ReadStringEnc(int(wordSize-2), utf16)

		entry := model.NewEntry(word).WithSimpleCode(code)
		entries = append(entries, entry)
		r.Seek(6, io.SeekCurrent) // 结尾标志 \0 和 4 个未知字节
		f.Debugf("%s\t%s\n", word, code)
	}
	return entries, nil
}

func (f *BaiduDef) Export(entries []*model.Entry, w io.Writer) error {
	slices.SortStableFunc(entries, func(a, b *model.Entry) int {
		return strings.Compare(a.Code.String(), b.Code.String())
	})

	// 首字母词条字节数统计
	lengthMap := make(map[byte]int)
	var buf bytes.Buffer
	buf.Write(make([]byte, 0x6D))

	for _, entry := range entries {
		word := entry.Word
		code := entry.Code.String()

		wordBytes := utf16.Encode(word)
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
		currBytes := util.To4Bytes(currNum)
		byteList = append(byteList, currBytes...)
	}

	// 写入头部
	data := buf.Bytes()
	copy(data[:0x6D], byteList)

	_, err := w.Write(data)
	return err
}
