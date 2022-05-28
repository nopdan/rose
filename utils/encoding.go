package util

import (
	"bufio"
	"io"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
)

// 将 io流 转换为 utf-8
func Decode(f io.Reader) (io.Reader, error) {

	brd := bufio.NewReader(f)
	buf, _ := brd.Peek(1024)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf) // 检测编码格式
	if err != nil {
		return brd, err
	}
	if cs.Confidence != 100 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}
	rd, err := charset.NewReaderLabel(cs.Charset, brd) // 转换字节流
	return rd, err
}

// utf-16le 转 utf-8
func DecUtf16le(b []byte) []byte {
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	ret, _ := decoder.Bytes(b)
	return ret
}

// utf-8 转 utf-16le
func ToUtf16le(b []byte) []byte {
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	ret, _ := encoder.Bytes(b)
	return ret
}

// GBK 转 utf-8
func DecGBK(b []byte) []byte {
	decoder := simplifiedchinese.GBK.NewDecoder()
	ret, _ := decoder.Bytes(b)
	return ret
}
