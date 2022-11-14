package pinyin

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/cxcn/dtool/pkg/encoder"
	"github.com/cxcn/dtool/pkg/util"
)

type MspyUDL struct{}

// 自学习词库，纯汉字
func (MspyUDL) Parse(filename string) Dict {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	r.Seek(0xC, 0)
	count := ReadUint32(r)
	ret := make(Dict, 0, count)
	r.Seek(4, 1)
	export_stamp := ReadUint32(r)
	export_time := time.Unix(int64(export_stamp), 0).Add(946684800 * time.Second)
	fmt.Printf("词条数: %d, 时间: %v\n", count, export_time)

	for i := 0; i < count; i++ {
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
		ret = append(ret, Entry{word, encoder.GetPinyin(word), 1})
		// fmt.Printf("时间: %v, 简拼: %s, 词: %s\n", insert_time, string(jianpin), word)
	}
	return ret
}
