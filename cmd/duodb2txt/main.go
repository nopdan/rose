package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Printf("多多输入法.duodb 转纯文本\n作者：单单 q37389732\n\n")
	var input string
	if len(os.Args) != 2 {
		fmt.Println("输入词库路径：")
		fmt.Scanln(&input)
	} else {
		input = os.Args[1]
	}
	data, err := os.ReadFile(input)
	if err != nil {
		os.Exit(0)
	}
	r := bytes.NewReader(data)

	wl := Unmarshal(r)
	var buf bytes.Buffer
	for _, v := range wl {
		buf.WriteString(v.Word)
		buf.WriteByte('\t')
		buf.WriteString(v.Code)
		buf.WriteByte('\n')
	}
	ext := filepath.Ext(input)
	name := strings.TrimSuffix(input, ext) + ".txt"
	os.WriteFile(name, buf.Bytes(), 0644)
}

type Entry struct {
	Word string
	Code string
}

func Unmarshal(r *bytes.Reader) []*Entry {
	d := make([]*Entry, 0, r.Size()>>8)
	r.Seek(0x4086C, 0)
	offsetList := make([]uint32, 0, 12)
	for {
		offset := ReadUint32(r)
		if offset == 0 {
			break
		}
		offsetList = append(offsetList, offset)
	}
	for _, offset := range offsetList {
		r.Seek(int64(offset), 0)
		r.Seek(4, 1)
		codeLen := ReadIntN(r, 1)
		code := string(ReadN(r, codeLen))
		wordSize := ReadIntN(r, 2)
		word := string(ReadN(r, wordSize))
		d = append(d, &Entry{word, code})
	}
	return d
}
