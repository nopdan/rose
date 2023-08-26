package data

import (
	"bytes"
	"compress/flate"
	"embed"
	"io"
	"os"

	"github.com/nopdan/rose/util"
)

//go:embed res
var fs embed.FS

var Duoyin, Pinyin, Correct, Wubilex io.Reader

func init() {
	Duoyin = decompress("res/duoyin.bin")
	Pinyin = decompress("res/pinyin.bin")
	Correct = decompress("res/correct.bin")
	Wubilex = decompress("res/wubilex.bin")
}

func compress(input, output string) {
	data, _ := os.ReadFile(input)
	var buffer bytes.Buffer
	w, _ := flate.NewWriter(&buffer, 9)
	w.Write(data)
	w.Flush()
	f, _ := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	buffer.WriteTo(f)
}

func decompress(name string) io.Reader {
	f, err := fs.Open(name)
	if err != nil {
		panic(err)
	}
	zrd := flate.NewReader(f)
	rd := util.NewReader(zrd)
	return rd
}
