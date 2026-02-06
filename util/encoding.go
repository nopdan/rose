package util

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

type Encoding struct {
	enc *encoding.Encoder
	dec *encoding.Decoder
}

// 无异常处理，必须输入正确的编码名称
// 例如：GB18030, UTF-8, GBK, Big5 等
func NewEncoding(label string) *Encoding {
	e, _ := charset.Lookup(label)
	if e == nil {
		panic(errors.New(""))
	}
	return &Encoding{
		enc: e.NewEncoder(),
		dec: e.NewDecoder(),
	}
}

func (e *Encoding) Encode(str string) []byte {
	ret, _ := e.enc.Bytes([]byte(str))
	return ret
}

func (e *Encoding) Decode(data []byte) string {
	ret, _ := e.dec.String(string(data))
	return strings.TrimFunc(ret, func(r rune) bool {
		return r == '\x00'
	})
}

// 查找编码
func LookupEnc(label string) encoding.Encoding {
	e, _ := charset.Lookup(label)
	return e
}

// 将 io.Reader 转换为 utf-8 编码
func NewUTF8Reader(rd io.Reader) io.Reader {
	r := bufio.NewReader(rd)
	buf, _ := r.Peek(1024)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf) // 检测编码格式
	if err != nil {
		return r
	}
	if cs.Confidence < 95 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}
	e := LookupEnc(cs.Charset)
	return e.NewDecoder().Reader(r)
}

// 读取文件，自动识别编码
func Read(path string) (io.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return f, err
	}
	return NewUTF8Reader(f), nil
}
