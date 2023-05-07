package rose

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "embed"

	"github.com/flowerime/rose/pkg/data"
	"github.com/nopdan/ku"
)

type MswbLex struct {
	Dict
	wwMap map[string]int
}

func NewMswbLex() *MswbLex {
	d := new(MswbLex)
	d.Name = "微软五笔.lex"
	d.Suffix = "lex"
	d.wwMap = make(map[string]int)
	d.LoadWordWeight()
	return d
}

func (d *MswbLex) Parse() {
	wl := make([]Entry, 0, d.size>>8)

	r := bytes.NewReader(d.data)
	r.Seek(0x0c, 0) // 文件头
	idx_start := ReadUint32(r)
	entry_start := ReadUint32(r)
	total_size := ReadUint32(r) // 词库总长度
	r.Seek(4, 1)
	create_stamp := ReadUint32(r) // 时间戳
	create_time := time.Unix(int64(create_stamp), 0)
	fmt.Printf("索引表开始：0x%x\n", idx_start)
	fmt.Printf("文件总大小：0x%x\n", total_size)
	fmt.Printf("时间：%v\n", create_time)

	r.Seek(int64(entry_start), 0)
	for r.Len() > 4 {
		length := ReadUint16(r) // 词条总字节数
		r.Seek(2, 1)
		// weight := ReadUint16(r)
		codeLen := ReadUint16(r) // 有效编码长
		tmp := make([]byte, 8)
		r.Read(tmp)
		tmp = tmp[:codeLen*2]
		code := DecodeMust(tmp, "UTF-16LE")

		tmp = make([]byte, length-16)
		r.Read(tmp)
		word := DecodeMust(tmp, "UTF-16LE")
		wl = append(wl, &WubiEntry{word, code, 1})
		// fmt.Println(length, codeLen, code, word, weight)
		r.Seek(2, 1)
	}
	d.WordLibrary = wl
}

func (d *MswbLex) GenFrom(wl WordLibrary) []byte {
	var buf bytes.Buffer
	buf.WriteString("imscwubi")
	buf.Write([]byte{1, 0, 1, 0})
	buf.Write([]byte{0x40, 0, 0, 0})
	buf.Write([]byte{0xA8, 0, 0, 0})
	buf.Write(make([]byte, 4)) // hold total size
	now := time.Now()
	buf.Write(ku.To4Bytes(now.Unix()))
	buf.Write(make([]byte, 0x40-0x1C))
	buf.Write(make([]byte, 0xA8-0x40))

	codeWeight := make(map[byte]uint32)
	for _, v := range wl {
		w := v.GetWord()
		word := EncodeMust(w, "UTF-16LE")
		length := 16 + len(word)
		buf.Write(ku.To2Bytes(length))
		weight, ok := d.wwMap[w]
		if !ok {
			weight = 60000
		}
		buf.Write(ku.To2Bytes(weight))
		c := v.GetCode()
		codeLen := 4
		if len(c) < 4 {
			codeLen = len(c)
		}
		buf.Write(ku.To2Bytes(codeLen))
		tmp := EncodeMust(c, "UTF-16LE")
		code := make([]byte, 4)
		codeWeight[code[0]]++
		copy(code, tmp)
		buf.Write(code)
		buf.Write(word)
		buf.Write([]byte{0, 0})
	}
	b := buf.Bytes()
	copy(b[0x14:0x18], ku.To4Bytes(len(b)))
	for i := 0; i < 26; i++ {
		copy(b[0x40+4*i:0x40+4*(i+1)], ku.To4Bytes(codeWeight['a'+byte(i)]))
	}
	return b
}

func (d *MswbLex) LoadWordWeight() {
	bin := data.WubiLex
	rd := data.Decompress(bin)

	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		text := scan.Text()
		tmp := strings.Split(text, "\t")
		fmt.Println(text)
		if len(tmp) == 2 {
			word := tmp[0]
			weight, _ := strconv.Atoi(tmp[1])
			d.wwMap[word] = weight
		}
	}
}
