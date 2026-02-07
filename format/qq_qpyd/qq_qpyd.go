package qq_qpyd

import (
	"compress/zlib"
	"encoding/binary"
	"io"
	"strings"
	"time"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var utf16 = util.NewEncoding("UTF-16LE")

type Qpyd struct {
	model.BaseFormat
}

func New() *Qpyd {
	return &Qpyd{
		BaseFormat: model.BaseFormat{
			ID:          "qpyd",
			Name:        "QQ拼音分类词库",
			Type:        model.FormatTypePinyin,
			Extension:   ".qpyd",
			Description: "QQ拼音v6以下分类词库格式",
		},
	}
}

func (f *Qpyd) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	entries := make([]*model.Entry, 0, r.Size()>>8)

	magic := r.ReadN(8)
	f.Infof("魔数: %X\n", magic)

	r.Seek(0x18, 0)
	filetime := r.ReadN(8)
	t := toTime(filetime)
	f.Infof("时间: %s\n", t.Format(time.RFC3339))

	r.Seek(8, 1) // 重复的时间戳

	version := r.ReadUint32()
	f.Infof("版本: %d\n", version)

	r.Seek(0x2C, 0)
	infoOffset := r.ReadUint32()
	infoSize := r.ReadUint32()
	f.Infof("信息偏移: 0x%X, 大小: 0x%X\n", infoOffset, infoSize)

	r.Seek(0x38, 0)
	zipOffset := r.ReadUint32()
	zipSize := r.ReadUint32()
	f.Infof("压缩数据偏移: 0x%X, 大小: 0x%X\n", zipOffset, zipSize)

	unzipSize := r.ReadUint32()
	f.Infof("解压后大小: 0x%X\n", unzipSize)

	r.Seek(0x44, 0)
	entryCount := r.ReadIntN(4)
	f.Infof("词条数量: %d\n", entryCount)

	r.Seek(int64(infoOffset), 0)
	info := r.ReadStringEnc(int(infoSize), utf16)
	f.Infof("%s\n", info)

	// 解压数据
	r.Seek(int64(zipOffset), 0)
	zrd, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zrd.Close()
	b, _ := io.ReadAll(zrd) // 确保解压完毕
	// if f.LogLevel >= model.LogDebug {
	// 	ff, _ := os.OpenFile("test.bin", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	// 	ff.Write(b)
	// 	ff.Close()
	// }
	r.Reset(b)

	for i := range entryCount {
		// 指向当前
		r.Seek(int64(10*i), 0)

		sizes := r.ReadN(2)
		r.Seek(4, 1)             // 未知
		offset := r.ReadUint32() // 词条偏移

		r.Seek(int64(offset), 0)
		// 读编码，自带 ' 分隔符
		code := r.ReadString(int(sizes[0]))
		// 读词
		word := r.ReadStringEnc(int(sizes[1]), utf16)

		pinyins := strings.Split(code, "'")
		entry := model.NewEntry(word).WithMultiCode(pinyins...)
		entries = append(entries, entry)

		f.Debugf("%s\t%v\n", word, code)
	}
	return entries, nil
}

// toTime converts an 8-byte Windows FILETIME (100ns since 1601-01-01) to time.Time.
func toTime(t []byte) time.Time {
	if len(t) < 8 {
		return time.Unix(0, 0)
	}
	ft := int64(binary.LittleEndian.Uint64(t))
	if ft == 0 {
		return time.Unix(0, 0)
	}
	// FILETIME epoch (1601-01-01) to Unix epoch (1970-01-01) in 100ns ticks.
	const unixEpochInFiletime = int64(116444736000000000)
	unixTicks := ft - unixEpochInFiletime
	secs := unixTicks / 10000000
	nanos := (unixTicks % 10000000) * 100
	return time.Unix(secs, nanos)
}
