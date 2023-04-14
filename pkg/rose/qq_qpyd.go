package rose

import (
	"bytes"
	"compress/zlib"
	"io"
	"log"
	"strings"
)

type QqQpyd struct{ Dict }

func NewQqQpyd() *QqQpyd {
	d := new(QqQpyd)
	d.IsPinyin = true
	d.IsBinary = true
	d.Name = "QQ 拼音.qpyd"
	d.Suffix = "qpyd"
	return d
}

func (d *QqQpyd) Parse() {
	pyt := make(PyTable, 0, d.size>>8)

	r := bytes.NewReader(d.data)
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
	// headStr, _ := Decode(head, "UTF-16LE")
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

	for i := _u32; i < dictLen; i++ {
		var tmp []byte
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
		word, _ := Decode(tmp, "UTF-16LE")

		pyt = append(pyt, &PinyinEntry{word, strings.Split(code, "'"), 1})
	}
	d.pyt = pyt
}
