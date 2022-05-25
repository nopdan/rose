package pinyin

import (
	"bytes"
	"compress/zlib"
	"log"
	"strings"

	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/text/encoding/unicode"
)

func ParseQqQpyd(rd io.Reader) []Pinyin {
	ret := make([]Pinyin, 0, 1e5)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)

	// utf-16le 转换器
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()

	// 0x38 后跟的是压缩数据开始的偏移量
	r.Seek(0x38, 0)
	tmp := make([]byte, 4)
	r.Read(tmp)
	startZip := bytesToInt(tmp)
	// 0x44 后4字节是词条数
	r.Seek(0x44, 0)
	r.Read(tmp)
	dictLen := bytesToInt(tmp)
	// 0x60 到zip数据前的一段是一些描述信息
	r.Seek(0x60, 0)
	head := make([]byte, startZip-0x60)
	r.Read(head)
	b, _ := decoder.Bytes(head)
	fmt.Println(string(b)) // 打印描述信息

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
		// 读码长、词长、索引
		addr := make([]byte, 10)
		r.Read(addr)
		idx := bytesToInt(addr[6:]) // 后4字节是索引
		r.Seek(int64(idx), 0)       // 指向索引
		codeSli := make([]byte, addr[0])
		r.Read(codeSli)
		wordSli := make([]byte, addr[1])
		r.Read(wordSli)
		wordSli, _ = decoder.Bytes(wordSli)
		ret = append(ret, Pinyin{string(wordSli), strings.Split(string(codeSli), "'"), 1})
		// 指向下一条
		r.Seek(int64(10*(i+1)), 0)
	}
	return ret
}
