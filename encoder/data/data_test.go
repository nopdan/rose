package data

import (
	"bufio"
	"os"
	"testing"
)

func TestCompress(t *testing.T) {
	compress("wubilex.txt", "res/wubilex.bin")
	compress("../../pinyin-data/pinyin.txt", "res/pinyin.bin")
	compress("../../pinyin-data/duoyin.txt", "res/duoyin.bin")
	compress("../../pinyin-data/correct.txt", "res/correct.bin")
}

func TestDecompress(t *testing.T) {
	buf := bufio.NewReader(Duoyin)
	f, _ := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	buf.WriteTo(f)
}
