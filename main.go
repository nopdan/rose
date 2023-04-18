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

var (
	input         string // 输入词库
	output        string // 保存路径
	input_format  string // 输入词库的格式
	output_format string // 保存的词库格式
)

func main() {
	switch len(os.Args) {
	case 1:
		ask()
		return
	case 2:
		switch os.Args[1] {
		case "-v", "version":
			fmt.Println("蔷薇词库转换v1.1.1\nhttps://github.com/flowerime/rose")
			return
		case "-h", "help":
			fmt.Printf("Usage: .\\rose.exe [input] [input_format]:[output_format] [output]\n")
			fmt.Printf("Example: .\\rose.exe .\\sogou.scel scel:rime rime.dict.yaml\n")
			return
		}
	}

	if len(os.Args) >= 3 {
		goto NORMAL
	}
	wrong()
	return

NORMAL:
	input = os.Args[1]
	fm := strings.Split(os.Args[2], ":")
	if len(fm) == 2 {
		input_format = fm[0]
		output_format = fm[1]
	}
	if len(os.Args) > 3 {
		output = os.Args[3]
	}
	// fmt.Println(input, output, input_format, output_format)
	convert(input, output, input_format, output_format)
}

func convert(input, output, input_format, output_format string) {

	ifm := rose.DetectFormat(input, input_format)
	fmt.Println("输入词库:", input)
	fmt.Println("词库格式:", ifm.GetDict().Name, input_format)
	d := rose.Parse(input, input_format)
	fmt.Println("解析成功: 词条数 ", len(d.WordLibrary))
	fmt.Println()

	ofm := rose.DetectFormat(output, output_format)
	if output == "" {
		output = filepath.Base(input) + "_" + ofm.GetDict().Name
	}
	fmt.Println("输出词库:", output)
	fmt.Println("词库格式:", ofm.GetDict().Name, output_format)
	data := ofm.GenFrom(d.WordLibrary)
	if len(data) == 0 {
		return
	}
	if err := os.WriteFile(output, data, 0666); err == nil {
		fmt.Println("转换成功")
	}
}

func wrong() {
	fmt.Println("输入参数有误")
}

// 交互式
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

	input = askOne("输入词库：")
	input_format = askOne("词库格式：")
	output = askOne("输出词库：")
	output_format = askOne("词库格式：")

	convert(input, output, input_format, output_format)
	fmt.Scanln()
}
