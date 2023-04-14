package rose

import (
	"bytes"
	"fmt"
	"time"

	util "github.com/flowerime/goutil"
	"github.com/flowerime/rose/pkg/zhuyin"
)

type MspyUDL struct{ Dict }

func NewMspyUDL() *MspyUDL {
	d := new(MspyUDL)
	d.IsPinyin = true
	d.IsBinary = true
	d.Name = "微软拼音自学习词汇.dat"
	d.Suffix = "dat"
	return d
}

// 自学习词库，纯汉字
func (d *MspyUDL) Parse() {
	r := bytes.NewReader(d.data)
	r.Seek(0xC, 0)
	count := ReadUint32(r)
	pyt := make(PyTable, 0, count)
	r.Seek(4, 1)
	export_stamp := ReadUint32(r)
	export_time := time.Unix(int64(export_stamp), 0).Add(946684800 * time.Second)
	fmt.Printf("词条数: %d, 时间: %v\n", count, export_time)

	for i := _u32; i < count; i++ {
		r.Seek(0x2400+60*int64(i), 0)
		// insert_stamp := ReadUint32(r)
		// insert_time := time.Unix(int64(insert_stamp), 0).Add(946684800 * time.Second)
		// jianpin := make([]byte, 4)
		// r.Read(jianpin)
		// ReadUint16(r)
		r.Seek(10, 1)
		wordLen, _ := r.ReadByte()
		r.ReadByte()
		wordSli := make([]byte, wordLen*2)
		r.Read(wordSli)
		word, _ := util.Decode(wordSli, "UTF-16LE")
		pyt = append(pyt, &PinyinEntry{word, zhuyin.Get(word), 1})
		// fmt.Printf("时间: %v, 简拼: %s, 词: %s\n", insert_time, string(jianpin), word)
	}
	d.pyt = pyt
}
