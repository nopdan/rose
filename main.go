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

	if len(os.Args) != 3 {
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
	os.WriteFile(filepath.Base(path)+"_"+oFormat+"."+od.Suffix, data, 0666)
}

func wrong() {
	fmt.Println("输入参数有误")
}
