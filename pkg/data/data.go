package data

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"os"

	"github.com/nopdan/ku"
)

//go:embed res/wubilex.bin
var WubiLex []byte

//go:embed res/pinyin.bin
var Pinyin []byte

//go:embed res/duoyin.bin
var Duoyin []byte

//go:embed res/correct.bin
var Correct []byte

func compress(input, output string) {
	data, _ := os.ReadFile(input)
	var buffer bytes.Buffer
	w := gzip.NewWriter(&buffer)
	w.Write(data)
	w.Flush()
	f, _ := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	buffer.WriteTo(f)
}

func Decompress(data []byte) io.Reader {
	brd := bytes.NewReader(data)
	zrd, err := gzip.NewReader(brd)
	if err != nil {
		panic(zrd)
	}
	rd := ku.NewReader(zrd)
	return rd
}
