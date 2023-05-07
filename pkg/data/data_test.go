package data

import "testing"

func TestCompress(t *testing.T) {
	compress("wubilex.txt", "res/wubilex.bin")
	compress("../../pinyin-data/pinyin.txt", "res/pinyin.bin")
	compress("../../pinyin-data/duoyin.txt", "res/duoyin.bin")
	compress("../../pinyin-data/correct.txt", "res/correct.bin")
}
