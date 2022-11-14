package table

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/cxcn/dtool/pkg/util"
)

type MswbLex struct{}

func (MswbLex) Parse(path string) Table {
	data, _ := os.ReadFile(path)
	r := bytes.NewReader(data)
	ret := make(Table, 0, r.Len()>>8)

	r.Seek(0x0c, 0) // 文件头
	idx_start := ReadUint32(r)
	entry_start := ReadUint32(r)
	total_size := ReadUint32(r)   // 词库总长度
	create_stamp := ReadUint32(r) // 时间戳
	create_time := time.Unix(int64(create_stamp), 0)
	fmt.Println(idx_start, entry_start, total_size, create_time)

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
		ret = append(ret, Entry{word, code, 1})
		// fmt.Println(length, codeLen, code, word, order)
		r.Seek(2, 1)
	}
	// sort.Slice(ret, func(i, j int) bool {
	// 	return ret[i].Order > ret[j].Order
	// })
	// sort.Slice(ret, func(i, j int) bool {
	// 	return ret[i].Code < ret[j].Code
	// })
	return ret
}
