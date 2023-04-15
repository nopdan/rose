/*
Copyright © 2023 nopdan <me@nopdan.com>
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flowerime/rose/pkg/rose"
)

func main() {

	switch len(os.Args) {
	case 1:
		ask()
		return
	case 2:
		switch os.Args[1] {
		case "-v", "version":
			fmt.Println("蔷薇词库转换v1.0.1\nhttps://github.com/flowerime/rose")
			return
		case "-h", "help":
			fmt.Printf("Usage: .\\rose.exe [path] [input_format]:[output_format]\n")
			fmt.Printf("Example: .\\rose.exe .\\sogou.scel scel:rime\n")
			return
		}
	case 3:
	default:
		wrong()
		return
	}

	path := os.Args[1]
	format := os.Args[2]
	tmp := strings.Split(format, ":")
	if len(tmp) != 2 {
		wrong()
		return
	}
	iFormat, oFormat := tmp[0], tmp[1]
	convert(path, iFormat, oFormat)
}

func convert(path, iFormat, oFormat string) {
	if iFormat == "" {
		iFormat = "dd"
	}
	d := rose.Parse(path, iFormat)
	if oFormat == "" {
		if d.IsPinyin {
			oFormat = "rime"
		} else {
			oFormat = "dd"
		}
	}
	data := rose.Generate(d, oFormat)
	od := rose.NewFormat(oFormat).GetDict()
	oPath := filepath.Base(path) + "_" + oFormat + "." + od.Suffix
	err := os.WriteFile(oPath, data, 0666)
	if err == nil {
		fmt.Println("转换成功，输出到", oPath)
	}
}

func wrong() {
	fmt.Println("输入参数有误")
}

func ask() {

	askOne := func(hint string) string {
		fmt.Printf("%s\n> ", hint)
		reader := bufio.NewReader(os.Stdin)
		var value string
		value, _ = reader.ReadString('\n')
		value = strings.ReplaceAll(value, "\r", "")
		value = strings.ReplaceAll(value, "\n", "")
		last := value[len(value)-1]
		if len(value) > 3 {
			if last == '"' && value[0] == '"' {
				value = value[1 : len(value)-1]
			} else if len(value) > 4 && last == '\'' && value[:3] == "& '" {
				value = value[3 : len(value)-1]
			}
		}

		fmt.Println()
		return value
	}

	path := askOne("待转换的码表：")
	iFormat := askOne("待转换码表的格式：")
	oFormat := askOne("转换为的格式：")

	convert(path, iFormat, oFormat)
}
