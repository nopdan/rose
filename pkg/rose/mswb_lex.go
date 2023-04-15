package rose

import (
	"bytes"
	"fmt"

	util "github.com/flowerime/goutil"
)

type MswbLex struct{ Dict }

func NewMswbLex() *MswbLex {
	d := new(MswbLex)
	d.Name = "微软五笔.lex"
	d.Suffix = "lex"
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
	// create_stamp := ReadUint32(r) // 时间戳
	// create_time := time.Unix(int64(create_stamp), 0)
	fmt.Printf("索引表开始：0x%x\n", idx_start)
	fmt.Printf("文件总大小：0x%x\n", total_size)
	// fmt.Printf("时间：%v\n", create_time)

	r.Seek(int64(entry_start), 0)
	for r.Len() > 4 {
		length := ReadUint16(r) // 词条总字节数
		r.Seek(2, 1)            // 未知
		codeLen := ReadUint16(r) << 1
		tmp := make([]byte, codeLen)
		r.Read(tmp)
		code, _ := util.Decode(tmp, "UTF-16LE")
		tmp = make([]byte, length-8-codeLen)
		r.Read(tmp)
		b := bytes.Split(tmp, []byte{0, 0})
		word, _ := util.Decode(b[len(b)-1], "UTF-16LE")
		wl = append(wl, &WubiEntry{word, code, 1})
		// fmt.Println(length, codeLen, code, word, order)
		r.Seek(2, 1)
	}
	// sort.Slice(ret, func(i, j int) bool {
	// 	return ret[i].Order > ret[j].Order
	// })
	// sort.Slice(ret, func(i, j int) bool {
	// 	return ret[i].Code < ret[j].Code
	// })
	d.WordLibrary = wl
}
