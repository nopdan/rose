package pinyin

import (
	"bytes"
	"compress/zlib"
	"log"
	"strings"

	"io"
	"io/ioutil"

	. "github.com/cxcn/dtool/pkg/util"
)

func ParseQqQpyd(filename string) WpfDict {
	data, _ := ioutil.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(WpfDict, 0, r.Len()>>8)
	var tmp []byte

	// 0x38 后跟的是压缩数据开始的偏移量
	r.Seek(0x38, 0)
	startZip := ReadUint32(r)
	// 0x44 后4字节是词条数
	r.Seek(0x44, 0)
	dictLen := ReadUint32(r)
	// 0x60 到zip数据前的一段是一些描述信息
	r.Seek(0x60, 0)
	head := make([]byte, startZip-0x60)
	r.Read(head)
	// headStr, _ := Decode(head, "utf16")
	// fmt.Println(headStr) // 打印描述信息

	// 解压数据
	zrd, err := zlib.NewReader(r)
	if err != nil {
		log.Panic(err)
	}
	defer zrd.Close()
	buf := new(bytes.Buffer)
	buf.Grow(r.Len())
	_, err = io.Copy(buf, zrd)
	if err != nil {
		log.Panic(err)
	}
	// 解压完了
	r.Reset(buf.Bytes())

	for i := 0; i < dictLen; i++ {
		// 指向当前
		r.Seek(int64(10*i), 0)

		// 读码长、词长、索引
		addr := make([]byte, 10)
		r.Read(addr)
		idx := BytesToInt(addr[6:]) // 后4字节是索引
		r.Seek(int64(idx), 0)       // 指向索引
		// 读编码，自带 ' 分隔符
		tmp = make([]byte, addr[0])
		r.Read(tmp)
		code := string(tmp)
		// 读词
		tmp = make([]byte, addr[1])
		r.Read(tmp)
		word, _ := Decode(tmp, "utf16")

		ret = append(ret, WordPyFreq{word, strings.Split(code, "'"), 1})
	}
	return ret
}
