package util

import (
	"bufio"
	"io"
	"os"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

// 读取编码
func LookupEnc(label string) encoding.Encoding {
	e, _ := charset.Lookup(label)
	return e
}

// 将 io.Reader 转换为 utf-8 编码
func NewReader(r io.Reader) io.Reader {
	brd := bufio.NewReader(r)
	buf, _ := brd.Peek(1024)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf) // 检测编码格式
	if err != nil {
		return brd
	}
	// fmt.Printf("cs: %+v\n", cs)
	if cs.Confidence < 95 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}
	e := LookupEnc(cs.Charset)
	return e.NewDecoder().Reader(brd)
}

// 读取文件，自动识别编码
func Read(path string) (io.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return f, err
	}
	return NewReader(f), nil
}

// 编码
func Encode(str, label string) ([]byte, error) {
	e := LookupEnc(label)
	return e.NewEncoder().Bytes([]byte(str))
}
func EncodeMust(str, label string) []byte {
	v, _ := Encode(str, label)
	return v
}

// 解码
func Decode(src []byte, label string) (string, error) {
	e := LookupEnc(label)
	return e.NewDecoder().String(string(src))
}
func DecodeMust(src []byte, label string) string {
	v, _ := Decode(src, label)
	return v
}
