package pinyin

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestKafan(t *testing.T) {
	k := NewKafan()
	b, _ := os.ReadFile("./sample/拼音词库备份2023-08-19.dict")
	wl := k.Unmarshal(bytes.NewReader(b))
	fmt.Println(wl)
}

func TestReadPinyin(t *testing.T) {
	k := NewKafan()
	b, _ := os.ReadFile("./sample/拼音词库备份2023-08-19.dict")
	wl := k.readPinyin(bytes.NewReader(b))
	fmt.Println(wl)
}
func TestReadTest(t *testing.T) {
	k := NewKafan()
	b, _ := os.ReadFile("./sample/拼音词库备份2023-08-19.dict")
	k.Test(bytes.NewReader(b))
}
