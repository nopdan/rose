package mswb_lex

import (
	"bufio"
	"bytes"
	_ "embed"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var utf16 = util.NewEncoding("UTF-16LE")

type MswbLex struct {
	model.BaseFormat
}

//go:embed priority.txt
var priorityData string

var wwMap = make(map[string]int)

func New() *MswbLex {
	f := &MswbLex{
		BaseFormat: model.BaseFormat{
			ID:          "mswb_lex",
			Name:        "微软五笔词典",
			Type:        model.FormatTypeWubi,
			Extension:   ".lex",
			Description: "微软五笔输入法词典格式",
		},
	}
	return f
}

func (f *MswbLex) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	entries := make([]*model.Entry, 0, r.Size()>>8)

	r.Seek(0x0c, 0) // 文件头
	idx_start := r.ReadUint32()
	entry_start := r.ReadUint32()
	total_size := r.ReadUint32() // 词库总长度
	r.Seek(4, 1)
	create_stamp := r.ReadUint32() // 时间戳
	create_time := time.Unix(int64(create_stamp), 0)
	f.Infof("索引表开始：0x%x\n", idx_start)
	f.Infof("文件总大小：0x%x\n", total_size)
	f.Infof("时间：%v\n", create_time)

	r.Seek(int64(entry_start), 0)
	for r.Len() > 4 {
		length := r.ReadUint16()   // 词条总字节数
		priority := r.ReadUint16() // 权重
		codeLen := r.ReadUint16()  // 有效编码长
		tmp := make([]byte, 8)
		r.Read(tmp)
		tmp = tmp[:codeLen*2]
		code := utf16.Decode(tmp)

		tmp = r.ReadN(int(length - 16))
		word := utf16.Decode(tmp)

		entries = append(entries, model.NewEntry(word).WithSimpleCode(code))
		f.Debugf("%s\t%s\t%d\n", word, code, priority)
		r.Seek(2, 1)
	}
	return entries, nil
}

func (d *MswbLex) Export(di []*model.Entry, w io.Writer) error {
	if len(wwMap) == 0 {
		d.loadPriority()
	}

	var buf bytes.Buffer
	buf.WriteString("imscwubi")
	buf.Write([]byte{1, 0, 1, 0})
	buf.Write([]byte{0x40, 0, 0, 0})
	buf.Write([]byte{0xA8, 0, 0, 0})
	buf.Write(make([]byte, 4)) // hold total size
	now := time.Now()
	buf.Write(util.To4Bytes(now.Unix()))
	buf.Write(make([]byte, 0x40-0x1C))
	buf.Write(make([]byte, 0xA8-0x40))

	codeWeight := make(map[byte]uint32)
	for _, v := range di {
		w := v.Word
		word := utf16.Encode(w)
		length := 16 + len(word)
		buf.Write(util.To2Bytes(length))
		weight, ok := wwMap[w]
		if !ok {
			weight = 60000
		}
		buf.Write(util.To2Bytes(weight))
		c := v.Code.String()
		codeLen := min(len(c), 4)
		buf.Write(util.To2Bytes(codeLen))
		tmp := utf16.Encode(c)
		code := make([]byte, 4)
		codeWeight[code[0]]++
		copy(code, tmp)
		buf.Write(code)
		buf.Write(word)
		buf.Write([]byte{0, 0})
	}
	b := buf.Bytes()
	copy(b[0x14:0x18], util.To4Bytes(len(b)))
	for i := range 26 {
		copy(b[0x40+4*i:0x40+4*(i+1)], util.To4Bytes(codeWeight['a'+byte(i)]))
	}
	_, err := w.Write(b)
	return err
}

func (d *MswbLex) loadPriority() {
	r := strings.NewReader(priorityData)
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		text := scan.Text()
		tmp := strings.Split(text, "\t")
		// fmt.Println(text)
		if len(tmp) == 2 {
			word := tmp[0]
			weight, _ := strconv.Atoi(tmp[1])
			wwMap[word] = weight
		}
	}
}
