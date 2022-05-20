package dtool

import (
	"bufio"
	"io"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
)

type A [][]byte

// 将 io流 转换为 utf-8
func ReadFile(f io.Reader) (io.Reader, error) {

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

// 笛卡尔积
func Product(a A) []string {
	if len(a) <= 1 {
		return []string{}
	}
	res := make(A, 0, len(a))
	for _, v := range a[0] {
		res = append(res, []byte{v})
	}
	for i := 1; i < len(a); i++ {
		res = productOne(res, a[i])
	}
	ret := make([]string, 0, len(res))
	for _, v := range res {
		ret = append(ret, string(v))
	}
	return ret
}

func productOne(a A, b []byte) A {
	ret := make([]string, 0, len(a)*len(b))
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			ret = append(ret, string(append(a[i], b[j])))
		}
	}
	tmp := make(A, 0, len(a)*len(b))
	for _, v := range ret {
		tmp = append(tmp, []byte(v))
	}
	return tmp
}
