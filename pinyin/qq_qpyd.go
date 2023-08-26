package pinyin

import (
	"bytes"
	"compress/zlib"
	"io"
	"log"
	"strings"

	"github.com/nopdan/rose/util"
)

type QqQpyd struct{ Template }

func NewQqQpyd() *QqQpyd {
	f := new(QqQpyd)
	f.Name = "QQ 拼音.qpyd"
	f.ID = "qpyd"
	return f
}

func (f *QqQpyd) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)

	r.Seek(0x2C, 0)
	startInfo := ReadUint32(r)
	// 0x38 后跟的是压缩数据开始的偏移量
	r.Seek(0x38, 0)
	startZip := ReadUint32(r)
	// 0x44 后4字节是词条数
	r.Seek(0x44, 0)
	dictLen := ReadIntN(r, 4)
	// 0x60 到zip数据前的一段是一些描述信息
	r.Seek(int64(startInfo), 0)
	util.Info(r, startZip-startInfo, "")

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
		addr := ReadN(r, 10)
		idx := BytesToInt(addr[6:]) // 后4字节是索引
		r.Seek(int64(idx), 0)       // 指向索引
		// 读编码，自带 ' 分隔符
		tmp := ReadN(r, int(addr[0]))
		code := string(tmp)
		// 读词
		tmp = ReadN(r, int(addr[1]))
		word := util.DecodeMust(tmp, "UTF-16LE")

		di = append(di, &Entry{word, strings.Split(code, "'"), 1})
	}
	return di
}
