package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nopdan/rose/pkg/wubi"
)

func main() {
	fmt.Printf("多多.dmg 转纯文本\n作者：单单 q37389732\n\n")
	var input string
	if len(os.Args) != 2 {
		fmt.Println("输入词库路径：")
		fmt.Scanln(&input)
	} else {
		input = os.Args[1]
	}
	data, err := os.ReadFile(input)
	if err != nil {
		os.Exit(1)
	}
	r := bytes.NewReader(data)

	f := wubi.NewDDdmg()
	di := f.Unmarshal(r)
	var buf bytes.Buffer
	for _, v := range di {
		buf.WriteString(v.Word)
		buf.WriteByte('\t')
		buf.WriteString(v.Code)
		buf.WriteByte('\n')
	}
	ext := filepath.Ext(input)
	name := strings.TrimSuffix(input, ext) + ".txt"
	os.WriteFile(name, buf.Bytes(), 0644)
}
