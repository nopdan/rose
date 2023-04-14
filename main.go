/*
Copyright © 2023 nopdan <me@nopdan.com>
*/
package main

import (
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
			fmt.Println("蔷薇词库转换v1.0\nhttps://github.com/flowerime/rose")
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
	var path string
	fmt.Print("待转换的码表：\n> ")
	fmt.Scanf("%s\n", &path)
	fmt.Println()

	var iFormat string
	fmt.Print("待转换码表的格式：\n> ")
	fmt.Scanf("%s\n", &iFormat)
	fmt.Println()

	var oFormat string
	fmt.Print("转换为的格式：\n> ")
	fmt.Scanf("%s\n", &oFormat)
	fmt.Println()

	convert(path, iFormat, oFormat)
}
